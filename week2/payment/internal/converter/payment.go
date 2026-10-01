package converter

import (
	"github.com/AMSt1010/microservises-course/week2/payment/internal/model"
	PaymentV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/payment/v1"
)

func PaymentToModel(req *PaymentV1.PayOrderRequest) model.Payment {
	if req == nil {
		return model.Payment{}
	}

	return model.Payment{
		OrderUUID:     req.GetOrderUuid(),
		UserUUID:      req.GetUserUuid(),
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	}
}

func PaymentToProto(p model.Payment) *PaymentV1.PayOrderRequest {
	return &PaymentV1.PayOrderRequest{
		OrderUuid:     p.OrderUUID,
		UserUuid:      p.UserUUID,
		PaymentMethod: PaymentV1.PaymentMethod(p.PaymentMethod),
	}
}
