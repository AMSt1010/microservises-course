package v1

import (
	"github.com/AMSt1010/microservises-course/week2/order/internal/service"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

type api struct {
	service service.OrderService
}

func NewAPI(s service.OrderService) *api {
	return &api{
		service: s,
	}
}

func (a *api) GenerateError(code int, message string) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: code,
		Response: orderV1.GenericError{
			Code:    code,
			Message: message,
		},
	}
}
