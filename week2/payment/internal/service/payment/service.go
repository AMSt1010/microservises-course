package payment

import (
	def "github.com/AMSt1010/microservises-course/week2/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
