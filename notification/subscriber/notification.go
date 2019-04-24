package subscriber

import (
	"context"
	"fmt"
	"github.com/micro/go-log"

	notification "shopping/notification/proto/notification"
)

type Notification struct{}

func (e *Notification) Handle(ctx context.Context, req *notification.SubmitRequest) error {
	log.Log(fmt.Sprintf("Handler Received message: ID为%v 的用户购买了商品ID为：%v 的物品" , req.Uid , req.ProductId))
	return nil
}

