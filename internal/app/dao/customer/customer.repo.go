package customer

import (
	"context"
	"gorm.io/gorm/clause"

	"github.com/google/wire"
	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
)

// CustomerSet Injection wire
var CustomerSet = wire.NewSet(wire.Struct(new(CustomerRepo), "*"))

type CustomerRepo struct {
	DB *gorm.DB
}

func (a *CustomerRepo) getQueryOption(opts ...schema.CustomerQueryOptions) schema.CustomerQueryOptions {
	var opt schema.CustomerQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *CustomerRepo) Query(ctx context.Context, params schema.CustomerQueryParam, opts ...schema.CustomerQueryOptions) (*schema.CustomerQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetCustomerDB(ctx, a.DB)
	if len(params.Keyword) > 0 {
		db.Where("email LIKE %?%", params.Keyword)
	}
	if len(params.IDs) > 0 {
		if params.Include {
			db.Where("id in ?", params.IDs)
		} else {
			db.Where("id not in ?", params.IDs)
		}
	}

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	var list Customers
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.CustomerQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaCustomers(),
	}

	return qr, nil
}

func (a *CustomerRepo) Get(ctx context.Context, id uint, opts ...schema.CustomerQueryOptions) (*schema.Customer, error) {
	var item Customer
	ok, err := util.FindOne(ctx, GetCustomerDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaCustomer(), nil
}

func (a *CustomerRepo) Create(ctx context.Context, item schema.Customer) error {
	eitem := SchemaCustomer(item).ToCustomer()
	result := GetCustomerDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *CustomerRepo) Update(ctx context.Context, id uint, item schema.Customer) error {
	eitem := SchemaCustomer(item).ToCustomer()
	result := GetCustomerDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *CustomerRepo) Delete(ctx context.Context, id uint) error {
	result := GetCustomerDB(ctx, a.DB).Where("id=?", id).Delete(Customer{})
	return errors.WithStack(result.Error)
}

func (a *CustomerRepo) UpdateStatus(ctx context.Context, id uint, status int) error {
	result := GetCustomerDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}

func (a *CustomerRepo) BatchCreate(ctx context.Context, customers []*Customer) error {
	result := GetCustomerDB(ctx, a.DB).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(customers)
	return errors.WithStack(result.Error)
}
