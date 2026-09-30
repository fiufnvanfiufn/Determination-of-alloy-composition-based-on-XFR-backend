package repository

import (
	"alloyDetermination/internal/app/ds"
)

// CreateUser создаёт нового пользователя.
func (r *Repository) CreateUser(login, password string) (*ds.User, error) {
	u := ds.User{Login: login, Password: password}
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByLogin ищет пользователя по логину.
func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	var u ds.User
	err := r.db.Where("user_login = ?", login).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
