package converter

import (
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

func CreateOrderToModel(orderApi orderV1.CreateOrderRequest) model.CreateOrderIn {
	var UUIDs []string
	if len(orderApi.GetPartUuids()) != 0 {
		UUIDs = make([]string, len(orderApi.GetPartUuids()))
		copy(UUIDs, orderApi.GetPartUuids())
	}
	return model.CreateOrderIn{
		UserUUID:  orderApi.UserUUID,
		PartUuids: UUIDs,
	}
}

func OrderToAPIOrder(order model.Order) orderV1.OrderDto {
	var parts []string
	if len(order.PartUuids) != 0 {
		parts = make([]string, len(order.PartUuids))
		copy(parts, order.PartUuids)
	}

	var transactionUUID orderV1.OptString
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.NewOptString(*order.TransactionUUID)
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.NewOptPaymentMethod(PaymentMethodToAPI(*order.PaymentMethod))
	}

	var orderStatus orderV1.OrderStatus
	if order.Status != nil {
		orderStatus = PaymentStatusToAPI(*order.Status)
	}

	return orderV1.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       parts,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderStatus,
	}
}

func PaymentStatusToAPI(st model.OrderStatus) orderV1.OrderStatus {
	switch st {
	case model.OrderStatusPENDINGPAYMENT:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPAID:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusUNKNOWN
	}
}

func PaymentMethodToAPI(m model.PaymentMethod) orderV1.PaymentMethod {
	switch m {
	case model.PaymentMethodCARD:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCREDITCARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodINVESTORMONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func PaymentMethodToModel(m orderV1.PaymentMethod) model.PaymentMethod {
	switch m {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}
