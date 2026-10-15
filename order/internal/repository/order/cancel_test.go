package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	"github.com/AMSt1010/microservises-course/week2/order/internal/repository/converter"
)

func (s *RepositorySuite) TestCancelOrderSuccess() {
	orderFake := s.GenerateOrder()
	orderFake.Status = lo.ToPtr(model.OrderStatusPENDINGPAYMENT)
	err := s.repo.CreateOrder(s.ctx, orderFake)

	s.Require().NoError(err)
	order, err := s.repo.CancelOrder(s.ctx, orderFake.OrderUUID)

	s.Require().NoError(err)

	s.Require().NotNil(order.Status)

	s.Require().Equal(*order.Status, model.OrderStatusCANCELLED)
	s.Require().Contains(s.repo.orders, order.OrderUUID)
	s.Require().Equal(order, converter.OrderToModel(s.repo.orders[order.OrderUUID]))
}

func (s *RepositorySuite) TestCancelOrderNotFound() {
	orderFake := s.GenerateOrder()
	order, err := s.repo.CancelOrder(s.ctx, orderFake.OrderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)

	s.Require().Empty(order)
}

func (s *RepositorySuite) TestCancelOrderContextDone() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	orderFake := s.GenerateOrder()
	order, err := s.repo.CancelOrder(ctx, orderFake.OrderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(order)
}

func (s *RepositorySuite) TestCancelOrderStatusPaid() {
	orderFake := s.GenerateOrder()
	orderFake.Status = lo.ToPtr(model.OrderStatusPAID)
	err := s.repo.CreateOrder(s.ctx, orderFake)

	s.Require().NoError(err)
	order, err := s.repo.CancelOrder(s.ctx, orderFake.OrderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrAlreadyPaid)

	s.Require().Empty(order)
}

func (s *RepositorySuite) TestCancelOrderStatusCancelled() {
	orderFake := s.GenerateOrder()
	orderFake.Status = lo.ToPtr(model.OrderStatusCANCELLED)
	err := s.repo.CreateOrder(s.ctx, orderFake)

	s.Require().NoError(err)
	order, err := s.repo.CancelOrder(s.ctx, orderFake.OrderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderCancelled)

	s.Require().Empty(order)
}
