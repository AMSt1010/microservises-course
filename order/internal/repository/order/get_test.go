package order

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

func (s *RepositorySuite) TestGetOrderSuccess() {
	orderExp := s.GenerateOrder()

	err := s.repo.CreateOrder(s.ctx, orderExp)

	s.Require().NoError(err)

	getOrder, err := s.repo.GetOrder(s.ctx, orderExp.OrderUUID)

	s.Require().NoError(err)

	s.Require().Equal(orderExp, getOrder)
}

func (s *RepositorySuite) TestGetOrderNotFound() {
	orderExp := s.GenerateOrder()

	getOrder, err := s.repo.GetOrder(s.ctx, orderExp.OrderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Empty(getOrder)
}

func (s *RepositorySuite) TestGetOrderContextDone() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	order := s.GenerateOrder()

	getOrder, err := s.repo.GetOrder(ctx, order.OrderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(getOrder)
}
