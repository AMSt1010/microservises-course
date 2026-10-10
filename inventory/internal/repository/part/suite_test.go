package part

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

type RepositorySuite struct {
	suite.Suite

	ctx  context.Context
	repo *repository
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = NewRepository()
}

func (s *RepositorySuite) TearDownTest() {}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

// generateRepoTestPartInfo генерирует валидные данные для тестов репозитория
func generateRepoTestPartInfo() model.PartInfo {
	return model.PartInfo{
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(5),
		Price:         gofakeit.Price(10, 500),
		StockQuantity: 50,
		Category:      model.CategoryEngine,
		Dimensions: model.Dimensions{
			Length: 10,
			Width:  10,
			Height: 10,
			Weight: 5,
		},
		Manufacturer: model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: "Germany",
			Website: gofakeit.URL(),
		},
		Tags: []string{"titanium", "turbo"},
		Metadata: map[string]model.CustomValue{
			"serial": model.NewStringValue(gofakeit.UUID()),
		},
	}
}
