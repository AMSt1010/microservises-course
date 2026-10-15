package part

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

func (s *RepositorySuite) TestCreatePartSuccess() {
	part := model.Part{
		UUID:      uuid.NewString(),
		Info:      generateRepoTestPartInfo(),
		CreatedAt: time.Now().UTC(),
	}

	err := s.repo.CreatePart(s.ctx, part)

	s.Require().NoError(err)

	// Проверяем прямое сохранение в памяти
	s.Require().Contains(s.repo.data, part.UUID)
	s.Require().Contains(s.repo.filterIndices, part.UUID)
}

func (s *RepositorySuite) TestCreatePartContextCancelled() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	part := model.Part{
		UUID:      uuid.NewString(),
		Info:      generateRepoTestPartInfo(),
		CreatedAt: time.Now().UTC(),
	}

	err := s.repo.CreatePart(ctx, part)

	s.Require().Error(err)
	s.Require().ErrorIs(err, context.Canceled)
}
