package v1

import (
	"github.com/AMSt1010/microservises-course/week2/payment/internal/service"
	PaymentV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/payment/v1"
)

type api struct {
	PaymentV1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewAPI(s service.PaymentService) *api {
	return &api{
		paymentService: s,
	}
}
