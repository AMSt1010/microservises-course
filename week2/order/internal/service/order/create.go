package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

// Проверяет, что все запрашиваемые детали существуют
func areAllPartsExist(checkUuids []string, parts []model.Part) bool {
	getParts := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		getParts[part.UUID] = struct{}{}
	}
	for _, checkuuid := range checkUuids {
		if _, ok := getParts[checkuuid]; !ok {
			return false
		}
	}
	return true
}

// Считает общую цену заказа с учетом количества каждого запрошенного UUID
func totalPrice(requestedUUIDs []string, parts []model.Part) float64 {
	priceMap := make(map[string]float64, len(parts))
	for _, part := range parts {
		priceMap[part.UUID] = part.Info.Price
	}
	var total float64
	for _, id := range requestedUUIDs {
		total += priceMap[id]
	}
	return total
}

// POST /api/v1/orders
// Создаёт новый заказ на основе выбранных пользователем деталей.
func (s *service) CreateOrder(ctx context.Context, order model.CreateOrderIn) (model.Order, error) {
	if err := ctx.Err(); err != nil {
		return model.Order{}, err
	}

	if err := order.Validate(); err != nil {
		return model.Order{}, err
	}

	// генерируем id заказа
	order_uuid := uuid.NewString()

	// список запрашиваемых деталей
	uuids := order.PartUuids

	// Обращаемся к сервису Inventory и получаем список деталей
	filter := model.PartsFilter{
		UUIDs: order.PartUuids,
	}

	listParts, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		return model.Order{}, fmt.Errorf("%w: %w", model.ErrInventoryUnavailable, err)
	}

	// Список деталей
	// Если хотя бы одной детали нет — возвращает ошибку.
	if ok := areAllPartsExist(uuids, listParts); !ok {
		return model.Order{}, model.ErrAtLeastOnePartWasNotFound
	}

	// общая сумма заказа
	total := totalPrice(uuids, listParts)

	// Создаем заказ с полученными данными
	orderInStorage := model.Order{
		OrderUUID:  order_uuid,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: float32(total),
		Status:     new(model.OrderStatusPENDINGPAYMENT),
	}

	// Обновляем данные в хранилище
	err = s.repo.CreateOrder(ctx, orderInStorage)
	if err != nil {
		return model.Order{}, err
	}

	return orderInStorage, nil
}
