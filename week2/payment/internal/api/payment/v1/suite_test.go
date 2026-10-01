package v1

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/AMSt1010/microservises-course/week2/payment/internal/model"
	"github.com/AMSt1010/microservises-course/week2/payment/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite
	ctx            context.Context
	api            *api
	paymentService *mocks.PaymentService
}

func (a *ApiSuite) SetupTest() {
	a.ctx = context.Background()
	a.paymentService = mocks.NewPaymentService(a.T())
	a.api = NewAPI(a.paymentService)
}

func (a *ApiSuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
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
