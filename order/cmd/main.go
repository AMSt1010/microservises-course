package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderV1 "github.com/alex100010/microservices-course/hw/hw_1/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/alex100010/microservices-course/hw/hw_1/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/alex100010/microservices-course/hw/hw_1/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	InventoryServiceAddress = "localhost:50051"
	PaymentServiceAddress   = "localhost:50050"
	httpPort                = "8080"
	// Таймауты для HTTP сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// Представляет потокобезопасное хранилище данных о заказах
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

// Конструктор - создает новое хранилище данных о заказах
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

// Создает новый заказ в Map или обновляет информацию о заказе
func (s *OrderStorage) SaveOrder(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()
	orderCopy := *order
	s.orders[order.OrderUUID] = &orderCopy
}

// Получает заказ из Map по Uuid заказа
func (s *OrderStorage) GetOrder(orderUUID string) (orderV1.OrderDto, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return orderV1.OrderDto{}, false
	}
	return *order, true
}

func (s *OrderStorage) MarkAsPaidCAS(orderUUID, transactionUUID string, method orderV1.PaymentMethod) (*orderV1.OrderDto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil, errors.New("not_found")
	}

	// Атомарно проверяем актуальный статус на момент завершения платежа
	if order.Status == orderV1.OrderStatusPAID {
		return nil, errors.New("already_paid")
	}
	if order.Status == orderV1.OrderStatusCANCELLED {
		return nil, errors.New("cancelled")
	}

	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptString(transactionUUID)
	order.PaymentMethod = orderV1.NewOptPaymentMethod(method)

	orderCopy := *order
	return &orderCopy, nil
}

func (s *OrderStorage) CancelOrder(orderUUID string) (orderV1.OrderStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return "", errors.New("not_found")
	}
	if order.Status == orderV1.OrderStatusPAID {
		return order.Status, errors.New("conflict")
	}
	if order.Status == orderV1.OrderStatusCANCELLED {
		return order.Status, nil // Уже отменен
	}

	order.Status = orderV1.OrderStatusCANCELLED
	return orderV1.OrderStatusCANCELLED, nil
}

// Реализиует интерфейс для обработки запросов к API заказов
type orderHandler struct {
	storage         *OrderStorage
	inventoryClient inventory_v1.InventoryServiceClient
	paymentClient   payment_v1.PaymentServiceClient
}

// Создает новый обработчик запросов к API заказов
func NewOrderHandler(s *OrderStorage, invClient inventory_v1.InventoryServiceClient, pClient payment_v1.PaymentServiceClient) *orderHandler {
	return &orderHandler{
		storage:         s,
		inventoryClient: invClient,
		paymentClient:   pClient,
	}
}

// Проверяет, что все запрашиваемые детали существуют
func areAllPartsExist(checkUuids []string, parts []*inventory_v1.Part) bool {
	getParts := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		getParts[part.GetUuid()] = struct{}{}
	}
	for _, checkuuid := range checkUuids {
		if _, ok := getParts[checkuuid]; !ok {
			return false
		}
	}
	return true
}

// Считает общую цену заказа с учетом количества каждого запрошенного UUID
func totalPrice(requestedUUIDs []string, parts []*inventory_v1.Part) float64 {
	priceMap := make(map[string]float64, len(parts))
	for _, part := range parts {
		priceMap[part.GetUuid()] = part.GetInfo().GetPrice()
	}

	var total float64
	for _, id := range requestedUUIDs {
		total += priceMap[id]
	}
	return total
}

// POST /api/v1/orders
// Создаёт новый заказ на основе выбранных пользователем деталей.
func (h *orderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error) {
	if req.GetUserUUID() == "" {
		return nil, errors.New("user_uuid cannot be empty")
	}

	if len(req.GetPartUuids()) == 0 {
		return nil, errors.New("part_uuids cannot be empty")
	}
	// генерируем id заказа
	order_uuid := uuid.NewString()

	// список запрашиваемых деталей
	uuids := req.GetPartUuids()
	// Обращаемся к сервису Inventory и получаем список деталей
	filter := &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: uuids,
		},
	}
	listPartsResponse, err := h.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("create order: inventory service ListParts failed: %w", err)
	}

	// Список деталей
	parts := listPartsResponse.GetParts()

	// Если хотя бы одной детали нет — возвращает ошибку.
	if ok := areAllPartsExist(uuids, parts); !ok {
		return nil, errors.New("create order: at least one part was not found")
	}

	// общая сумма заказа
	total := totalPrice(uuids, parts)

	// Создаем заказ с полученными данными
	order := &orderV1.OrderDto{
		OrderUUID:  order_uuid,
		UserUUID:   req.GetUserUUID(),
		PartUuids:  uuids,
		TotalPrice: float32(total),
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	// Обновляем данные в хранилище
	h.storage.SaveOrder(order)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order_uuid,
		TotalPrice: float32(total),
	}, nil
}

