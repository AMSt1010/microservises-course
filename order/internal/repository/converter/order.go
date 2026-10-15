package converter

import (
	"github.com/samber/lo"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	repoModel "github.com/AMSt1010/microservises-course/week2/order/internal/repository/model"
)

func OrderToRepoModel(order model.Order) repoModel.Order {
	partUUIDs := make([]string, len(order.PartUuids))
	copy(partUUIDs, order.PartUuids)
	var paymentMethod *repoModel.PaymentMethod

	if order.PaymentMethod != nil {
		paymentMethod = lo.ToPtr(repoModel.PaymentMethod(*order.PaymentMethod))
	}
	var orderStatus *repoModel.OrderStatus
	if order.Status != nil {
		orderStatus = lo.ToPtr(repoModel.OrderStatus(*order.Status))
	}

	var transactionUUID *string
	if order.TransactionUUID != nil {
		transactionUUID = lo.ToPtr(string(*order.TransactionUUID))
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
		paymentMethod = lo.ToPtr(model.PaymentMethod(*order.PaymentMethod))
	}
	var orderStatus *model.OrderStatus
	if order.Status != nil {
		orderStatus = lo.ToPtr(model.OrderStatus(*order.Status))
	}
	var transactionUUID *string
	if order.TransactionUUID != nil {
		transactionUUID = lo.ToPtr(string(*order.TransactionUUID))
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
