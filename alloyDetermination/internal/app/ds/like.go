package ds

type Like struct {
	ID      uint `gorm:"column:like_id;primaryKey"`
	UserID  uint `gorm:"column:user_id;not null;uniqueIndex:idx_user_alloy"`
	AlloyID uint `gorm:"column:alloy_id;not null;uniqueIndex:idx_user_alloy"`

	User  User  `gorm:"foreignKey:UserID;references:user_id"`
	Alloy Alloy `gorm:"foreignKey:AlloyID;references:alloy_id"`
}
