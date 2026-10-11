package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	"github.com/AMSt1010/microservises-course/week2/order/internal/repository/converter"
	repoModel "github.com/AMSt1010/microservises-course/week2/order/internal/repository/model"
)

// Помечает заказ как оплаченный с использованием сравнения и замены (Compare-And-Swap)
func (r *repository) PayOrder(ctx context.Context, orderUUID, transactionUUID string, method model.PaymentMethod) (model.Order, error) {
	if err := ctx.Err(); err != nil {
		return model.Order{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[orderUUID]

	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	if order.Status != nil && *order.Status == repoModel.OrderStatusPAID {
		return model.Order{}, model.ErrAlreadyPaid
	}
	if order.Status != nil && *order.Status == repoModel.OrderStatusCANCELLED {
		return model.Order{}, model.ErrOrderCancelled
	}

	order.Status = lo.ToPtr(repoModel.OrderStatusPAID)
	order.TransactionUUID = lo.ToPtr(transactionUUID)
	order.PaymentMethod = lo.ToPtr(repoModel.PaymentMethod(method))

	r.orders[orderUUID] = order

	result := converter.OrderToModel(order)

	return result, nil
}
