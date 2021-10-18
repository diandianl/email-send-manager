package service

import (
	"context"
	"email-send-manager/internal/app/dao/customer"
	"github.com/xuri/excelize/v2"
	"io"
	"regexp"
	"strconv"

	"github.com/google/wire"

	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
)

var emailRegExp *regexp.Regexp

func init() {
	emailRegExp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
}

var CustomerSet = wire.NewSet(wire.Struct(new(CustomerSrv), "*"))

type CustomerSrv struct {
	TransRepo    *dao.TransRepo
	CustomerRepo *dao.CustomerRepo
}

func (a *CustomerSrv) Query(ctx context.Context, params schema.CustomerQueryParam, opts ...schema.CustomerQueryOptions) (*schema.CustomerQueryResult, error) {
	result, err := a.CustomerRepo.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *CustomerSrv) Get(ctx context.Context, id uint, opts ...schema.CustomerQueryOptions) (*schema.Customer, error) {
	item, err := a.CustomerRepo.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *CustomerSrv) Create(ctx context.Context, item schema.Customer) (*schema.IDResult, error) {

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.CustomerRepo.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *CustomerSrv) Update(ctx context.Context, id uint, item schema.Customer) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	item.ID = oldItem.ID
	item.CreatedAt = oldItem.CreatedAt

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.CustomerRepo.Update(ctx, id, item)
	})
}

func (a *CustomerSrv) Delete(ctx context.Context, id uint) error {
	oldItem, err := a.CustomerRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.CustomerRepo.Delete(ctx, id)
	})
}

func (a *CustomerSrv) UpdateStatus(ctx context.Context, id uint, status int) error {
	oldItem, err := a.CustomerRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Status == status {
		return nil
	}

	return a.CustomerRepo.UpdateStatus(ctx, id, status)
}

func (a *CustomerSrv) Import(ctx context.Context, r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return errors.Wrap400Response(err, "无法打开文件")
	}
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return errors.Wrap400Response(err, "无法打开Sheet1")
	}

	var customers []*customer.Customer
	for i, row := range rows {
		if len(row) < 1 {
			return errors.New400Response("无效的Excel，需要满足行格式 [email | name | status]（name, status 可选）")
		}

		if !emailRegExp.MatchString(row[0]) {
			return errors.New400Response("第%d行无效的邮箱格式，'%s'", i + 1, row[0])
		}

		c := &customer.Customer{Email: row[0]}
		if len(row) > 1 {
			c.Name = row[1]
		}
		if len(row) > 2 {
			c.Status, _ = strconv.Atoi(row[2])
		}
		if c.Status == 0 {
			c.Status = 1
		}
		customers = append(customers, c)
	}
	return a.CustomerRepo.BatchCreate(ctx, customers)
}
