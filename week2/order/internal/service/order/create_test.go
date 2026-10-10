package order

import (
	"context"
	"slices"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderContextDone() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()
	createOrderIn := s.GenerateCreateOrderIn()
	order, err := s.service.CreateOrder(ctx, createOrderIn)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestCreateOrderEmptyUserUUID() {
	createOrderIn := s.GenerateCreateOrderIn()
	createOrderIn.UserUUID = ""
	order, err := s.service.CreateOrder(s.ctx, createOrderIn)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUserUUIDCannotBeEmpty)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestCreateOrderEmptyPartsUUIDs() {
	createOrderIn := s.GenerateCreateOrderIn()
	createOrderIn.PartUuids = nil
	order, err := s.service.CreateOrder(s.ctx, createOrderIn)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartsUUIDsCannotBeEmpty)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestCreateOrderInventoryUnavailable() {
	createOrderIn := s.GenerateCreateOrderIn()
	filter := model.PartsFilter{
		UUIDs: createOrderIn.PartUuids,
	}
	s.inventoryClient.On("ListParts", s.ctx, filter).Return(nil, gofakeit.Error())
	order, err := s.service.CreateOrder(s.ctx, createOrderIn)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInventoryUnavailable)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestCreateOrderAtLeastOnePartWasNotFound() {
	createOrderIn := s.GenerateCreateOrderIn()
	filter := model.PartsFilter{
		UUIDs: createOrderIn.PartUuids,
	}
	s.inventoryClient.On("ListParts", s.ctx, filter).Return([]model.Part{}, nil)
	order, err := s.service.CreateOrder(s.ctx, createOrderIn)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrAtLeastOnePartWasNotFound)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestCreateOrderSuccess() {
	createOrderIn := s.GenerateCreateOrderIn()
	filter := model.PartsFilter{
		UUIDs: createOrderIn.PartUuids,
	}

	parts := s.GenerateParts(filter.UUIDs)

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(parts, nil)
	s.repo.On("CreateOrder", s.ctx, mock.MatchedBy(func(o model.Order) bool {
		return o.OrderUUID != "" && o.UserUUID == createOrderIn.UserUUID && slices.Equal(o.PartUuids, createOrderIn.PartUuids) && o.TotalPrice == s.TotalPrice(createOrderIn.PartUuids, parts) && o.Status != nil && *o.Status == model.OrderStatusPENDINGPAYMENT
	})).Return(nil)
	order, err := s.service.CreateOrder(s.ctx, createOrderIn)

	s.Require().NoError(err)

	s.Require().NotEmpty(order)
	s.Require().Equal(createOrderIn.UserUUID, order.UserUUID)
	s.Require().Equal(createOrderIn.PartUuids, order.PartUuids)
	s.Require().Equal(s.TotalPrice(createOrderIn.PartUuids, parts), order.TotalPrice)

	s.Require().NotNil(order.Status)
	s.Require().Equal(model.OrderStatusPENDINGPAYMENT, *order.Status)
}

func (s *ServiceSuite) TestCreateOrderRepoError() {
	createOrderIn := s.GenerateCreateOrderIn()
	filter := model.PartsFilter{
		UUIDs: createOrderIn.PartUuids,
	}

	parts := s.GenerateParts(filter.UUIDs)

	randomErr := gofakeit.Error()
	s.inventoryClient.On("ListParts", s.ctx, filter).Return(parts, nil)
	s.repo.On("CreateOrder", s.ctx, mock.MatchedBy(func(o model.Order) bool {
		return o.OrderUUID != "" && o.UserUUID == createOrderIn.UserUUID && slices.Equal(o.PartUuids, createOrderIn.PartUuids) && o.TotalPrice == s.TotalPrice(createOrderIn.PartUuids, parts) && o.Status != nil && *o.Status == model.OrderStatusPENDINGPAYMENT
	})).Return(randomErr)
	order, err := s.service.CreateOrder(s.ctx, createOrderIn)

	s.Require().Error(err)
	s.Require().ErrorIs(err, randomErr)

	s.Require().Empty(order)
}
