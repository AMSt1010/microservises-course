package converter

import (
	"time"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	InvV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ==========================================
// Proto -> Model
// ==========================================

// PartToModel преобразует protobuf Part в доменную модель Part.
func PartToModel(part *InvV1.Part) model.Part {
	if part == nil {
		return model.Part{}
	}

	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		updatedAt = lo.ToPtr(part.UpdatedAt.AsTime())
	}

	var createdAt time.Time
	if part.CreatedAt != nil {
		createdAt = part.CreatedAt.AsTime()
	}

	return model.Part{
		UUID:      part.Uuid,
		Info:      PartInfoToModel(part.Info),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// PartInfoToModel преобразует protobuf PartInfo в доменную модель PartInfo.
func PartInfoToModel(info *InvV1.Info) model.PartInfo {
	if info == nil {
		return model.PartInfo{}
	}

	var tags []string
	if info.Tags != nil {
		tags = make([]string, len(info.Tags))
		copy(tags, info.Tags)
	}

	var metadata map[string]model.CustomValue
	if info.Metadata != nil {
		metadata = make(map[string]model.CustomValue, len(info.Metadata))
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

// DimensionsToModel преобразует protobuf Dimensions в доменную модель.
func DimensionsToModel(dim *InvV1.Dimensions) model.Dimensions {
	if dim == nil {
		return model.Dimensions{}
	}

	return model.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

// ManufacturerToModel преобразует protobuf Manufacturer в доменную модель.
func ManufacturerToModel(m *InvV1.Manufacturer) model.Manufacturer {
	if m == nil {
		return model.Manufacturer{}
	}

	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// CustomValueToModel преобразует protobuf CustomValue (oneof) в доменную модель.
func CustomValueToModel(cv *InvV1.CustomValue) model.CustomValue {
	if cv == nil {
		return model.CustomValue{}
	}

	var res model.CustomValue

	switch val := cv.Kind.(type) {
	case *InvV1.CustomValue_StringValue:
		res.StringValue = lo.ToPtr(val.StringValue)
	case *InvV1.CustomValue_Int64Value:
		res.Int64Value = lo.ToPtr(val.Int64Value)
	case *InvV1.CustomValue_DoubleValue:
		res.DoubleValue = lo.ToPtr(val.DoubleValue)
	case *InvV1.CustomValue_BoolValue:
		res.BoolValue = lo.ToPtr(val.BoolValue)
	}

	return res
}

// PartsFilterToModel преобразует protobuf PartsFilter в доменную модель.
func PartsFilterToModel(filter *InvV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	var categories []model.Category
	if filter.Categories != nil {
		categories = lo.Map(filter.Categories, func(item InvV1.Category, _ int) model.Category {
			return model.Category(item)
		})
	}
	/*
			var categories []model.Category
			if filter.Categories != nil {
		    	// 1. Выделяем память ровно под нужный объем за одну аллокацию
		    	categories = make([]model.Category, len(filter.Categories))

		    	// 2. Явно приводим каждый элемент по индексу
		    	for i, item := range filter.Categories {
		        categories[i] = model.Category(item)
		    	}
			}
	*/

	var uuids []string
	if filter.Uuids != nil {
		uuids = make([]string, len(filter.Uuids))
		copy(uuids, filter.Uuids)
	}

	var names []string
	if filter.Names != nil {
		names = make([]string, len(filter.Names))
		copy(names, filter.Names)
	}

	var countries []string
	if filter.ManufacturerCountries != nil {
		countries = make([]string, len(filter.ManufacturerCountries))
		copy(countries, filter.ManufacturerCountries)
	}

	var tags []string
	if filter.Tags != nil {
		tags = make([]string, len(filter.Tags))
		copy(tags, filter.Tags)
	}

	return model.PartsFilter{
		UUIDs:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: countries,
		Tags:                  tags,
	}
}

// ==========================================
// Model -> Proto
// ==========================================

// PartToProto преобразует доменную модель Part в protobuf Part.
func PartToProto(part model.Part) *InvV1.Part {
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &InvV1.Part{
		Uuid:      part.UUID,
		Info:      PartInfoToProto(part.Info),
		CreatedAt: timestamppb.New(part.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// PartInfoToProto преобразует доменную модель PartInfo в protobuf PartInfo.
func PartInfoToProto(info model.PartInfo) *InvV1.Info {
	var tags []string
	if info.Tags != nil {
		tags = make([]string, len(info.Tags))
		copy(tags, info.Tags)
	}

	var metadata map[string]*InvV1.CustomValue
	if info.Metadata != nil {
		metadata = make(map[string]*InvV1.CustomValue, len(info.Metadata))
		for k, v := range info.Metadata {
			metadata[k] = CustomValueToProto(v)
		}
	}

	return &InvV1.Info{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Category:      InvV1.Category(info.Category),
		Dimensions:    DimensionsToProto(info.Dimensions),
		Manufacturer:  ManufacturerToProto(info.Manufacturer),
		Tags:          tags,
		Metadata:      metadata,
	}
}

// DimensionsToProto преобразует Dimensions в protobuf Dimensions.
func DimensionsToProto(dim model.Dimensions) *InvV1.Dimensions {
	return &InvV1.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

// ManufacturerToProto преобразует Manufacturer в protobuf Manufacturer.
func ManufacturerToProto(m model.Manufacturer) *InvV1.Manufacturer {
	return &InvV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// CustomValueToProto преобразует доменное CustomValue в protobuf CustomValue.
func CustomValueToProto(cv model.CustomValue) *InvV1.CustomValue {
	switch {
	case cv.StringValue != nil:
		return &InvV1.CustomValue{
			Kind: &InvV1.CustomValue_StringValue{
				StringValue: *cv.StringValue,
			},
		}
	case cv.Int64Value != nil:
		return &InvV1.CustomValue{
			Kind: &InvV1.CustomValue_Int64Value{
				Int64Value: *cv.Int64Value,
			},
		}
	case cv.DoubleValue != nil:
		return &InvV1.CustomValue{
			Kind: &InvV1.CustomValue_DoubleValue{
				DoubleValue: *cv.DoubleValue,
			},
		}
	case cv.BoolValue != nil:
		return &InvV1.CustomValue{
			Kind: &InvV1.CustomValue_BoolValue{
				BoolValue: *cv.BoolValue,
			},
		}
	default:
		return nil
	}
}

// PartsToProto преобразует срез доменных моделей Part в срез protobuf Part с помощью lo.Map.
func PartsToProto(parts []model.Part) []*InvV1.Part {
	if parts == nil {
		return nil
	}

	return lo.Map(parts, func(item model.Part, _ int) *InvV1.Part {
		return PartToProto(item)
	})
}
