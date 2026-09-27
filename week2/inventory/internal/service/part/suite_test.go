package part

import (
	"context"
	"testing"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	invRepository *mocks.InventoryRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.invRepository = mocks.NewInventoryRepository(s.T())

	s.service = NewService(s.invRepository)
}

func (s *ServiceSuite) TearDownTest() {

}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
