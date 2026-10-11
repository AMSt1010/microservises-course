package v1

import (
	"fmt"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/AMSt1010/microservises-course/week2/order/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

func (s *ApiSuite) TestGetOrderByUuidEmptyOrderUUID() {
	orderUUID := ""
	s.service.On("GetOrderByUuid", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderUUIDCannotBeEmpty)
	resp, err := s.api.GetOrderByUuid(s.ctx, orderV1.GetOrderByUuidParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	badReq, ok := resp.(*orderV1.BadRequestError)

	s.Require().True(ok)
	s.Require().Equal(http.StatusBadRequest, badReq.Code)
	s.Require().Equal(model.ErrOrderUUIDCannotBeEmpty.Error(), badReq.Message)
}

func (s *ApiSuite) TestGetOrderByUuidOrderNotFound() {
	orderUUID := uuid.NewString()
	s.service.On("GetOrderByUuid", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound)
	resp, err := s.api.GetOrderByUuid(s.ctx, orderV1.GetOrderByUuidParams{
		OrderUUID: orderUUID,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	notFound, ok := resp.(*orderV1.NotFoundError)

	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFound.Code)
	s.Require().Equal(fmt.Sprintf("Order by uuid '%s' not found", orderUUID), notFound.Message)
}

func (s *ApiSuite) TestGetOrderByUuidGenericError() {
	orderUUID := uuid.NewString()
	randomErr := gofakeit.Error()
	s.service.On("GetOrderByUuid", s.ctx, orderUUID).Return(model.Order{}, randomErr)
	resp, err := s.api.GetOrderByUuid(s.ctx, orderV1.GetOrderByUuidParams{
		OrderUUID: orderUUID,
	})
	s.Require().Error(err)
	s.Require().Nil(resp)

	var genErr *orderV1.GenericErrorStatusCode

	s.Require().ErrorAs(err, &genErr)
	s.Require().Equal(http.StatusInternalServerError, genErr.StatusCode)
	s.Require().Equal(randomErr.Error(), genErr.Response.Message)
}

func (s *ApiSuite) TestGetOrderByUuidSuccess() {
	orderUUID := uuid.NewString()
	expectedOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.NewString(),
		PartUuids:  []string{uuid.NewString(), uuid.NewString()},
		TotalPrice: 250.75,
		Status:     lo.ToPtr(model.OrderStatusPENDINGPAYMENT),
	}
	s.service.On("GetOrderByUuid", s.ctx, orderUUID).Return(expectedOrder, nil)
	resp, err := s.api.GetOrderByUuid(s.ctx, orderV1.GetOrderByUuidParams{
		OrderUUID: orderUUID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	orderResp, ok := resp.(*orderV1.GetOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(converter.OrderToAPIOrder(expectedOrder), orderResp.Order)
}
