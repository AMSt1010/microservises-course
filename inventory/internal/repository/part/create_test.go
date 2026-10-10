package part

func (s *RepositorySuite) TestCreatePartSuccess() {
	info := generateRepoTestPartInfo()

	uuid, err := s.repo.CreatePart(s.ctx, info)

	s.Require().NoError(err)
	s.Require().NotEmpty(uuid)

	// Проверяем прямое сохранение в памяти
	s.Require().Contains(s.repo.data, uuid)
	s.Require().Contains(s.repo.filterIndices, uuid)
}
