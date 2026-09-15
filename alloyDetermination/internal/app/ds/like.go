package ds

type Like struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint `gorm:"not null"`
	AlloyID uint `gorm:"not null"`

	User  User  `gorm:"foreignKey:UserID"`
	Alloy Alloy `gorm:"foreignKey:AlloyID"`
}
