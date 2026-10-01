package model

// Выбранный способ оплаты
type PaymentMethod int32

const (
	// Неизвестная категория
	PaymentMethodUnspecified PaymentMethod = iota
	// Банковская карта
	PaymentMethodCard
	// Система быстрых платежей
	PaymentMethodSBP
	// Кредитная карта
	PaymentMethodCreditCard
	// Деньги инвестора (внутренний метод)
	PaymentMethodInvestorMoney
)

type Payment struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}

func (p Payment) Validate() error {
	if p.OrderUUID == "" || p.UserUUID == "" || p.PaymentMethod <= PaymentMethodUnspecified || p.PaymentMethod > PaymentMethodInvestorMoney {
		return ErrInvalidArgument
	}
	return nil
}
