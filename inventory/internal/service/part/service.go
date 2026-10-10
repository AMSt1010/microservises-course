package part

import (
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/repository"
	def "github.com/AMSt1010/microservises-course/week2/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	invRepository repository.InventoryRepository
}

func NewService(repo repository.InventoryRepository) *service {
	return &service{
		invRepository: repo,
	}
}
