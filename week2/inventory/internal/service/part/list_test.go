package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

// GenerateFakePartsFilter генерирует структуру фильтрации со случайными критериями.
func GenerateFakePartsFilter() model.PartsFilter {
	return model.PartsFilter{
		UUIDs: []string{
			gofakeit.UUID(),
			gofakeit.UUID(),
		},
		Names: []string{
			gofakeit.Word(),
		},
		Categories: []model.Category{
			model.CategoryEngine,
			model.CategoryWing,
		},
		ManufacturerCountries: []string{
			gofakeit.Country(),
		},
		Tags: []string{
			gofakeit.Noun(),
		},
	}
}

func (s *ServiceSuite) TestListPartsSuccess() {
	filter := GenerateFakePartsFilter()

	count := gofakeit.Number(2, 5)
	expectedParts := make([]model.Part, 0, count)
	for range count {
		expectedParts = append(expectedParts, GenerateFakePart())
	}

	s.invRepository.On("ListParts", s.ctx, filter).Return(expectedParts, nil)

	actualParts, err := s.service.ListParts(s.ctx, filter)

	s.Require().NoError(err)
	s.Require().Equal(expectedParts, actualParts)
	s.Require().Len(actualParts, count)
}

func (s *ServiceSuite) TestListPartsEmptyResult() {
	filter := model.PartsFilter{}
	expectedParts := []model.Part{}

	s.invRepository.On("ListParts", s.ctx, filter).Return(expectedParts, nil)

	actualParts, err := s.service.ListParts(s.ctx, filter)

	s.Require().NoError(err)
	s.Require().Empty(actualParts)
}

func (s *ServiceSuite) TestListPartsRepoError() {
	filter := GenerateFakePartsFilter()
	repoErr := gofakeit.Error()

	s.invRepository.On("ListParts", s.ctx, filter).Return([]model.Part{}, repoErr)

	actualParts, err := s.service.ListParts(s.ctx, filter)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(actualParts)
}
