package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	payment_v1 "github.com/alex100010/microservices-course/hw/hw_1/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50050

type paymentService struct {
	payment_v1.UnimplementedPaymentServiceServer
}

// Обрабатывает команду на оплату и возвращает `transaction_uuid`.
func (pS *paymentService) PayOrder(_ context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	transaction_uuid := uuid.NewString()

	log.Printf("Оплата прошла успешно , transaction_uuid: %v\nUUID заказа: %v\nUUID пользователя: %v\nСпособ оплаты: %v\n", transaction_uuid, req.GetOrderUuid(), req.GetUserUuid(), req.GetPaymentMethod())

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transaction_uuid,
	}, nil
}

func main() {
	// Создаем tcp соединение
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	// Создаем новый grpc сервер
	s := grpc.NewServer()

	// Регистрируем сервис на сервере
	payment_v1.RegisterPaymentServiceServer(s, &paymentService{})

	// Подключаем рефлексию
	reflection.Register(s)

	// Запускаем grpc сервер
	go func() {
		log.Printf("grpc server listening on: %v", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down grpc server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
