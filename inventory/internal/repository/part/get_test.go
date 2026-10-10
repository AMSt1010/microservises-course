package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

func (s *RepositorySuite) TestGetPartSuccess() {
	info := generateRepoTestPartInfo()
	uuid, err := s.repo.CreatePart(s.ctx, info)
	s.Require().NoError(err)

	part, err := s.repo.GetPart(s.ctx, uuid)

	s.Require().NoError(err)
	s.Require().Equal(uuid, part.UUID)
	s.Require().Equal(info.Name, part.Info.Name)
	s.Require().Equal(info.Price, part.Info.Price)
}

func (s *RepositorySuite) TestGetPartEmptyUUID() {
	_, err := s.repo.GetPart(s.ctx, "")

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUUIDCannotBeEmpty)
}

func (s *RepositorySuite) TestGetPartNotFound() {
	_, err := s.repo.GetPart(s.ctx, gofakeit.UUID())

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartNotFound)
}
