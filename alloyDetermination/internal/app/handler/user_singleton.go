package handler

import (
	"alloyDetermination/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCurrentUser(ctx *gin.Context) *ds.User {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return &ds.User{ID: 0}
	}
	return &ds.User{ID: userID.(uint)}
}
