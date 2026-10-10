package converter

import (
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	repoModel "github.com/AMSt1010/microservises-course/week2/order/internal/repository/model"
)

func OrderToRepoModel(order model.Order) repoModel.Order {
	partUUIDs := make([]string, len(order.PartUuids))
	copy(partUUIDs, order.PartUuids)
	var paymentMethod *repoModel.PaymentMethod

	if order.PaymentMethod != nil {
		paymentMethod = new(repoModel.PaymentMethod(*order.PaymentMethod))
	}
	var orderStatus *repoModel.OrderStatus
	if order.Status != nil {
		orderStatus = new(repoModel.OrderStatus(*order.Status))
	}

	var transactionUUID *string
	if order.TransactionUUID != nil {
		transactionUUID = new(string(*order.TransactionUUID))
	}
	return repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderStatus,
	}
}

func OrderToModel(order repoModel.Order) model.Order {
	partUUIDs := make([]string, len(order.PartUuids))
	copy(partUUIDs, order.PartUuids)
	var paymentMethod *model.PaymentMethod

	if order.PaymentMethod != nil {
		paymentMethod = new(model.PaymentMethod(*order.PaymentMethod))
	}
	var orderStatus *model.OrderStatus
	if order.Status != nil {
		orderStatus = new(model.OrderStatus(*order.Status))
	}
	var transactionUUID *string
	if order.TransactionUUID != nil {
		transactionUUID = new(string(*order.TransactionUUID))
	}
	return model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderStatus,
	}
}
