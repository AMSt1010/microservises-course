package part

import (
	"time"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

// GenerateFakePart генерирует валидную доменную модель Part для тестов.
func GenerateFakePart() model.Part {
	now := time.Now().Truncate(time.Second)

	return model.Part{
		UUID:      gofakeit.UUID(),
		Info:      GenerateFakePartInfo(),
		CreatedAt: now,
		UpdatedAt: new(now.Add(time.Hour)),
	}
}

func (s *ServiceSuite) TestGetPartSuccess() {
	expectedPart := GenerateFakePart()
	targetUUID := expectedPart.UUID

	s.invRepository.On("GetPart", s.ctx, targetUUID).Return(expectedPart, nil)

	actualPart, err := s.service.GetPart(s.ctx, targetUUID)

	s.Require().NoError(err)
	s.Require().Equal(expectedPart, actualPart)
}

func (s *ServiceSuite) TestGetPartRepoError() {
	targetUUID := gofakeit.UUID()
	repoErr := gofakeit.Error()

	s.invRepository.On("GetPart", s.ctx, targetUUID).Return(model.Part{}, repoErr)

	actualPart, err := s.service.GetPart(s.ctx, targetUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(actualPart)
}

func (s *ServiceSuite) TestGetPartNotFound() {
	targetUUID := gofakeit.UUID()

	s.invRepository.On("GetPart", s.ctx, targetUUID).Return(model.Part{}, model.ErrPartNotFound)

	actualPart, err := s.service.GetPart(s.ctx, targetUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartNotFound)
	s.Require().Empty(actualPart)
}
