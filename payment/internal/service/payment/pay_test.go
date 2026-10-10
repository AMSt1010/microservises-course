package payment

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/payment/internal/model"
)

func (s *ServiceSuite) TestPayServiceSuccess() {
	payment := fakeOrderGenerate()

	transactionUUID, err := s.service.PayOrder(s.ctx, *payment)

	s.Require().NoError(err)
	s.Require().NotZero(transactionUUID)
}

func (s *ServiceSuite) TestPayServiceOrderUUIDIsEmpty() {
	payment := fakeOrderGenerate()
	payment.OrderUUID = ""

	transactionUUID, err := s.service.PayOrder(s.ctx, *payment)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInvalidArgument)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayServiceUserUUIDIsEmpty() {
	payment := fakeOrderGenerate()
	payment.UserUUID = ""

	transactionUUID, err := s.service.PayOrder(s.ctx, *payment)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInvalidArgument)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayServiceIncorrectPaymentMethod() {
	payment := fakeOrderGenerate()
	payment.PaymentMethod = model.PaymentMethodUnspecified

	transactionUUID, err := s.service.PayOrder(s.ctx, *payment)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInvalidArgument)
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayServiceCancelledContext() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	payment := fakeOrderGenerate()

	transactionUUID, err := s.service.PayOrder(ctx, *payment)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(transactionUUID)
}
