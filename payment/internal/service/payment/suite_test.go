package payment

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/AMSt1010/microservises-course/week2/payment/internal/model"
)

type ServiceSuite struct {
	ctx context.Context
	suite.Suite
	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = NewService()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func fakeOrderGenerate() *model.Payment {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := []model.PaymentMethod{
		model.PaymentMethodCard,
		model.PaymentMethodSBP,
		model.PaymentMethodCreditCard,
		model.PaymentMethodInvestorMoney,
	}[gofakeit.IntN(4)]
	return &model.Payment{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PaymentMethod: paymentMethod,
	}
}
