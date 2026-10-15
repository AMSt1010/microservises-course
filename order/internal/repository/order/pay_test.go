package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	"github.com/AMSt1010/microservises-course/week2/order/internal/repository/converter"
)

func (s *RepositorySuite) TestPayOrderSuccess() {
	orderFake := s.GenerateOrder()
	orderFake.TransactionUUID = nil
	orderFake.Status = lo.ToPtr(model.OrderStatusPENDINGPAYMENT)
	err := s.repo.CreateOrder(s.ctx, orderFake)

	s.Require().NoError(err)
	s.Require().NotNil(orderFake.PaymentMethod)
	order, err := s.repo.PayOrder(s.ctx, orderFake.OrderUUID, gofakeit.UUID(), *orderFake.PaymentMethod)

	s.Require().NoError(err)

	s.Require().NotNil(order.Status)

	s.Require().Equal(*order.Status, model.OrderStatusPAID)
	s.Require().Contains(s.repo.orders, order.OrderUUID)
	s.Require().Equal(order, converter.OrderToModel(s.repo.orders[order.OrderUUID]))
}

func (s *RepositorySuite) TestPayOrderNotFound() {
	orderFake := s.GenerateOrder()
	s.Require().NotNil(orderFake.PaymentMethod)
	s.Require().NotNil(orderFake.TransactionUUID)
	order, err := s.repo.PayOrder(s.ctx, orderFake.OrderUUID, *orderFake.TransactionUUID, *orderFake.PaymentMethod)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)

	s.Require().Empty(order)
}

func (s *RepositorySuite) TestPayOrderContextDone() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	orderFake := s.GenerateOrder()
	s.Require().NotNil(orderFake.PaymentMethod)
	s.Require().NotNil(orderFake.TransactionUUID)
	order, err := s.repo.PayOrder(ctx, orderFake.OrderUUID, *orderFake.TransactionUUID, *orderFake.PaymentMethod)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(order)
}

func (s *RepositorySuite) TestPayOrderStatusPaid() {
	orderFake := s.GenerateOrder()
	orderFake.Status = lo.ToPtr(model.OrderStatusPAID)
	err := s.repo.CreateOrder(s.ctx, orderFake)

	s.Require().NoError(err)
	s.Require().NotNil(orderFake.PaymentMethod)
	s.Require().NotNil(orderFake.TransactionUUID)
	order, err := s.repo.PayOrder(s.ctx, orderFake.OrderUUID, *orderFake.TransactionUUID, *orderFake.PaymentMethod)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrAlreadyPaid)

	s.Require().Empty(order)
}

func (s *RepositorySuite) TestPayOrderStatusCancelled() {
	orderFake := s.GenerateOrder()
	orderFake.Status = lo.ToPtr(model.OrderStatusCANCELLED)
	err := s.repo.CreateOrder(s.ctx, orderFake)

	s.Require().NoError(err)
	s.Require().NotNil(orderFake.PaymentMethod)
	s.Require().NotNil(orderFake.TransactionUUID)
	order, err := s.repo.PayOrder(s.ctx, orderFake.OrderUUID, *orderFake.TransactionUUID, *orderFake.PaymentMethod)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderCancelled)

	s.Require().Empty(order)
}
