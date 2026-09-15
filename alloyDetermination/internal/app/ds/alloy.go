package ds

import (
	"database/sql"
	"time"
)

type Alloy struct {
	ID           uint         `gorm:"primaryKey"`
	Name         string       `gorm:"type:varchar(100);not null"` // наименование
	Description  string       `gorm:"type:varchar(200)"`          // краткое описание
	Status       string       `gorm:"type:varchar(15);not null;default:'черновик'"`
	ImgURL       string       `gorm:"type:varchar(255)"`
	VideoURL     string       `gorm:"type:varchar(255)"`
	EnergyKev    float64      `gorm:"type:numeric(10,2)"` // поле 1 по теме
	IntensityCps float64      `gorm:"type:numeric(12,2)"` // поле 2 по теме
	DateCreate   time.Time    `gorm:"not null"`
	DateFormed   sql.NullTime `gorm:"default:null"`
	CreatorID    uint         `gorm:"not null"`

	Creator User `gorm:"foreignKey:CreatorID"`
}
