package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/protobuf/proto"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	invV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestListPartsSuccess() {
	protoFilter := &invV1.PartsFilter{
		Uuids:                 []string{gofakeit.UUID()},
		Names:                 []string{gofakeit.ProductName()},
		Categories:            []invV1.Category{invV1.Category_CATEGORY_FUEL},
		ManufacturerCountries: []string{gofakeit.Country()},
		Tags:                  []string{gofakeit.Noun()},
	}

	req := &invV1.ListPartsRequest{
		Filter: protoFilter,
	}

	now := time.Now().Truncate(time.Second)
	domainParts := []model.Part{
		{
			UUID: protoFilter.Uuids[0],
			Info: model.PartInfo{
				Name:     protoFilter.Names[0],
				Category: model.CategoryFuel,
			},
			CreatedAt: now,
		},
	}

	expectedDomainFilter := converter.PartsFilterToModel(protoFilter)

	s.invService.On("ListParts", s.ctx, expectedDomainFilter).Return(domainParts, nil)

	res, err := s.api.ListParts(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.GetParts(), len(domainParts))

	expectedProtoParts := converter.PartsToProto(domainParts)
	for i, p := range res.GetParts() {
		s.Require().True(proto.Equal(expectedProtoParts[i], p))
	}
}

func (s *APISuite) TestListPartsServiceError() {
	req := &invV1.ListPartsRequest{
		Filter: &invV1.PartsFilter{},
	}

	expectedDomainFilter := converter.PartsFilterToModel(req.GetFilter())
	serviceErr := gofakeit.Error()

	s.invService.On("ListParts", s.ctx, expectedDomainFilter).Return(nil, serviceErr)

	res, err := s.api.ListParts(s.ctx, req)

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
