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

// PartToModel преобразует репозиторную модель Part в доменную модель.
func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		UUID:      part.UUID,
		Info:      PartInfoToModel(part.Info),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

// PartInfoToModel преобразует PartInfo, выполняя глубокую изоляцию срезов и мап.
func PartInfoToModel(info repoModel.PartInfo) model.PartInfo {
	var tags []string
	if info.Tags != nil {
		tags = make([]string, len(info.Tags))
		copy(tags, info.Tags)
	}

	metadata := make(map[string]model.CustomValue, len(info.Metadata))
	if info.Metadata != nil {
		for k, v := range info.Metadata {
			metadata[k] = CustomValueToModel(v)
		}
	}

	return model.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Category:      model.Category(info.Category),
		Dimensions:    DimensionsToModel(info.Dimensions),
		Manufacturer:  ManufacturerToModel(info.Manufacturer),
		Tags:          tags,
		Metadata:      metadata,
	}
}

// DimensionsToModel преобразует Dimensions в доменную модель.
func DimensionsToModel(dim repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

// ManufacturerToModel преобразует Manufacturer в доменную модель.
func ManufacturerToModel(m repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// CustomValueToModel преобразует CustomValue с копированием указателей.
func CustomValueToModel(cv repoModel.CustomValue) model.CustomValue {
	return model.CustomValue{
		StringValue: cv.StringValue,
		Int64Value:  cv.Int64Value,
		DoubleValue: cv.DoubleValue,
		BoolValue:   cv.BoolValue,
	}
}

// PartsToModel преобразует срез деталей в пакетном режиме с предвыделением памяти.
func PartsToModel(parts []repoModel.Part) []model.Part {
	if parts == nil {
		return nil
	}

	res := make([]model.Part, len(parts))
	for i, p := range parts {
		res[i] = PartToModel(p)
	}
	return res
}
