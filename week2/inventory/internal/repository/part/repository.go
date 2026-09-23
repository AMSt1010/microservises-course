package part

import (
	"strings"
	"sync"

	def "github.com/AMSt1010/microservises-course/week2/inventory/internal/repository"

	repoModel "github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu            sync.RWMutex
	data          map[string]repoModel.Part
	filterIndices map[string]partsFilterStorage
}

func NewRepository() *repository {
	return &repository{
		data:          make(map[string]repoModel.Part),
		filterIndices: make(map[string]partsFilterStorage),
	}
}

// Плоская структура для быстрого поиска по фильтру
type partsFilterStorage struct {
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

func newPartsFilterStorage(part repoModel.Part) partsFilterStorage {
	return partsFilterStorage{
		uuid:                part.UUID,
		name:                strings.ToLower(strings.TrimSpace(part.Info.Name)),
		category:            int32(part.Info.Category),
		manufacturerCountry: strings.ToLower(strings.TrimSpace(part.Info.Manufacturer.Country)),
		tags:                normalizationSearchStrings(part.Info.Tags),
	}
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
