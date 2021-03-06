package service

import (
	"context"
	"crypto/tls"
	"email-send-manager/pkg/errors"
	"email-send-manager/pkg/logger"
	"encoding/json"
	"gopkg.in/gomail.v2"
	gotempalte "html/template"
	"io"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/wire"

	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/schema"
)

var SendBatchSet = wire.NewSet(
	wire.Struct(new(SendBatchSrv),
		"TransRepo", "CustomerRepo", "TemplateRepo", "SettingRepo", "RecordRepo",
	))

type SendBatchSrv struct {
	sync.Mutex
	current *current

	TransRepo    *dao.TransRepo
	CustomerRepo *dao.CustomerRepo
	TemplateRepo *dao.TemplateRepo
	SettingRepo  *dao.SettingRepo
	RecordRepo   *dao.RecordRepo
}

func (a *SendBatchSrv) Current(ctx context.Context) (*schema.SendBatchProgress, error) {
	defer a.Unlock()
	a.Lock()
	if a.current == nil {
		return nil, nil
	}

	defer a.current.Unlock()
	a.current.Lock()

	cur := a.current
	if cur.done {
		a.current = nil
	}

	var err string
	if cur.err != nil {
		err = cur.err.Error()
	}

	return &schema.SendBatchProgress {
		StartAt:      cur.startAt,
		TemplateName: cur.tpl,
		Total:        cur.total,
		Success:      cur.success,
		Failure:      cur.failure,
		Error:        err,
	}, nil
}

func (a *SendBatchSrv) TerminateCurrent(ctx context.Context) error {
	defer a.Unlock()
	a.Lock()
	if a.current == nil {
		//return errors.New("无运行中的任务")
		return nil
	}
	a.current.Cancel()
	a.current = nil
	return nil
}

type config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Interval int64  `json:"interval"`
}

