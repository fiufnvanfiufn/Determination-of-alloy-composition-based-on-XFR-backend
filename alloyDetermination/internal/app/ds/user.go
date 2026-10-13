package ds

import "alloyDetermination/internal/app/role"

type User struct {
	ID       uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	Login    string    `gorm:"column:user_login;unique;not null" json:"login"`
	Password string    `gorm:"column:user_password;not null" json:"-"`
	Role     role.Role `gorm:"column:role;default:0" json:"role"`
}
