package model

import "slices"

type Order struct {
	// Uuid заказа.
	OrderUUID string
	// Uuid пользователя.
	UserUUID string
	// Uuids деталей заказа.
	PartUuids []string
	// Общая цена.
	TotalPrice float32
	// Uuid транзакции.
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          *OrderStatus
}

type CreateOrderIn struct {
	// Uuid пользователя.
	UserUUID string
	// Uuids деталей.
	PartUuids []string
}

func (order CreateOrderIn) Validate() error {
	if order.UserUUID == "" {
		return ErrUserUUIDCannotBeEmpty
	}

	if len(order.PartUuids) == 0 {
		return ErrPartsUUIDsCannotBeEmpty
	}
	return nil
}

type PayOrderIn struct {
	PaymentMethod PaymentMethod
	OrderUUID     string
}

func (payOrder PayOrderIn) Validate() error {
	if payOrder.OrderUUID == "" {
		return ErrOrderUUIDCannotBeEmpty
	}
	if !slices.Contains([]PaymentMethod{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}, payOrder.PaymentMethod) {
		return ErrInsufficientPaymentMethod
	}
	return nil
}

type Payment struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}

func (p Payment) Validate() error {
	if p.OrderUUID == "" || p.UserUUID == "" || !slices.Contains([]PaymentMethod{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}, p.PaymentMethod) {
		return ErrInsufficientPaymentMethod
	}
	return nil
}

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusUNKNOWN        OrderStatus = "UNKNOWN"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)
