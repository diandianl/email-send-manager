package record

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
)

// RecordSet Injection wire
var RecordSet = wire.NewSet(wire.Struct(new(RecordRepo), "*"))

type RecordRepo struct {
	DB *gorm.DB
}

func (a *RecordRepo) getQueryOption(opts ...schema.RecordQueryOptions) schema.RecordQueryOptions {
	var opt schema.RecordQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *RecordRepo) Query(ctx context.Context, params schema.RecordQueryParam, opts ...schema.RecordQueryOptions) (*schema.RecordQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetRecordDB(ctx, a.DB)
	// TODO: 查询条件

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	var list Records
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.RecordQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRecords(),
	}

	return qr, nil
}

func (a *RecordRepo) Get(ctx context.Context, id uint64, opts ...schema.RecordQueryOptions) (*schema.Record, error) {
	var item Record
	ok, err := util.FindOne(ctx, GetRecordDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRecord(), nil
}

func (a *RecordRepo) Create(ctx context.Context, item schema.Record) error {
	eitem := SchemaRecord(item).ToRecord()
	result := GetRecordDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *RecordRepo) Update(ctx context.Context, id uint64, item schema.Record) error {
	eitem := SchemaRecord(item).ToRecord()
	result := GetRecordDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *RecordRepo) Delete(ctx context.Context, id uint64) error {
	result := GetRecordDB(ctx, a.DB).Where("id=?", id).Delete(Record{})
	return errors.WithStack(result.Error)
}

func (a *RecordRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
	result := GetRecordDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
