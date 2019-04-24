package model

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	OrderId string
	Status int //订单状态
	ProductId uint32 //商品ID
	Uid uint32 //用户ID
}
