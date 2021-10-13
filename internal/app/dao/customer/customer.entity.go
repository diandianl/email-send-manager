package customer

import (
	"context"

	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/util/structure"
)

// GetCustomerDB Get Customer db model
func GetCustomerDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Customer))
}

// SchemaCustomer Customer schema
type SchemaCustomer schema.Customer

// ToCustomer Convert to Customer entity
func (a SchemaCustomer) ToCustomer() *Customer {
	item := new(Customer)
	structure.Copy(a, item)
	return item
}

// Customer Customer entity
type Customer struct {
	util.Model
	Name    string `gorm:"size:50;index;"`                // 名称
	Email   string `gorm:"size:100;index;unique;"`               // 邮箱
	Status  int    `gorm:"type:tinyint;index;default:0;"` // 状态(1:启用 2:停用)
}

// ToSchemaCustomer Convert to Customer schema
func (a Customer) ToSchemaCustomer() *schema.Customer {
	item := new(schema.Customer)
	structure.Copy(a, item)
	return item
}

// Customers Customer entity list
type Customers []*Customer

// ToSchemaCustomers Convert to Customer schema list
func (a Customers) ToSchemaCustomers() []*schema.Customer {
	list := make([]*schema.Customer, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaCustomer()
	}
	return list
}
