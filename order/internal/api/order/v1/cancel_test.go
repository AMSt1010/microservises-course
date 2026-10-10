package v1

import (
	"fmt"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

func (s *ApiSuite) TestOrderCancelByUuidEmptyOrderUuid() {
	orderUUID := ""
	s.service.On("CancelOrderByUuid", s.ctx, orderUUID).Return(model.ErrOrderUUIDCannotBeEmpty).Once()
	resp, err := s.api.OrderCancelByUuid(s.ctx, orderV1.OrderCancelByUuidParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	badReqErr, ok := resp.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusBadRequest, badReqErr.Code)
	s.Require().Equal(model.ErrOrderUUIDCannotBeEmpty.Error(), badReqErr.Message)
}

func (s *ApiSuite) TestOrderCancelByUuidAlreadyPaid() {
	orderUUID := gofakeit.UUID()
	s.service.On("CancelOrderByUuid", s.ctx, orderUUID).Return(model.ErrAlreadyPaid).Once()
	resp, err := s.api.OrderCancelByUuid(s.ctx, orderV1.OrderCancelByUuidParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	conflictError, ok := resp.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusConflict, conflictError.Code)
	s.Require().Equal("The order has already been paid for or cancelled", conflictError.Message)
}

func (s *ApiSuite) TestOrderCancelByUuidOrderCancelled() {
	orderUUID := gofakeit.UUID()
	s.service.On("CancelOrderByUuid", s.ctx, orderUUID).Return(model.ErrOrderCancelled).Once()
	resp, err := s.api.OrderCancelByUuid(s.ctx, orderV1.OrderCancelByUuidParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	conflictError, ok := resp.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusConflict, conflictError.Code)
	s.Require().Equal("The order has already been paid for or cancelled", conflictError.Message)
}

func (s *ApiSuite) TestOrderCancelByUuidOrderNotFound() {
	orderUUID := gofakeit.UUID()
	s.service.On("CancelOrderByUuid", s.ctx, orderUUID).Return(model.ErrOrderNotFound).Once()
	resp, err := s.api.OrderCancelByUuid(s.ctx, orderV1.OrderCancelByUuidParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	notFoundError, ok := resp.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundError.Code)
	s.Require().Equal(fmt.Sprintf("Order by uuid '%s' not found", orderUUID), notFoundError.Message)
}

func (s *ApiSuite) TestOrderCancelByUuidGenericError() {
	orderUUID := gofakeit.UUID()
	s.service.On("CancelOrderByUuid", s.ctx, orderUUID).Return(gofakeit.Error()).Once()
	resp, err := s.api.OrderCancelByUuid(s.ctx, orderV1.OrderCancelByUuidParams{
		OrderUUID: orderUUID,
	})

	s.Require().Error(err)
	s.Require().Nil(resp)

	s.Require().Equal("an unexpected error occurred. Please try again later", err.Error())
}

func (s *ApiSuite) TestOrderCancelByUuidSuccess() {
	orderUUID := gofakeit.UUID()
	s.service.On("CancelOrderByUuid", s.ctx, orderUUID).Return(nil).Once()
	resp, err := s.api.OrderCancelByUuid(s.ctx, orderV1.OrderCancelByUuidParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	_, ok := resp.(*orderV1.OrderCancelByUuidNoContent)
	s.Require().True(ok, "Ожидался ответ типа *orderV1.OrderCancelByUuidNoContent (HTTP 204)")
}
