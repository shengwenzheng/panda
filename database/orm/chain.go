package orm

import "gorm.io/gorm"

type Chain struct {
	gorm.Model
	Height int64 `gorm:"not null"`
}