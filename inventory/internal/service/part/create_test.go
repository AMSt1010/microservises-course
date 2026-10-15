package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

// GenerateFakePartInfo генерирует валидный экземпляр структуры PartInfo с реалистичными тестовыми данными.
func GenerateFakePartInfo() model.PartInfo {
	validCategories := []model.Category{
		model.CategoryEngine,
		model.CategoryFuel,
		model.CategoryPorthole,
		model.CategoryWing,
	}
	randomCategory := validCategories[gofakeit.Number(0, len(validCategories)-1)]

	dimensions := model.Dimensions{
		Length: gofakeit.Float64Range(10.0, 500.0),
		Width:  gofakeit.Float64Range(5.0, 200.0),
		Height: gofakeit.Float64Range(2.0, 150.0),
		Weight: gofakeit.Float64Range(0.5, 50.0),
	}

	manufacturer := model.Manufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}

	tagsCount := gofakeit.Number(2, 5)
	tags := make([]string, 0, tagsCount)
	for i := 0; i < tagsCount; i++ {
		tags = append(tags, gofakeit.Noun())
	}

	metadata := map[string]model.CustomValue{
		"serial_code":  model.NewStringValue(gofakeit.UUID()),
		"batch_number": model.NewInt64Value(gofakeit.Int64()),
		"tolerance":    model.NewDoubleValue(gofakeit.Float64Range(0.001, 0.05)),
		"is_certified": model.NewBoolValue(gofakeit.Bool()),
	}

	return model.PartInfo{
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.ProductDescription(),
		Price:         gofakeit.Price(50.0, 10000.0),
		StockQuantity: int64(gofakeit.Number(0, 500)),
		Category:      randomCategory,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          tags,
		Metadata:      metadata,
	}
}

func (s *ServiceSuite) TestCreatePartSuccess() {
	info := GenerateFakePartInfo()

	s.invRepository.On("CreatePart", s.ctx, mock.MatchedBy(func(p model.Part) bool {
		return p.UUID != "" && p.Info.Name == info.Name && p.Info.Price == info.Price
	})).Return(nil).Once()

	newUUID, err := s.service.CreatePart(s.ctx, info)

	s.Require().NoError(err)
	s.Require().NotEmpty(newUUID)
}

func (s *ServiceSuite) TestCreatePartValidationError() {
	info := GenerateFakePartInfo()
	info.Name = "" // Пустое имя нарушает валидацию

	newUUID, err := s.service.CreatePart(s.ctx, info)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartNameCannotBeEmpty)
	s.Require().Empty(newUUID)
}

func (s *ServiceSuite) TestCreatePartContextCancelled() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	info := GenerateFakePartInfo()

	newUUID, err := s.service.CreatePart(ctx, info)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
	s.Require().Empty(newUUID)
}

func (s *ServiceSuite) TestCreatePartRepoError() {
	info := GenerateFakePartInfo()
	repoErr := gofakeit.Error()

	s.invRepository.On("CreatePart", s.ctx, mock.Anything).Return(repoErr).Once()

	newUUID, err := s.service.CreatePart(s.ctx, info)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(newUUID)
}
