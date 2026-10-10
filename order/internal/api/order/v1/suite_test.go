package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/AMSt1010/microservises-course/week2/order/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite
	ctx     context.Context
	service *mocks.OrderService
	api     *api
}

func (s *ApiSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = mocks.NewOrderService(s.T())
	s.api = NewAPI(s.service)
}

func (s *ApiSuite) TearDownTest() {}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
