package converter

import (
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	paymentV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(m model.PaymentMethod) paymentV1.PaymentMethod {
	switch m {
	case model.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func PaymentToProto(p model.Payment) *paymentV1.PayOrderRequest {
	return &paymentV1.PayOrderRequest{
		OrderUuid:     p.OrderUUID,
		UserUuid:      p.UserUUID,
		PaymentMethod: PaymentMethodToProto(p.PaymentMethod),
	}
}
