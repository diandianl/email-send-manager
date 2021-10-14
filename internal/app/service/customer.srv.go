package service

import (
	"context"
	"email-send-manager/internal/app/dao/customer"
	"github.com/xuri/excelize/v2"
	"io"

	"github.com/google/wire"

	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
)

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
		return err
	}
	rows, err := xlsx.GetRows("Sheet1")

	var customers []*customer.Customer
	for _, row := range rows {
		if len(row) < 2 {
			return errors.Errorf("invalid format excel file")
		}
		name, email := row[0], row[1]
		customers = append(customers, &customer.Customer{Name: name, Email: email})
	}
	return a.CustomerRepo.BatchCreate(ctx, customers)
}
