package handler

import (
	"context"
	"github.com/micro/go-log"
	"shopping/order/model"
	"shopping/order/repository"
	"github.com/bwmarrin/snowflake"
	order "shopping/order/proto/order"
	product "shopping/product/proto/product"
)

type Order struct{
	Order *repository.Order
	ProductCli product.ProductService
}

func (h *Order) Submit (ctx context.Context , req *order.SubmitRequest, rsp *order.Response) error{
	log.Log("Received Product.Search request")

	node, err := snowflake.NewNode(1)
	if err != nil {
		return err
	}

	// Generate a snowflake ID.
	orderId := node.Generate().String()

	order := &model.Order{
		Status:1,
		OrderId:orderId,
		ProductId:req.ProductId,
	}

	if err = h.Order.Create(order); err != nil{
		return err
	}

	rsp.Code = "200"
	rsp.Msg = "订单提交成功"
	return nil
}


func (h *Order) OrderDetail (ctx context.Context , req *order.OrderDetailRequest, rsp *order.Response) error {
	orderDetail , err := h.Order.Find(req.OrderId)

	if err != nil{
		return err
	}

	productDetail,err := h.ProductCli.Detail(context.TODO() , &product.DetailRequest{Id:orderDetail.ProductId})

	rsp.Code = "200"
	rsp.Msg = "订单详情如下：订单号为："+orderDetail.OrderId+"。购买的产品名字为："+productDetail.Product.Name
	return nil

}