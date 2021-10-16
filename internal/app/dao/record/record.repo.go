package record

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
	if params.Status > 0 {
		db.Where("status = ?", params.Status)
	}
	if params.TemplateID > 0 {
		db.Where("template_id = ?", params.TemplateID)
	}
	if len(params.Email) > 0 {
		db.Where("customer_id in (SELECT id FROM tb_customer WHERE email LIKE ?)", fmt.Sprintf("%%%s%%", params.Email))
	}

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	db = db.Preload("Template", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("Customer", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "email")
	})

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

func (a *RecordRepo) Get(ctx context.Context, id uint, opts ...schema.RecordQueryOptions) (*schema.Record, error) {
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
	result := GetRecordDB(ctx, a.DB).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "template_id"}, {Name: "customer_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "reason"}),
	}).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *RecordRepo) Update(ctx context.Context, id uint, item schema.Record) error {
	eitem := SchemaRecord(item).ToRecord()
	result := GetRecordDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *RecordRepo) Delete(ctx context.Context, id uint) error {
	result := GetRecordDB(ctx, a.DB).Where("id=?", id).Delete(Record{})
	return errors.WithStack(result.Error)
}

func (a *RecordRepo) DeleteByTemplateId(ctx context.Context, templateId uint) error {
	result := GetRecordDB(ctx, a.DB).Where("template_id=?", templateId).Delete(Record{})
	return errors.WithStack(result.Error)
}

func (a *RecordRepo) UpdateStatus(ctx context.Context, id uint, status int) error {
	result := GetRecordDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
