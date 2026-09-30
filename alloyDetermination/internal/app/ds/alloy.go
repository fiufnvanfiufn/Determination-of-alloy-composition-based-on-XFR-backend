package ds

import (
	"database/sql"
	"time"
)

type Alloy struct {
	ID           uint         `gorm:"column:alloy_id;primaryKey" json:"alloy_id"`
	Name         string       `gorm:"column:alloy_name;type:varchar(100);not null" json:"alloy_name"`
	Description  string       `gorm:"column:alloy_description;type:varchar(200)" json:"alloy_description"`
	Status       string       `gorm:"column:alloy_status;type:varchar(15);not null;default:'черновик'" json:"alloy_status"`
	ImgURL       string       `gorm:"column:alloy_image_url;type:varchar(255)" json:"alloy_image_url"`
	VideoURL     string       `gorm:"column:alloy_video_url;type:varchar(255)" json:"alloy_video_url"`
	EnergyKev    float64      `gorm:"column:alloy_energy_kev;type:numeric(10,2)" json:"alloy_energy_kev"`
	IntensityCps float64      `gorm:"column:alloy_intensity_cps;type:numeric(12,2)" json:"alloy_intensity_cps"`
	DateCreate   time.Time    `gorm:"column:alloy_created_at;not null" json:"alloy_created_at"`
	DateFormed   sql.NullTime `gorm:"column:alloy_formed_at;default:null" json:"alloy_formed_at"`
	CreatorID    uint         `gorm:"column:creator_id;not null" json:"creator_id"`

	Creator User `gorm:"foreignKey:CreatorID;references:user_id" json:"-"`
}
