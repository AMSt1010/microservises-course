package order

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	clientsMocks "github.com/AMSt1010/microservises-course/week2/order/internal/client/grpc/mocks"
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	"github.com/AMSt1010/microservises-course/week2/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	service *service

	paymentClient *clientsMocks.PaymentClient

	inventoryClient *clientsMocks.InventoryClient

	repo *mocks.OrderRepository
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewOrderRepository(s.T())
	s.paymentClient = clientsMocks.NewPaymentClient(s.T())
	s.inventoryClient = clientsMocks.NewInventoryClient(s.T())
	s.service = NewService(s.repo, s.inventoryClient, s.paymentClient)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) GenerateCreateOrderIn() model.CreateOrderIn {
	num := gofakeit.IntRange(1, 5)
	parts := make([]string, num)

	for i := range parts {
		parts[i] = uuid.NewString()
	}
	createOrderIn := model.CreateOrderIn{
		UserUUID:  gofakeit.UUID(),
		PartUuids: parts,
	}
	return createOrderIn
}

func (s *ServiceSuite) GenerateParts(partsUUIDs []string) []model.Part {
	parts := make([]model.Part, len(partsUUIDs))

	for i := range parts {
		partInfo := model.PartInfo{
			Price: gofakeit.Float64Range(0.0, 100.0),
		}
		part := model.Part{
			UUID: partsUUIDs[i],
			Info: partInfo,
		}
		parts[i] = part
	}
	return parts
}

func (s *ServiceSuite) TotalPrice(requestedUUIDs []string, parts []model.Part) float32 {
	priceMap := make(map[string]float64, len(parts))
	for _, part := range parts {
		priceMap[part.UUID] = part.Info.Price
	}
	var total float64
	for _, id := range requestedUUIDs {
		total += priceMap[id]
	}
	return float32(total)
}
