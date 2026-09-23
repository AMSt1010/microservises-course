package converter

import (
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	repoModel "github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/model"
)

// PartToRepoModel конвертирует корневую сущность детали в модель репозитория
func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		UUID:      part.UUID,
		Info:      PartInfoToRepoModel(part.Info),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

// PartInfoToRepoModel конвертирует характеристики детали в модель репозитория
func PartInfoToRepoModel(info model.PartInfo) repoModel.PartInfo {
	return repoModel.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Category:      CategoryToRepoModel(info.Category),
		Dimensions:    DimensionsToRepoModel(info.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(info.Manufacturer),
		Tags:          info.Tags,
		Metadata:      MetadataToRepoModel(info.Metadata),
	}
}

// DimensionsToRepoModel конвертирует габариты в модель репозитория
func DimensionsToRepoModel(d model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

// ManufacturerToRepoModel конвертирует данные производителя в модель репозитория
func ManufacturerToRepoModel(m model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// CategoryToRepoModel конвертирует категорию в модель репозитория
func CategoryToRepoModel(c model.Category) repoModel.Category {
	return repoModel.Category(c)
}

// CustomValueToRepoModel конвертирует CustomValue в модель репозитория
func CustomValueToRepoModel(cv model.CustomValue) repoModel.CustomValue {
	return repoModel.CustomValue{
		StringValue: cv.StringValue,
		Int64Value:  cv.Int64Value,
		DoubleValue: cv.DoubleValue,
		BoolValue:   cv.BoolValue,
	}
}

// MetadataToRepoModel выполняет маппинг словаря метаданных
func MetadataToRepoModel(meta map[string]model.CustomValue) map[string]repoModel.CustomValue {
	if meta == nil {
		return nil
	}

	res := make(map[string]repoModel.CustomValue, len(meta))
	for k, v := range meta {
		res[k] = CustomValueToRepoModel(v)
	}
	return res
}

// PartsFilterToRepoModel конвертирует фильтр выборки в модель репозитория
func PartsFilterToRepoModel(filter model.PartsFilter) repoModel.PartsFilter {
	var categories []repoModel.Category
	if filter.Categories != nil {
		categories = make([]repoModel.Category, len(filter.Categories))
		for i, cat := range filter.Categories {
			categories[i] = CategoryToRepoModel(cat)
		}
	}

	return repoModel.PartsFilter{
		UUIDs:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
