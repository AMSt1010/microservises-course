package model

import (
	"errors"
	"time"
)

// PartInfo содержит базовые характеристики детали.
type PartInfo struct {
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]CustomValue
}

// Dimensions описывает физические размеры и вес.
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

// Manufacturer описывает данные производителя.
type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// Part описывает сущность детали в бизнес-логике.
type Part struct {
	UUID      string
	Info      PartInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// PartsFilter задает параметры фильтрации при выборке деталей.
type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

// Category описывает доступные категории детали.
type Category int32

const (
	CategoryUnspecified Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

// CustomValue представляет строго одно из допустимых значений метаданных.
type CustomValue struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

func NewStringValue(v string) CustomValue {
	return CustomValue{StringValue: &v}
}

func NewInt64Value(v int64) CustomValue {
	return CustomValue{Int64Value: &v}
}

func NewDoubleValue(v float64) CustomValue {
	return CustomValue{DoubleValue: &v}
}

func NewBoolValue(v bool) CustomValue {
	return CustomValue{BoolValue: &v}
}

func (c CustomValue) Validate() error {
	count := 0
	if c.StringValue != nil {
		count++
	}
	if c.Int64Value != nil {
		count++
	}
	if c.DoubleValue != nil {
		count++
	}
	if c.BoolValue != nil {
		count++
	}

	if count != 1 {
		return errors.New("custom value must have exactly one field populated")
	}
	return nil
}
