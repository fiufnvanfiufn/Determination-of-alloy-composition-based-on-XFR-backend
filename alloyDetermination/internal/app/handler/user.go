package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostRegisterUser(ctx *gin.Context) {
	var req struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.CreateUser(req.Login, req.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id": user.ID,
			"login":   user.Login,
		},
	})
}

// Заглушка для 4-й лабы — сейчас возвращает фиксированного singleton-пользователя
func (h *Handler) PostLoginUser(ctx *gin.Context) {
	user := h.GetCurrentUser()
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    gin.H{"user_id": user.ID, "login": user.Login},
		"message": "stub for lab 4",
	})
}

// Заглушка — ничего не делает
func (h *Handler) PostLogoutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "stub for lab 4",
	})
}