func mapPaymentMethodToProto(m orderV1.PaymentMethod) payment_v1.PaymentMethod {
	switch m {
	case orderV1.PaymentMethodCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

// POST /api/v1/orders/{order_uuid}/pay
// Проводит оплату ранее созданного заказа.
func (h *orderHandler) OrderPay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.OrderPayParams) (orderV1.OrderPayRes, error) {
	// 1. Быстрое чтение и первичная валидация
	order, ok := h.storage.GetOrder(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order for uuid '%s' not found", params.OrderUUID),
		}, nil
	}

	if order.GetStatus() == orderV1.OrderStatusPAID {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Order for uuid '%s' has already been paid for", params.OrderUUID),
		}, nil
	}
	if order.GetStatus() == orderV1.OrderStatusCANCELLED {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Cannot pay cancelled order '%s'", params.OrderUUID),
		}, nil
	}

	protoMethod := mapPaymentMethodToProto(req.PaymentMethod)
	if protoMethod == payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid or unsupported payment method '%s'", req.PaymentMethod),
		}, nil
	}

	payment := &payment_v1.PayOrderRequest{
		OrderUuid:     order.GetOrderUUID(),
		UserUuid:      order.GetUserUUID(),
		PaymentMethod: protoMethod,
	}

	// 2. Сетевой вызов (мьютекс не держим!)
	transactionUuid, err := h.paymentClient.PayOrder(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("order pay: payment service PayOrder failed: %w", err)
	}

	// 3. Атомарная фиксация статуса
	updatedOrder, err := h.storage.MarkAsPaidCAS(params.OrderUUID, transactionUuid.GetTransactionUuid(), req.PaymentMethod)
	if err != nil {
		switch err.Error() {
		case "not_found":
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Order for uuid '%s' not found", params.OrderUUID),
			}, nil
		case "already_paid":
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Order for uuid '%s' has already been paid for", params.OrderUUID),
			}, nil
		case "cancelled":
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Order for uuid '%s' is cancelled and cannot be paid", params.OrderUUID),
			}, nil
		default:
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		}
	}
	txUUID, _ := updatedOrder.TransactionUUID.Get()
	return &orderV1.PayOrderResponse{
		TransactionUUID: txUUID,
	}, nil
}

// GET /api/v1/orders/{order_uuid}
// Возвращает информацию о заказе.
func (h *orderHandler) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order, ok := h.storage.GetOrder(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order for uuid '" + params.OrderUUID + "' not found",
		}, nil
	}
	return &orderV1.GetOrderResponse{
		Order: order,
	}, nil
}

// POST /api/v1/orders/{order_uuid}/cancel
// Отменяет заказ.
func (h *orderHandler) OrderCancelByUuid(ctx context.Context, params orderV1.OrderCancelByUuidParams) (orderV1.OrderCancelByUuidRes, error) {
	_, err := h.storage.CancelOrder(params.OrderUUID)
	if err != nil {
		switch err.Error() {
		case "not_found":
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Order by uuid '%s' not found", params.OrderUUID),
			}, nil
		case "conflict":
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: "The order has already been paid for and cannot be cancelled",
			}, nil
		default:
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		}
	}

	return &orderV1.OrderCancelByUuidNoContent{}, nil
}

// Used for common default response.
func (h *orderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

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
	inventoryClient := inventory_v1.NewInventoryServiceClient(connInventory)

	// создаем пул grpc соединений для сервиса Payment
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
	paymentClient := payment_v1.NewPaymentServiceClient(connPayment)

	// Создаем хранилище для заказов
	orderStorage := NewOrderStorage()

	// Создаем обработчик API заказов
	orderHandler := NewOrderHandler(orderStorage, inventoryClient, paymentClient)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
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

	log.Println("🛑 Завершение работы сервера")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
