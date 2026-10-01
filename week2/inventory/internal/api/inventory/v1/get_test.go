package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	invV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPartSuccess() {
	targetUUID := gofakeit.UUID()
	now := time.Now().Truncate(time.Second)

	domainPart := model.Part{
		UUID: targetUUID,
		Info: model.PartInfo{
			Name:          gofakeit.ProductName(),
			Description:   gofakeit.ProductDescription(),
			Price:         199.99,
			StockQuantity: 10,
			Category:      model.CategoryPorthole,
			Dimensions: model.Dimensions{
				Length: 10,
				Width:  10,
				Height: 5,
				Weight: 2,
			},
			Manufacturer: model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{gofakeit.Word()},
			Metadata: map[string]model.CustomValue{
				"code": model.NewStringValue(gofakeit.UUID()),
			},
		},
		CreatedAt: now,
		UpdatedAt: new(now.Add(time.Hour)),
	}

	req := &invV1.GetPartRequest{
		Uuid: targetUUID,
	}

	s.invService.On("GetPart", s.ctx, targetUUID).Return(domainPart, nil)

	res, err := s.api.GetPart(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().True(proto.Equal(converter.PartToProto(domainPart), res.GetPart()))
}

func (s *APISuite) TestGetPartEmptyUUIDError() {
	req := &invV1.GetPartRequest{
		Uuid: "",
	}

	s.invService.On("GetPart", s.ctx, "").Return(model.Part{}, model.ErrUUIDCannotBeEmpty)

	res, err := s.api.GetPart(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
	s.Require().Equal("uuid cannot be empty", st.Message())
}

func (s *APISuite) TestGetPartNotFoundError() {
	targetUUID := gofakeit.UUID()
	req := &invV1.GetPartRequest{
		Uuid: targetUUID,
	}

	s.invService.On("GetPart", s.ctx, targetUUID).Return(model.Part{}, model.ErrPartNotFound)

	res, err := s.api.GetPart(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
	s.Require().Contains(st.Message(), targetUUID)
}

func (s *APISuite) TestGetPartServiceInternalError() {
	targetUUID := gofakeit.UUID()
	req := &invV1.GetPartRequest{
		Uuid: targetUUID,
	}
	unexpectedErr := gofakeit.Error()

	s.invService.On("GetPart", s.ctx, targetUUID).Return(model.Part{}, unexpectedErr)

	res, err := s.api.GetPart(s.ctx, req)

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpectedErr)
	s.Require().Nil(res)
}
