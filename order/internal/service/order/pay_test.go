package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderContextDone() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()
	order, err := s.service.PayOrder(ctx, model.PayOrderIn{})

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestPayOrderEmptyOrderUUID() {
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodCARD,
		OrderUUID:     "",
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderUUIDCannotBeEmpty)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderInsufficientPaymentMethod() {
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodUNKNOWN,
		OrderUUID:     gofakeit.UUID(),
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInsufficientPaymentMethod)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderRepoGetError() {
	orderUUID := gofakeit.UUID()
	randomError := gofakeit.Error()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, randomError)
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodCARD,
		OrderUUID:     orderUUID,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, randomError)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderOrderStatusPaid() {
	orderUUID := gofakeit.UUID()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		Status:    lo.ToPtr(model.OrderStatusPAID),
	}, nil)
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodCARD,
		OrderUUID:     orderUUID,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrAlreadyPaid)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderOrderStatusCancelled() {
	orderUUID := gofakeit.UUID()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		Status:    lo.ToPtr(model.OrderStatusCANCELLED),
	}, nil)
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodCARD,
		OrderUUID:     orderUUID,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderCancelled)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderPaymentClientUnavailable() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    lo.ToPtr(model.OrderStatusPENDINGPAYMENT),
	}, nil)
	s.paymentClient.On("PayOrder", s.ctx, model.Payment{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}).Return("", gofakeit.Error())
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodCARD,
		OrderUUID:     orderUUID,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPaymentUnavailable)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderPaymentClientContextDone() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	s.repo.On("GetOrder", ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    lo.ToPtr(model.OrderStatusPENDINGPAYMENT),
	}, nil)
	s.paymentClient.On("PayOrder", ctx, model.Payment{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}).Run(func(args mock.Arguments) {
		cancel() // отменяем ctx прямо во время обращения к внешнему сервису
	}).Return("", context.Canceled)
	transactionUUID, err := s.service.PayOrder(ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodCARD,
		OrderUUID:     orderUUID,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderRepoPayOrderError() {
	orderUUID := gofakeit.UUID()
	transactionUUIDExp := gofakeit.UUID()
	randomError := gofakeit.Error()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		Status:    lo.ToPtr(model.OrderStatusPENDINGPAYMENT),
	}, nil)
	s.paymentClient.On("PayOrder", s.ctx, model.Payment{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}).Return(transactionUUIDExp, nil)
	s.repo.On("PayOrder", s.ctx, orderUUID, transactionUUIDExp, model.PaymentMethodCARD).Return(model.Order{}, randomError)
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodCARD,
		OrderUUID:     orderUUID,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, randomError)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderSuccess() {
	orderUUID := gofakeit.UUID()
	transactionUUIDExp := gofakeit.UUID()
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		Status:    lo.ToPtr(model.OrderStatusPENDINGPAYMENT),
	}, nil)
	s.paymentClient.On("PayOrder", s.ctx, model.Payment{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}).Return(transactionUUIDExp, nil)
	s.repo.On("PayOrder", s.ctx, orderUUID, transactionUUIDExp, model.PaymentMethodCARD).Return(model.Order{}, nil)
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PayOrderIn{
		PaymentMethod: model.PaymentMethodCARD,
		OrderUUID:     orderUUID,
	})

	s.Require().NoError(err)
	s.Require().Equal(transactionUUIDExp, transactionUUID)
}
