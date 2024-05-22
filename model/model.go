// model/models.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64         `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Version   int            `gorm:"default:1"`
}

type Type struct {
	BaseModel
	Namespace   string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Family      string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Name        string `gorm:"uniqueIndex:idx_namespace_family_name"`
	ElementType string
	WidgetType  string
}

type Validation struct {
	BaseModel
	Namespace        string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Family           string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Name             string `gorm:"uniqueIndex:idx_namespace_family_name"`
	RuleName         string
	ValidationParams string `gorm:"type:json"`
}

type Attribute struct {
	BaseModel
	Namespace  string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Family     string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Name       string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Label      string
	DesignSpec string `gorm:"type:json"`
}

type Form struct {
	BaseModel
	Namespace  string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Family     string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Name       string `gorm:"uniqueIndex:idx_namespace_family_name"`
	ActionName string
	Attributes string `gorm:"type:json"`
}
