package part

import (
	"context"
	"strings"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/converter"
)

func containsAny(target string, searchTerms []string) bool {
	if len(searchTerms) == 0 {
		return true
	}
	for _, term := range searchTerms {
		if strings.Contains(target, term) {
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
			if partTag == searchTag {
				return true
			}
		}
	}
	return false
}

func matchCategories(itemCategory int32, categories []model.Category) bool {
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

func (r *repository) collectSourceByUUIDs(uuids []string) map[string]partsFilterStorage {
	source := make(map[string]partsFilterStorage, len(uuids))
	for _, id := range uuids {
		if filterItem, ok := r.filterIndices[id]; ok {
			source[id] = filterItem
		}
	}
	return source
}

// Возвращает список деталей с возможностью фильтрации.
func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if isFilterEmpty(filter) {
		result := make([]model.Part, 0, len(r.data))
		for _, v := range r.data {
			result = append(result, converter.PartToModel(v))
		}
		return result, nil
	}

	lowerNames := normalizationSearchStrings(filter.Names)
	lowerCountries := normalizationSearchStrings(filter.ManufacturerCountries)
	lowerTags := normalizationSearchStrings(filter.Tags)
	categories := filter.Categories

	var source map[string]partsFilterStorage
	if len(filter.UUIDs) > 0 {
		source = r.collectSourceByUUIDs(filter.UUIDs)
	} else {
		source = r.filterIndices
	}

	result := make([]model.Part, 0, len(source))
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
		if !containsAnyTag(item.tags, lowerTags) {
			continue
		}

		part, ok := r.data[item.uuid]
		if !ok {
			continue // Защита от рассинхронизации индекса и данных
		}

		result = append(result, converter.PartToModel(part))
	}

	return result, nil
}

func isFilterEmpty(f model.PartsFilter) bool {
	return len(f.UUIDs) == 0 &&
		len(f.Names) == 0 &&
		len(f.Categories) == 0 &&
		len(f.ManufacturerCountries) == 0 &&
		len(f.Tags) == 0
}
