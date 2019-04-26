package handler

import (
	"context"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"shopping/order/model"
	"shopping/order/repository"
	"github.com/bwmarrin/snowflake"
	order "shopping/order/proto/order"
	product "shopping/product/proto/product"
)

type Order struct{
	Order *repository.Order
	ProductCli product.ProductService
	Publisher micro.Publisher
}

func (h *Order) Submit (ctx context.Context , req *order.SubmitRequest, rsp *order.Response) error{
	log.Log("Received Order.Submit request")

	//查询商品的库存数量
	productDetail,err := h.ProductCli.Detail(ctx , &product.DetailRequest{Id:req.ProductId})
	if productDetail.Product.Number < 1 {
		return errors.BadRequest("go.micro.srv.order" , "库存不足")
	}

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
		Uid : req.Uid,
	}

	if err = h.Order.Create(order); err != nil{
		return err
	}

	//减库存
	reduce,err := h.ProductCli.ReduceNumber(ctx , &product.ReduceNumberRequest{Id:req.ProductId})
	if reduce == nil || reduce.Code != "200" {
		return errors.BadRequest("go.micro.srv.order" , err.Error())
	}

	//异步发送通知给用户订单信息
	if err := h.Publisher.Publish(ctx , req);err != nil {
		return errors.BadRequest("notification" , err.Error())
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