package ds

import (
	"database/sql"
	"time"
)

type Alloy struct {
	ID           uint         `gorm:"column:alloy_id;primaryKey"`
	Name         string       `gorm:"column:alloy_name;type:varchar(100);not null"`
	Description  string       `gorm:"column:alloy_description;type:varchar(200)"`
	Status       string       `gorm:"column:alloy_status;type:varchar(15);not null;default:'черновик'"`
	ImgURL       string       `gorm:"column:alloy_image_url;type:varchar(255)"`
	VideoURL     string       `gorm:"column:alloy_video_url;type:varchar(255)"`
	EnergyKev    float64      `gorm:"column:alloy_energy_kev;type:numeric(10,2)"`
	IntensityCps float64      `gorm:"column:alloy_intensity_cps;type:numeric(12,2)"`
	DateCreate   time.Time    `gorm:"column:alloy_created_at;not null"`
	DateFormed   sql.NullTime `gorm:"column:alloy_formed_at;default:null"`
	CreatorID    uint         `gorm:"column:creator_id;not null"`

	Creator User `gorm:"foreignKey:CreatorID;references:user_id"`
}
