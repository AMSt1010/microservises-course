package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/AMSt1010/microservises-course/week2/order/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

func (s *ApiSuite) TestCreateOrderSuccess() {
	userUUID := uuid.NewString()
	partUUIDs := []string{uuid.NewString(), uuid.NewString()}

	req := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	}

	expectedOrder := model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: 250.75,
		Status:     new(model.OrderStatusPENDINGPAYMENT),
	}

	s.service.On("CreateOrder", s.ctx, converter.CreateOrderToModel(*req)).
		Return(expectedOrder, nil).
		Once()

	res, err := s.api.CreateOrder(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedOrder.OrderUUID, res.OrderUUID)
	s.Require().Equal(expectedOrder.TotalPrice, res.TotalPrice)
}

func (s *ApiSuite) TestCreateOrderEmptyUserUUID() {
	partsUUIDS := []string{uuid.NewString(), uuid.NewString()}
	req := orderV1.CreateOrderRequest{
		PartUuids: partsUUIDS,
	}
	s.service.On("CreateOrder", s.ctx, converter.CreateOrderToModel(req)).Return(model.Order{}, model.ErrUserUUIDCannotBeEmpty).
		Once()

	resp, err := s.api.CreateOrder(s.ctx, &req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	var genericErr *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &genericErr)
	s.Require().Equal(http.StatusBadRequest, genericErr.StatusCode)
	s.Require().Equal(model.ErrUserUUIDCannotBeEmpty.Error(), genericErr.Response.Message)
}

func (s *ApiSuite) TestCreateOrderNilRequest() {
	resp, err := s.api.CreateOrder(s.ctx, nil)

	s.Require().Error(err)
	s.Require().Nil(resp)

	var genericErr *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &genericErr)
	s.Require().Equal(http.StatusBadRequest, genericErr.StatusCode)
	s.Require().Equal("request body cannot be nil", genericErr.Response.Message)
}

func (s *ApiSuite) TestCreateOrderEmptyPartsUUIDs() {
	userUUID := uuid.NewString()
	req := orderV1.CreateOrderRequest{
		UserUUID: userUUID,
	}
	s.service.On("CreateOrder", s.ctx, converter.CreateOrderToModel(req)).Return(model.Order{}, model.ErrPartsUUIDsCannotBeEmpty).
		Once()

	resp, err := s.api.CreateOrder(s.ctx, &req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	var genericErr *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &genericErr)
	s.Require().Equal(http.StatusBadRequest, genericErr.StatusCode)
	s.Require().Equal(model.ErrPartsUUIDsCannotBeEmpty.Error(), genericErr.Response.Message)
}

func (s *ApiSuite) TestCreateOrderInventoryUnavailable() {
	userUUID := uuid.NewString()
	partsUUIDS := []string{uuid.NewString(), uuid.NewString()}
	req := orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partsUUIDS,
	}
	s.service.On("CreateOrder", s.ctx, converter.CreateOrderToModel(req)).Return(model.Order{}, model.ErrInventoryUnavailable).
		Once()

	resp, err := s.api.CreateOrder(s.ctx, &req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	var genericErr *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &genericErr)
	s.Require().Equal(http.StatusServiceUnavailable, genericErr.StatusCode)
	s.Require().Equal("Inventory service is temporarily unavailable", genericErr.Response.Message)
}

func (s *ApiSuite) TestCreateOrderAtLeastOnePartWasNotFound() {
	userUUID := uuid.NewString()
	partsUUIDS := []string{uuid.NewString(), uuid.NewString()}
	req := orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partsUUIDS,
	}
	s.service.On("CreateOrder", s.ctx, converter.CreateOrderToModel(req)).Return(model.Order{}, model.ErrAtLeastOnePartWasNotFound).
		Once()

	resp, err := s.api.CreateOrder(s.ctx, &req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	var genericErr *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &genericErr)
	s.Require().Equal(http.StatusBadRequest, genericErr.StatusCode)
	s.Require().Equal(model.ErrAtLeastOnePartWasNotFound.Error(), genericErr.Response.Message)
}

func (s *ApiSuite) TestCreateOrderUnknownError() {
	userUUID := uuid.NewString()
	partsUUIDS := []string{uuid.NewString(), uuid.NewString()}
	randomErr := gofakeit.Error()
	req := orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partsUUIDS,
	}
	s.service.On("CreateOrder", s.ctx, converter.CreateOrderToModel(req)).Return(model.Order{}, randomErr).
		Once()

	resp, err := s.api.CreateOrder(s.ctx, &req)

	s.Require().Error(err)
	s.Require().Nil(resp)

	var genericErr *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &genericErr)
	s.Require().Equal(http.StatusInternalServerError, genericErr.StatusCode)
	s.Require().Equal(randomErr.Error(), genericErr.Response.Message)
}
