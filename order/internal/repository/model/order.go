package model

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
