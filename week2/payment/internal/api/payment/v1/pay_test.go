package v1

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AMSt1010/microservises-course/week2/payment/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/payment/internal/model"
	PaymentV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/payment/v1"
)

func (s *ApiSuite) TestPayOrderSuccess() {
	payment := fakeOrderGenerate()
	expectedUUID := gofakeit.UUID()
	s.paymentService.On("PayOrder", s.ctx, *payment).Return(expectedUUID, nil)

	transactionUUID, err := s.api.PayOrder(s.ctx, converter.PaymentToProto(*payment))

	s.Require().NoError(err)
	s.Require().Equal(&PaymentV1.PayOrderResponse{
		TransactionUuid: expectedUUID,
	}, transactionUUID)
}

func (s *ApiSuite) TestPayOrderInvalidArgument() {
	payment := fakeOrderGenerate()
	s.paymentService.On("PayOrder", s.ctx, *payment).Return("", model.ErrInvalidArgument)

	transactionUUID, err := s.api.PayOrder(s.ctx, converter.PaymentToProto(*payment))

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
	s.Require().Equal("all payment fields must be filled", st.Message())
	s.Require().Nil(transactionUUID)
}

func (s *ApiSuite) TestPayOrderCancelledContext() {
	payment := fakeOrderGenerate()
	s.paymentService.On("PayOrder", s.ctx, *payment).Return("", context.Canceled)

	transactionUUID, err := s.api.PayOrder(s.ctx, converter.PaymentToProto(*payment))

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Canceled, st.Code())
	s.Require().Nil(transactionUUID)
}

func (s *ApiSuite) TestPayOrderGeneralError() {
	payment := fakeOrderGenerate()
	expectedErr := gofakeit.Error()
	s.paymentService.On("PayOrder", s.ctx, *payment).Return("", expectedErr)

	transactionUUID, err := s.api.PayOrder(s.ctx, converter.PaymentToProto(*payment))

	st, ok := status.FromError(err)
	s.Require().False(ok)
	s.Require().Equal(codes.Unknown, st.Code())
	s.Require().Equal(expectedErr.Error(), st.Message())
	s.Require().Nil(transactionUUID)
}
