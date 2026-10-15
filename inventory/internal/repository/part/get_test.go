package part

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

func (s *RepositorySuite) TestGetPartSuccess() {
	expectedPart := model.Part{
		UUID:      uuid.NewString(),
		Info:      generateRepoTestPartInfo(),
		CreatedAt: time.Now().UTC(),
	}
	err := s.repo.CreatePart(s.ctx, expectedPart)
	s.Require().NoError(err)

	part, err := s.repo.GetPart(s.ctx, expectedPart.UUID)

	s.Require().NoError(err)
	s.Require().Equal(expectedPart.UUID, part.UUID)
	s.Require().Equal(expectedPart.Info.Name, part.Info.Name)
	s.Require().Equal(expectedPart.Info.Price, part.Info.Price)
}

func (s *RepositorySuite) TestGetPartNotFound() {
	_, err := s.repo.GetPart(s.ctx, gofakeit.UUID())

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartNotFound)
}

func (s *RepositorySuite) TestGetPartContextCancelled() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	_, err := s.repo.GetPart(ctx, gofakeit.UUID())

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
}
