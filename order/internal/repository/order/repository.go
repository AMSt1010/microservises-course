package order

import (
	"sync"

	def "github.com/AMSt1010/microservises-course/week2/order/internal/repository"
	repoModel "github.com/AMSt1010/microservises-course/week2/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu     sync.RWMutex
	orders map[string]repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]repoModel.Order),
	}
}
