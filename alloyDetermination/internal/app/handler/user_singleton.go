package handler

import (
	"alloyDetermination/internal/app/ds"
)

var currentUser *ds.User

func (h *Handler) GetCurrentUser() *ds.User {
	if currentUser == nil {
		currentUser = &ds.User{
			ID:    1,
			Login: "test",
		}
	}
	return currentUser
}
