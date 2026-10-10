package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderByUuidContextDone() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()
	order, err := s.service.GetOrderByUuid(ctx, gofakeit.UUID())

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestGetOrderByUuidEmptyOrderUUID() {
	order, err := s.service.GetOrderByUuid(s.ctx, "")

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderUUIDCannotBeEmpty)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestGetOrderByUuidRepoError() {
	orderUUID := gofakeit.UUID()
	randomError := gofakeit.Error()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, randomError)
	order, err := s.service.GetOrderByUuid(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, randomError)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestGetOrderByUuidSuccess() {
	orderUUID := gofakeit.UUID()

	orderExp := model.Order{
		OrderUUID: orderUUID,
	}
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(orderExp, nil)
	order, err := s.service.GetOrderByUuid(s.ctx, orderUUID)

	s.Require().NoError(err)
	s.Require().Equal(orderExp, order)
}
