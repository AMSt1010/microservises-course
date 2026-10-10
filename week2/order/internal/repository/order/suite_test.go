package order

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

type RepositorySuite struct {
	suite.Suite
	ctx  context.Context
	repo *repository
}

func (r *RepositorySuite) SetupTest() {
	r.ctx = context.Background()
	r.repo = NewRepository()
}

func (r *RepositorySuite) TearDownTest() {
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (r *RepositorySuite) GenerateOrder() model.Order {
	parts := make([]string, 5)
	for i := range parts {
		parts[i] = gofakeit.UUID()
	}

	payment := []model.PaymentMethod{model.PaymentMethodUNKNOWN, model.PaymentMethodCARD, model.PaymentMethodSBP, model.PaymentMethodCREDITCARD, model.PaymentMethodINVESTORMONEY}[gofakeit.IntN(5)]

	status := []model.OrderStatus{model.OrderStatusUNKNOWN, model.OrderStatusPENDINGPAYMENT, model.OrderStatusPAID, model.OrderStatusCANCELLED}[gofakeit.IntN(4)]

	order := model.Order{
		OrderUUID:       gofakeit.UUID(),
		UserUUID:        gofakeit.UUID(),
		PartUuids:       parts,
		TotalPrice:      gofakeit.Float32(),
		TransactionUUID: new(gofakeit.UUID()),
		PaymentMethod:   &payment,
		Status:          &status,
	}

	return order
}
