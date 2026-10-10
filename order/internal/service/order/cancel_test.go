package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

func (s *ServiceSuite) TestCancelOrderByUuidContextDone() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()
	orderUUID := gofakeit.UUID()
	err := s.service.CancelOrderByUuid(ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
}

func (s *ServiceSuite) TestCancelOrderByUuidEmptyOrderUUID() {
	err := s.service.CancelOrderByUuid(s.ctx, "")

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderUUIDCannotBeEmpty)
}

func (s *ServiceSuite) TestCancelOrderByUuidRepoGetError() {
	orderUUID := gofakeit.UUID()
	randomErr := gofakeit.Error()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, randomErr)
	err := s.service.CancelOrderByUuid(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, randomErr)
}

func (s *ServiceSuite) TestCancelOrderByUuidOrderStatusPaid() {
	orderUUID := gofakeit.UUID()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		Status:    new(model.OrderStatusPAID),
	}, nil)
	err := s.service.CancelOrderByUuid(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrAlreadyPaid)
}

func (s *ServiceSuite) TestCancelOrderByUuidOrderStatusCancelled() {
	orderUUID := gofakeit.UUID()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		Status:    new(model.OrderStatusCANCELLED),
	}, nil)
	err := s.service.CancelOrderByUuid(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}

func (s *ServiceSuite) TestCancelOrderByUuidRepoCancelError() {
	orderUUID := gofakeit.UUID()
	randomError := gofakeit.Error()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		Status:    new(model.OrderStatusPENDINGPAYMENT),
	}, nil)
	s.repo.On("CancelOrder", s.ctx, orderUUID).Return(model.Order{}, randomError)
	err := s.service.CancelOrderByUuid(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, randomError)
}

func (s *ServiceSuite) TestCancelOrderByUuidSuccess() {
	orderUUID := gofakeit.UUID()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		Status:    new(model.OrderStatusPENDINGPAYMENT),
	}, nil)
	s.repo.On("CancelOrder", s.ctx, orderUUID).Return(model.Order{}, nil)
	err := s.service.CancelOrderByUuid(s.ctx, orderUUID)

	s.Require().NoError(err)
}
