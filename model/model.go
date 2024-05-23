package model

import (
	"time"

	"gorm.io/gorm"
)

type Type struct {
	ID          uint64 `gorm:"primaryKey"`
	Namespace   string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Family      string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Name        string `gorm:"uniqueIndex:idx_namespace_family_name"`
	ElementType string
	WidgetType  string
	CreatedAt   time.Time      `gorm:"autoCreateTime:milli"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:milli"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Version     int            `gorm:"default:1"`
}

type Validation struct {
	ID               uint64 `gorm:"primaryKey"`
	Namespace        string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Family           string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Name             string `gorm:"uniqueIndex:idx_namespace_family_name"`
	RuleName         string
	ValidationParams string
	CreatedAt        time.Time      `gorm:"autoCreateTime:milli"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime:milli"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Version          int            `gorm:"default:1"`
	Attributes       []Attribute    `gorm:"many2many:attribute_validations;"`
}

type Attribute struct {
	ID          uint64 `gorm:"primaryKey"`
	Namespace   string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Family      string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Name        string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Label       string
	DesignSpec  string `gorm:"type:json"`
	TypeID      uint64
	CreatedAt   time.Time      `gorm:"autoCreateTime:milli"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:milli"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Version     int            `gorm:"default:1"`
	Type        Type
	Validations []Validation `gorm:"many2many:attribute_validations;"`
}

type Form struct {
	ID         uint64 `gorm:"primaryKey"`
	Namespace  string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Family     string `gorm:"uniqueIndex:idx_namespace_family_name"`
	Name       string `gorm:"uniqueIndex:idx_namespace_family_name"`
	ActionName string
	CreatedAt  time.Time      `gorm:"autoCreateTime:milli"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime:milli"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Version    int            `gorm:"default:1"`
	Attributes []Attribute    `gorm:"many2many:form_attributes;"`
}
