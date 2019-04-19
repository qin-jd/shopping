package repository

import (
	"github.com/jinzhu/gorm"
	"shopping/order/model"
)

type Repository interface {
	Create(*model.Order) error
	Find(orderId string) (*model.Order, error)
	Update(model.Order, int64) (model.Order, error)
}

type Order struct {
	Db *gorm.DB
}


func (repo *Order) Create(order *model.Order)error{
	if err := repo.Db.Create(order).Error; err != nil {
		return err
	}
	return nil
}


func (repo *Order) Find(orderId string)(*model.Order ,error){
	order :=  &model.Order{}
	order.OrderId = orderId
	if err := repo.Db.First(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}


func (repo *Order) Update(order *model.Order)(*model.Order ,error){
	if err := repo.Db.Model(&order).Updates(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}