func (a *SendBatchSrv) StartSendBatch(ctx context.Context, item schema.SendBatch) (*schema.IDResult, error) {

	if item.Include && len(item.CustomerIDs) == 0 {
		return nil, errors.New400Response("包含选择模式，但未提供客户Id列表")
	}

	a.Lock()
	if a.current != nil && !a.current.done {
		a.Unlock()
		return nil, errors.New400Response("当前批次还未处理完")
	}
	a.current = nil
	a.Unlock()

	setting, err := a.SettingRepo.Get(ctx, "email_send_setting")
	if err != nil {
		return nil, err
	}
	if setting == nil || setting.Value == nil {
		return nil, errors.New500Response("请先设置邮件服务器配置")
	}
	data, err := setting.Value.MarshalJSON()
	if err != nil {
		return nil, errors.New400Response("无效的邮件服务器配置")
	}

	cfg := config{}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, errors.New400Response("无效的邮件服务器配置")
	}

	tpl, err := a.TemplateRepo.Get(ctx, item.TemplateID)
	if err != nil {
		return nil, err
	}
	if tpl == nil {
		return nil, errors.New400Response("没有找到模板 '%d'", item.TemplateID)
	}

	if tpl.Status != 1 {
		return nil, errors.New400Response("该模板处于禁用状态")
	}

	tempalteEngine, err := gotempalte.New("email-content").Parse(tpl.Content)

	if err != nil {
		return nil, errors.New400Response("模板编译失败 '%d'", item.TemplateID)
	}

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sender, err := d.Dial()

	if err != nil {
		return nil, errors.New500Response("连接邮件服务器异常： %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan *schema.Customer)

	cur := current{cancel: cancel, startAt: time.Now(), tpl: tpl.Name}

	cp := customerProvider {
		ctx:     ctx,
		batch:   &item,
		repo:    a.CustomerRepo,
		ch:      ch,
		current: &cur,
	}

	t := task {
		ctx:                     ctx,
		tpl:                     tpl,
		subjects:                strings.Split(tpl.Subject, "_::_"),
		current:                 &cur,
		customers:               ch,
		sender:                  sender,
		templateEngine:          tempalteEngine,
		RecordRepo:              a.RecordRepo,
		interval:                cfg.Interval,
	}

	defer a.Unlock()
	a.Lock()
	if a.current != nil {
		return nil, errors.New400Response("当前批次还未处理完")
	}

	go cp.start()

	go t.start()

	a.current = &cur

	return schema.NewIDResult(0), nil
}

type current struct {
	sync.Mutex

	startAt time.Time
	tpl     string

	total   int32
	success int32
	failure int32
	err     error

	cancel  context.CancelFunc
	done    bool
}

func(c *current) StopOnError(err error) {
	defer c.Unlock()
	c.Lock()
	c.err = err
	c.cancel()
}

func(c *current) Done() {
	defer c.Unlock()
	c.Lock()
	c.done = true
}

func(c *current) Cancel() {
	defer c.Unlock()
	c.Lock()
	c.cancel()
}

func(c *current) Total(num int32) {
	atomic.AddInt32(&c.total, num)
}

func(c *current) getTotal() int32 {
	return atomic.LoadInt32(&c.total)
}

func(c *current) Success() {
	atomic.AddInt32(&c.success, 1)
}

func(c *current) getSuccess() int32 {
	return atomic.LoadInt32(&c.success)
}

func(c *current) Failure() {
	atomic.AddInt32(&c.failure, 1)
}

func(c *current) getFailure() int32 {
	return atomic.LoadInt32(&c.failure)
}

type customerProvider struct {
	ctx     context.Context
	batch   *schema.SendBatch
	repo    *dao.CustomerRepo
	ch      chan <- *schema.Customer
	current *current
}

func (p *customerProvider) start() {

	logger.Debugf("starting customer provider")

	defer close(p.ch)

	pageSize := 100

	params := schema.CustomerQueryParam {
		PaginationParam: schema.PaginationParam{PageSize: pageSize, Pagination: true},
		IDs:             p.batch.CustomerIDs,
		Include:         p.batch.Include,
		Status:          1, // 1 for enabled
	}

	for {
		params.Current += 1

		logger.Debugf("query customers with params: %v", params)

		ret, err := p.repo.Query(p.ctx, params)

		if err != nil {
			p.current.StopOnError(err)
			return
		}

		if params.Current == 1 {
			p.current.Total(int32(ret.PageResult.Total))
		}

		for _, c := range ret.Data {
			select {
				case p.ch <- c:
				case <- p.ctx.Done():
					return
			}
		}

		if len(ret.Data) < pageSize {
			break
		}
	}
}

type task struct {
	ctx                     context.Context

	tpl                     *schema.Template
	subjects                []string

	current                 *current
	customers               <-chan *schema.Customer
	sender                  gomail.SendCloser
	templateEngine          *gotempalte.Template
	RecordRepo              *dao.RecordRepo
	consecutiveSendFailures int
	interval                int64
}

func (s *task) start()  {
	logger.Debugf("starting task")
	defer s.sender.Close()
	for {
		select {
		case customer, ok := <- s.customers:
			if !ok {
				s.current.Done()
				return
			}
			if err := s.sendTo(customer); err != nil {
				logger.Errorf("failure to send email(%s) to %s", s.tpl.Name, customer.Email)
				s.current.StopOnError(err)
				return
			}
			if s.interval > 0 {
				time.Sleep(time.Millisecond * time.Duration(s.interval))
			}
		case <- s.ctx.Done():
			logger.Infof("task cancel")
			return
		}
	}
}

func (s *task) sendTo(customer *schema.Customer) error {
	logger.Debugf("prepare sending email to %s(%s)", customer.Name, customer.Email)
	m := gomail.NewMessage()

	m.SetHeader("From", s.tpl.From)
	m.SetHeader("To", customer.Email)
	m.SetHeader("Subject", s.nextSubject())

	data := map[string]string {
		"CustomerName": customer.Name,
		"Date": time.Now().Format("2006-01-02"),
	}

	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return s.templateEngine.Execute(w, data)
	})

	record := schema.Record {
		TemplateID: s.tpl.ID,
		CustomerID: customer.ID,
		Status:     1,
	}

	if err := gomail.Send(s.sender, m); err != nil {
		logger.Errorf("failure sent email: %s", err)
		record.Reason = err.Error()
		record.Status = 2
		s.current.Failure()
		s.consecutiveSendFailures += 1
	} else {
		logger.Debugf("success sent email")
		s.current.Success()
		s.consecutiveSendFailures = 0
	}

	if s.consecutiveSendFailures > 2 {
		return errors.Errorf("%d次连续发送失败", s.consecutiveSendFailures)
	}

	return s.RecordRepo.Create(s.ctx, record)
}

func (s *task) nextSubject() string {
	length := len(s.subjects)
	if length == 1 {
		return s.subjects[0]
	}
	return s.subjects[rand.Intn(length)]
}
