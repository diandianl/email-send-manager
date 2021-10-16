package template

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
)

// TemplateSet Injection wire
var TemplateSet = wire.NewSet(wire.Struct(new(TemplateRepo), "*"))

type TemplateRepo struct {
	DB *gorm.DB
}

func (a *TemplateRepo) getQueryOption(opts ...schema.TemplateQueryOptions) schema.TemplateQueryOptions {
	var opt schema.TemplateQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *TemplateRepo) Query(ctx context.Context, params schema.TemplateQueryParam, opts ...schema.TemplateQueryOptions) (*schema.TemplateQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetTemplateDB(ctx, a.DB)
	if params.Status > 0 {
		db.Where("status = ?", params.Status)
	}
	if len(params.Name) > 0 {
		db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", params.Name))
	}

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	var list Templates
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.TemplateQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaTemplates(),
	}

	return qr, nil
}

func (a *TemplateRepo) Get(ctx context.Context, id uint, opts ...schema.TemplateQueryOptions) (*schema.Template, error) {
	var item Template
	ok, err := util.FindOne(ctx, GetTemplateDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaTemplate(), nil
}

func (a *TemplateRepo) Create(ctx context.Context, item schema.Template) error {
	eitem := SchemaTemplate(item).ToTemplate()
	result := GetTemplateDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *TemplateRepo) Update(ctx context.Context, id uint, item schema.Template) error {
	eitem := SchemaTemplate(item).ToTemplate()
	result := GetTemplateDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *TemplateRepo) Delete(ctx context.Context, id uint) error {
	result := GetTemplateDB(ctx, a.DB).Where("id=?", id).Delete(Template{})
	return errors.WithStack(result.Error)
}

func (a *TemplateRepo) UpdateStatus(ctx context.Context, id uint, status int) error {
	result := GetTemplateDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
