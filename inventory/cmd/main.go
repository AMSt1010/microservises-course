package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	inventory_v1 "github.com/alex100010/microservices-course/hw/hw_1/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

// Плоская структура для быстрого поиска по фильтру
type PartsFilterStorage struct {
	// Уникальный идентификатор детали
	uuid string

	// Название детали
	name string

	// Категория
	category int32

	// Страна производителя
	manufacturerCountry string

	// Теги для быстрого поиска
	tags []string
}

func NewPartsFilterStorage(part *inventory_v1.Part) PartsFilterStorage {
	return PartsFilterStorage{
		uuid:                part.GetUuid(),
		name:                part.GetInfo().GetName(),
		category:            int32(part.GetInfo().GetCategory()),
		manufacturerCountry: part.GetInfo().GetManufacturer().GetCountry(),
		tags:                normalizationSearchStrings(part.GetInfo().GetTags()),
	}
}

type inventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer
	mu           sync.RWMutex
	storage      map[string]*inventory_v1.Part
	filterSlices []PartsFilterStorage
}

func newInventoryService() *inventoryService {
	return &inventoryService{
		storage:      make(map[string]*inventory_v1.Part),
		filterSlices: make([]PartsFilterStorage, 0, 1000),
	}
}

// Записывает информацию о новой детали
func (inv *inventoryService) CreatePart(_ context.Context, req *inventory_v1.CreatePartRequest) (*inventory_v1.CreatePartResponse, error) {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	newUuid := uuid.NewString()

	part := &inventory_v1.Part{
		Uuid:      newUuid,
		Info:      req.GetInfo(),
		CreatedAt: timestamppb.New(time.Now()),
	}

	inv.storage[newUuid] = part

	partsFilter := NewPartsFilterStorage(part)
	inv.filterSlices = append(inv.filterSlices, partsFilter)

	log.Printf("Создана деталь с Uuid %v", newUuid)
	return &inventory_v1.CreatePartResponse{
		Uuid: newUuid,
	}, nil
}

func normalizationSearchStrings(inputStr []string) []string {
	normalizedString := make([]string, 0, len(inputStr))
	for _, n := range inputStr {
		if trimmed := strings.ToLower(strings.TrimSpace(n)); trimmed != "" {
			normalizedString = append(normalizedString, trimmed)
		}
	}
	return normalizedString
}

func containsAny(target string, searchTerms []string) bool {
	if len(searchTerms) == 0 {
		return true
	}
	lowerTarget := strings.ToLower(target)
	for _, term := range searchTerms {
		if strings.Contains(lowerTarget, term) {
			return true
		}
	}
	return false
}

func containsAnyTag(partTags, searchTags []string) bool {
	if len(searchTags) == 0 {
		return true
	}
	for _, searchTag := range searchTags {
		for _, partTag := range partTags {
			if strings.EqualFold(strings.TrimSpace(partTag), searchTag) {
				return true
			}
		}
	}
	return false
}

func matchCategories(itemCategory int32, categories []inventory_v1.Category) bool {
	if len(categories) == 0 {
		return true
	}
	for _, cat := range categories {
		if itemCategory == int32(cat) {
			return true
		}
	}
	return false
}

func (inv *inventoryService) collectSourceByUUIDs(uuids []string) []PartsFilterStorage {
	source := make([]PartsFilterStorage, 0, len(uuids))
	for _, id := range uuids {
		if part, ok := inv.storage[id]; ok && part != nil {
			source = append(source, NewPartsFilterStorage(part))
		}
	}
	return source
}

// Возвращает список деталей с возможностью фильтрации.
func (inv *inventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	filter := req.GetFilter()
	if filter == nil {
		result := make([]*inventory_v1.Part, 0, len(inv.storage))
		for _, v := range inv.storage {
			if v != nil {
				result = append(result, v)
			}
		}
		return &inventory_v1.ListPartsResponse{Parts: result}, nil
	}

	lowerNames := normalizationSearchStrings(filter.GetNames())
	lowerCountries := normalizationSearchStrings(filter.GetManufacturerCountries())
	categories := filter.GetCategories()

	var source []PartsFilterStorage
	if len(filter.GetUuids()) > 0 {
		source = inv.collectSourceByUUIDs(filter.GetUuids())
	} else {
		source = inv.filterSlices
	}

	result := make([]*inventory_v1.Part, 0, len(source))
	for _, item := range source {
		if !matchCategories(item.category, categories) {
			continue
		}
		if !containsAny(item.name, lowerNames) {
			continue
		}
		if !containsAny(item.manufacturerCountry, lowerCountries) {
			continue
		}
		if !containsAnyTag(item.tags, filter.GetTags()) {
			continue
		}

		if part := inv.storage[item.uuid]; part != nil {
			result = append(result, part)
		}
	}

	return &inventory_v1.ListPartsResponse{
		Parts: result,
	}, nil
}

// Возвращает информацию о детали по её UUID
func (inv *inventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	Uuid := req.GetUuid()

	if Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "Uuid cannot be empty")
	}

	part, ok := inv.storage[Uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with Uuid %s not exist", Uuid)
	}

	log.Printf("Деталь с Uuid %s получена: %+v", Uuid, part)

	return &inventory_v1.GetPartResponse{
		Part: part,
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
	service := newInventoryService()

	inventory_v1.RegisterInventoryServiceServer(s, service)

	// Подключаем рефлексию
	reflection.Register(s)

	// Запускаем grpc сервер
	go func() {
		log.Printf("grpc server listening on: %v", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v", err)
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
