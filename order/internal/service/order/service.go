package order

import (
	grpcClient "github.com/AMSt1010/microservises-course/week2/order/internal/client/grpc"
	"github.com/AMSt1010/microservises-course/week2/order/internal/repository"
	def "github.com/AMSt1010/microservises-course/week2/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	repo            repository.OrderRepository
	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient
}

func NewService(repo repository.OrderRepository, inventoryClient grpcClient.InventoryClient, paymentClient grpcClient.PaymentClient) *service {
	return &service{
		repo:            repo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
