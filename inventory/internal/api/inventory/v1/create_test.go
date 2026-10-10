package v1

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/converter"
	invV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestCreatePartSuccess() {
	protoInfo := &invV1.Info{
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.ProductDescription(),
		Price:         gofakeit.Price(10, 1000),
		StockQuantity: int64(gofakeit.Number(1, 100)),
		Category:      invV1.Category_CATEGORY_ENGINE,
		Dimensions: &invV1.Dimensions{
			Length: 100.5,
			Width:  50.2,
			Height: 30.1,
			Weight: 15.0,
		},
		Manufacturer: &invV1.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{gofakeit.Noun(), gofakeit.Noun()},
		Metadata: map[string]*invV1.CustomValue{
			"serial": {
				Kind: &invV1.CustomValue_StringValue{
					StringValue: gofakeit.UUID(),
				},
			},
		},
	}

	req := &invV1.CreatePartRequest{
		Info: protoInfo,
	}

	expectedDomainInfo := converter.PartInfoToModel(protoInfo)
	expectedUUID := gofakeit.UUID()

	s.invService.On("CreatePart", s.ctx, expectedDomainInfo).Return(expectedUUID, nil)

	res, err := s.api.CreatePart(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedUUID, res.GetUuid())
}

func (s *APISuite) TestCreatePartServiceError() {
	protoInfo := &invV1.Info{
		Name: gofakeit.ProductName(),
	}

	req := &invV1.CreatePartRequest{
		Info: protoInfo,
	}

	expectedDomainInfo := converter.PartInfoToModel(protoInfo)
	serviceErr := gofakeit.Error()

	s.invService.On("CreatePart", s.ctx, expectedDomainInfo).Return("", serviceErr)

	res, err := s.api.CreatePart(s.ctx, req)

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
