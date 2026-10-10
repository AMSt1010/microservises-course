package v1

import (
	"fmt"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/AMSt1010/microservises-course/week2/order/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

func (s *ApiSuite) TestOrderPayReqNil() {
	resp, err := s.api.OrderPay(s.ctx, nil, orderV1.OrderPayParams{
		OrderUUID: gofakeit.UUID(),
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	badReqErr, ok := resp.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusBadRequest, badReqErr.Code)
	s.Require().Equal("The request body is invalid or malformed", badReqErr.Message)
}

func (s *ApiSuite) TestOrderPayEmptyOrderUUID() {
	modelPaymentMethod := model.PaymentMethodCARD
	orderUUID := ""
	s.service.On("PayOrder", s.ctx, model.PayOrderIn{
		PaymentMethod: modelPaymentMethod,
		OrderUUID:     orderUUID,
	}).Return("", model.ErrOrderUUIDCannotBeEmpty).Once()
	resp, err := s.api.OrderPay(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: converter.PaymentMethodToAPI(modelPaymentMethod),
	}, orderV1.OrderPayParams{
		OrderUUID: orderUUID,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	badReqErr, ok := resp.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusBadRequest, badReqErr.Code)
	s.Require().Equal(model.ErrOrderUUIDCannotBeEmpty.Error(), badReqErr.Message)
}

func (s *ApiSuite) TestOrderPayInsufficientPaymentMethod() {
	orderUUID := gofakeit.UUID()
	modelPaymentMethod := model.PaymentMethodUNKNOWN

	s.service.On("PayOrder", s.ctx, model.PayOrderIn{
		PaymentMethod: modelPaymentMethod,
		OrderUUID:     orderUUID,
	}).Return("", model.ErrInsufficientPaymentMethod).Once()

	resp, err := s.api.OrderPay(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: converter.PaymentMethodToAPI(modelPaymentMethod),
	}, orderV1.OrderPayParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	badReqErr, ok := resp.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusBadRequest, badReqErr.Code)
	s.Require().Equal(model.ErrInsufficientPaymentMethod.Error(), badReqErr.Message)
}

func (s *ApiSuite) TestOrderPayOrderNotFound() {
	orderUUID := gofakeit.UUID()
	modelPaymentMethod := model.PaymentMethodCARD

	s.service.On("PayOrder", s.ctx, model.PayOrderIn{
		PaymentMethod: modelPaymentMethod,
		OrderUUID:     orderUUID,
	}).Return("", model.ErrOrderNotFound).Once()

	resp, err := s.api.OrderPay(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: converter.PaymentMethodToAPI(modelPaymentMethod),
	}, orderV1.OrderPayParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	notFoundErr, ok := resp.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundErr.Code)
	s.Require().Equal(fmt.Sprintf("Order by uuid '%s' not found", orderUUID), notFoundErr.Message)
}

func (s *ApiSuite) TestOrderPayAlreadyPaid() {
	orderUUID := gofakeit.UUID()
	modelPaymentMethod := model.PaymentMethodCARD

	s.service.On("PayOrder", s.ctx, model.PayOrderIn{
		PaymentMethod: modelPaymentMethod,
		OrderUUID:     orderUUID,
	}).Return("", model.ErrAlreadyPaid).Once()

	resp, err := s.api.OrderPay(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: converter.PaymentMethodToAPI(modelPaymentMethod),
	}, orderV1.OrderPayParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	conflictErr, ok := resp.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusConflict, conflictErr.Code)
	s.Require().Equal("The order has already been paid for or cancelled", conflictErr.Message)
}

func (s *ApiSuite) TestOrderPayOrderCancelled() {
	orderUUID := gofakeit.UUID()
	modelPaymentMethod := model.PaymentMethodCARD

	s.service.On("PayOrder", s.ctx, model.PayOrderIn{
		PaymentMethod: modelPaymentMethod,
		OrderUUID:     orderUUID,
	}).Return("", model.ErrOrderCancelled).Once()

	resp, err := s.api.OrderPay(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: converter.PaymentMethodToAPI(modelPaymentMethod),
	}, orderV1.OrderPayParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	conflictErr, ok := resp.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusConflict, conflictErr.Code)
	s.Require().Equal("The order has already been paid for or cancelled", conflictErr.Message)
}

func (s *ApiSuite) TestOrderPayPaymentUnavailable() {
	orderUUID := gofakeit.UUID()
	modelPaymentMethod := model.PaymentMethodCARD

	s.service.On("PayOrder", s.ctx, model.PayOrderIn{
		PaymentMethod: modelPaymentMethod,
		OrderUUID:     orderUUID,
	}).Return("", model.ErrPaymentUnavailable).Once()

	resp, err := s.api.OrderPay(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: converter.PaymentMethodToAPI(modelPaymentMethod),
	}, orderV1.OrderPayParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	badGatewayError, ok := resp.(*orderV1.BadGatewayError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusBadGateway, badGatewayError.Code)
	s.Require().Equal("Bad Gateway: Invalid response from upstream service", badGatewayError.Message)
}

func (s *ApiSuite) TestOrderPayGeneralError() {
	orderUUID := gofakeit.UUID()
	modelPaymentMethod := model.PaymentMethodCARD
	s.service.On("PayOrder", s.ctx, model.PayOrderIn{
		PaymentMethod: modelPaymentMethod,
		OrderUUID:     orderUUID,
	}).Return("", gofakeit.Error()).Once()

	resp, err := s.api.OrderPay(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: converter.PaymentMethodToAPI(modelPaymentMethod),
	}, orderV1.OrderPayParams{
		OrderUUID: orderUUID,
	})

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Require().Equal("an unexpected error occurred. Please try again later", err.Error())
}

func (s *ApiSuite) TestOrderPaySuccess() {
	orderUUID := gofakeit.UUID()
	transactionUUID := gofakeit.UUID()
	modelPaymentMethod := model.PaymentMethodCARD
	s.service.On("PayOrder", s.ctx, model.PayOrderIn{
		PaymentMethod: modelPaymentMethod,
		OrderUUID:     orderUUID,
	}).Return(transactionUUID, nil).Once()

	resp, err := s.api.OrderPay(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: converter.PaymentMethodToAPI(modelPaymentMethod),
	}, orderV1.OrderPayParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	payOrderResp, ok := resp.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(transactionUUID, payOrderResp.TransactionUUID)
}
