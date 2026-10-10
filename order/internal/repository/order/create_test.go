package order

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	"github.com/AMSt1010/microservises-course/week2/order/internal/repository/converter"
)

func (s *RepositorySuite) TestCreateOrderSuccess() {
	order := s.GenerateOrder()

	err := s.repo.CreateOrder(s.ctx, order)

	s.Require().NoError(err)

	s.Require().Contains(s.repo.orders, order.OrderUUID)
	s.Require().Equal(s.repo.orders[order.OrderUUID], converter.OrderToRepoModel(order))
}

func (s *RepositorySuite) TestCreateOrderAlreadyExist() {
	order := s.GenerateOrder()

	err := s.repo.CreateOrder(s.ctx, order)

	s.Require().NoError(err)

	err = s.repo.CreateOrder(s.ctx, order)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrAlreadyExists)
}

func (s *RepositorySuite) TestCreateOrderContextDone() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	order := s.GenerateOrder()

	err := s.repo.CreateOrder(ctx, order)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
}
