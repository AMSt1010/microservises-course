package v1

import (
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/service"
	invV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

type api struct {
	invV1.UnimplementedInventoryServiceServer
	InvService service.InventoryService
}

func NewApi(service service.InventoryService) *api {
	return &api{
		InvService: service,
	}
}
