package send_batch

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
)

// SendBatchSet Injection wire
var SendBatchSet = wire.NewSet(wire.Struct(new(SendBatchRepo), "*"))

type SendBatchRepo struct {
	DB *gorm.DB
}

func (a *SendBatchRepo) getQueryOption(opts ...schema.SendBatchQueryOptions) schema.SendBatchQueryOptions {
	var opt schema.SendBatchQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *SendBatchRepo) Query(ctx context.Context, params schema.SendBatchQueryParam, opts ...schema.SendBatchQueryOptions) (*schema.SendBatchQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetSendBatchDB(ctx, a.DB)
	// TODO: 查询条件

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	var list SendBatches
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.SendBatchQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaSendBatches(),
	}

	return qr, nil
}

func (a *SendBatchRepo) Get(ctx context.Context, id uint64, opts ...schema.SendBatchQueryOptions) (*schema.SendBatch, error) {
	var item SendBatch
	ok, err := util.FindOne(ctx, GetSendBatchDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaSendBatch(), nil
}

func (a *SendBatchRepo) Create(ctx context.Context, item schema.SendBatch) error {
	eitem := SchemaSendBatch(item).ToSendBatch()
	result := GetSendBatchDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *SendBatchRepo) Update(ctx context.Context, id uint64, item schema.SendBatch) error {
	eitem := SchemaSendBatch(item).ToSendBatch()
	result := GetSendBatchDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *SendBatchRepo) Delete(ctx context.Context, id uint64) error {
	result := GetSendBatchDB(ctx, a.DB).Where("id=?", id).Delete(SendBatch{})
	return errors.WithStack(result.Error)
}

func (a *SendBatchRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
	result := GetSendBatchDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
