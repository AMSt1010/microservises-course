package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx        context.Context
	invService *mocks.InventoryService
	api        *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()
	s.invService = mocks.NewInventoryService(s.T())
	s.api = NewApi(s.invService)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
