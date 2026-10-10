package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

// GenerateFakePartInfo генерирует валидный экземпляр структуры PartInfo с реалистичными тестовыми данными.
func GenerateFakePartInfo() model.PartInfo {
	// 1. Выбор случайной категории из допустимого диапазона (исключая CategoryUnspecified)
	validCategories := []model.Category{
		model.CategoryEngine,
		model.CategoryFuel,
		model.CategoryPorthole,
		model.CategoryWing,
	}
	randomCategory := validCategories[gofakeit.Number(0, len(validCategories)-1)]

	// 2. Генерация физических габаритов
	dimensions := model.Dimensions{
		Length: gofakeit.Float64Range(10.0, 500.0),
		Width:  gofakeit.Float64Range(5.0, 200.0),
		Height: gofakeit.Float64Range(2.0, 150.0),
		Weight: gofakeit.Float64Range(0.5, 50.0),
	}

	// 3. Генерация информации о производителе
	manufacturer := model.Manufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}

	// 4. Генерация тегов (от 2 до 5 уникальных слов)
	tagsCount := gofakeit.Number(2, 5)
	tags := make([]string, 0, tagsCount)
	for i := 0; i < tagsCount; i++ {
		tags = append(tags, gofakeit.Noun())
	}

	// 5. Генерация метаданных с гарантией соблюдения валидации CustomValue (ровно одно поле заполнено)
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

	expectedUUID := gofakeit.UUID()

	s.invRepository.On("CreatePart", s.ctx, info).Return(expectedUUID, nil)

	UUID, err := s.service.CreatePart(s.ctx, info)

	s.Require().NoError(err)
	s.Require().Equal(expectedUUID, UUID)
}

func (s *ServiceSuite) TestCreatePartRepoError() {
	info := GenerateFakePartInfo()

	repoErr := gofakeit.Error()

	s.invRepository.On("CreatePart", s.ctx, info).Return("", repoErr)

	UUID, err := s.service.CreatePart(s.ctx, info)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(UUID)
}
