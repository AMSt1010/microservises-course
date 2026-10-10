package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/AMSt1010/microservises-course/week2/order/internal/api/order/v1"
	inventoryV1Client "github.com/AMSt1010/microservises-course/week2/order/internal/client/grpc/inventory/v1"
	paymentV1Client "github.com/AMSt1010/microservises-course/week2/order/internal/client/grpc/payment/v1"
	repo "github.com/AMSt1010/microservises-course/week2/order/internal/repository/order"
	service "github.com/AMSt1010/microservises-course/week2/order/internal/service/order"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/payment/v1"
)

const (
	InventoryServiceAddress = "localhost:50051"
	PaymentServiceAddress   = "localhost:50050"
	httpPort                = "8080"
	// Таймауты для HTTP сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	// создаем пул grpc соединений для сервиса Inventory
	connInventory, err := grpc.NewClient(
		InventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect inventory service: %v", err)
		return
	}
	// закрываем соединение
	defer func() {
		if closeErr := connInventory.Close(); closeErr != nil {
			log.Printf("failed to close inventory connection: %v", closeErr)
		}
	}()

	// создаем клиента для сервиса Inventory
	inventoryGeneratedClient := inventoryV1.NewInventoryServiceClient(connInventory)

	inventoryClient := inventoryV1Client.NewClient(inventoryGeneratedClient)

	// создаем пул grpc соединений для сервиса Payment
	// create a pool of grpc connection for payment service
	connPayment, err := grpc.NewClient(
		PaymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect payment service: %v", err)
		return
	}

	defer func() {
		if closeErr := connPayment.Close(); closeErr != nil {
			log.Printf("failed to close payment connection: %v", closeErr)
		}
	}()
	// создаем клиента для сервиса Payment
	// Create a client for payment service
	paymentGeneratedClient := paymentV1.NewPaymentServiceClient(connPayment)

	paymentClient := paymentV1Client.NewClient(paymentGeneratedClient)

	// Создаем структуру репозитория
	// Create a repository struct
	repo := repo.NewRepository()

	// Создаем структуру сервисного слоя
	// Create the service layer struct

	service := service.NewService(repo, inventoryClient, paymentClient)

	// Создаем структуру API слоя
	// Create the API layer struct
	api := api.NewAPI(service)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
	}
	// инициализируем роутер chi
	r := chi.NewRouter()

	// добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем Http сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 Http сервер запущен на порту %s \n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// log.Println("🛑 Завершение работы сервера")
	log.Println("🛑 Graceful shutdown of the server")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		// log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
		log.Printf("❌ Error during graceful shutdown: %v\n", err)
	}

	// log.Println("✅ Сервер остановлен")
	log.Println("✅ Server stopped")
}
