package ds

type User struct {
	ID       uint   `gorm:"column:user_id;primaryKey"`
	Login    string `gorm:"column:user_login;type:varchar(25);unique;not null"`
	Password string `gorm:"column:user_password;type:varchar(100);not null"`
}
