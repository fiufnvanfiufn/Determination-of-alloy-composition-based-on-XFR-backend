package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"alloyDetermination/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")

	// Каталог и лента
	api.GET("/alloys", h.GetCatalogAlloy)
	api.GET("/alloys/feed", h.GetFeedAlloy)

	// Черновик
	api.GET("/alloys/draft", h.GetDraftAlloy)
	api.POST("/alloy/:id/draft", h.PostDraftAlloy)

	// Публикация
	api.PUT("/alloys/:id/publish", h.PutPublishAlloy)

	// Удаление и лайк
	api.DELETE("/alloys/:id", h.DeleteAlloy)
	api.POST("/alloys/:id/like", h.PostLikeAlloy)

	// Пользователь
	api.POST("/user/register", h.PostRegisterUser)
	api.POST("/user/login", h.PostLoginUser)
	api.POST("/user/logout", h.PostLogoutUser)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
	router.Static("/media", "./resources/media")
}

func (h *Handler) errorHandler(ctx *gin.Context, code int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(code, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
