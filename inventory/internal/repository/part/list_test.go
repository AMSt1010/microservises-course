package part

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

func (s *RepositorySuite) TestListPartsEmptyFilter() {
	info1 := generateRepoTestPartInfo()
	info2 := generateRepoTestPartInfo()

	_, err := s.repo.CreatePart(s.ctx, info1)
	s.Require().NoError(err)
	_, err = s.repo.CreatePart(s.ctx, info2)
	s.Require().NoError(err)

	parts, err := s.repo.ListParts(s.ctx, model.PartsFilter{})

	s.Require().NoError(err)
	s.Require().Len(parts, 2)
}

func (s *RepositorySuite) TestListPartsFilters() {
	info1 := generateRepoTestPartInfo()
	info1.Name = "Quantum Engine"
	info1.Category = model.CategoryEngine
	info1.Manufacturer.Country = "Germany"
	info1.Tags = []string{"quantum", "speed"}

	info2 := generateRepoTestPartInfo()
	info2.Name = "Solar Wing"
	info2.Category = model.CategoryWing
	info2.Manufacturer.Country = "Japan"
	info2.Tags = []string{"solar", "light"}

	uuid1, err := s.repo.CreatePart(s.ctx, info1)
	s.Require().NoError(err)

	uuid2, err := s.repo.CreatePart(s.ctx, info2)
	s.Require().NoError(err)

	// 1. Фильтр по UUID
	byUUID, err := s.repo.ListParts(s.ctx, model.PartsFilter{UUIDs: []string{uuid1}})
	s.Require().NoError(err)
	s.Require().Len(byUUID, 1)
	s.Require().Equal(uuid1, byUUID[0].UUID)

	// 2. Фильтр по имени и категории
	byNameAndCat, err := s.repo.ListParts(s.ctx, model.PartsFilter{
		Names:      []string{"quantum"},
		Categories: []model.Category{model.CategoryEngine},
	})
	s.Require().NoError(err)
	s.Require().Len(byNameAndCat, 1)
	s.Require().Equal(uuid1, byNameAndCat[0].UUID)

	// 3. Фильтр по стране и тегу
	byCountryAndTag, err := s.repo.ListParts(s.ctx, model.PartsFilter{
		ManufacturerCountries: []string{"japan"},
		Tags:                  []string{"solar"},
	})
	s.Require().NoError(err)
	s.Require().Len(byCountryAndTag, 1)
	s.Require().Equal(uuid2, byCountryAndTag[0].UUID)

	// 4. Фильтр без совпадений
	emptyRes, err := s.repo.ListParts(s.ctx, model.PartsFilter{
		Names: []string{"NonExistentPart"},
	})
	s.Require().NoError(err)
	s.Require().Empty(emptyRes)
}

func (s *RepositorySuite) TestListPartsCanceledContext() {
	canceledCtx, cancel := context.WithCancel(s.ctx)
	cancel()

	parts, err := s.repo.ListParts(canceledCtx, model.PartsFilter{})

	s.Require().Error(err)
	s.Require().Nil(parts)
}
