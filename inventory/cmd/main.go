package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	invV1API "github.com/AMSt1010/microservises-course/week2/inventory/internal/api/inventory/v1"
	invRepository "github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/part"
	invService "github.com/AMSt1010/microservises-course/week2/inventory/internal/service/part"
	invV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем grpc сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис

	r := invRepository.NewRepository()
	service := invService.NewService(r)
	api := invV1API.NewApi(service)

	invV1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	// запускаем сервер
	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("Server stopped gracefully")
	case <-ctx.Done():
		log.Println("Shutdown timeout exceeded, forcing stop")
		s.Stop()
	}
}
