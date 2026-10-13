package repository

import (
	"alloyDetermination/internal/app/ds"
	"errors"

	"gorm.io/gorm"
)

func (r *Repository) Register(user *ds.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) CreateUser(login, password string) (*ds.User, error) {
	u := ds.User{Login: login, Password: password}
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	var u ds.User
	err := r.db.Where("user_login = ?", login).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // ← юзер не найден = не ошибка
	}
	if err != nil {
		return nil, err // ← реальная ошибка БД
	}
	return &u, nil
}

func (r *Repository) GetUserByID(id uint) (*ds.User, error) {
	var u ds.User
	err := r.db.Where("user_id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
