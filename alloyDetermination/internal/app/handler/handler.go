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
	api.GET("/alloyCatalog", h.GetCatalogAlloy)
	api.GET("/alloyFeed", h.GetFeedAlloy)

	// Черновик
	api.GET("/alloyDraft", h.GetDraftAlloy)
	api.POST("/alloyDraft", h.PostDraftAlloy) // ← убрать router., поставить api.

	// Публикация
	api.PUT("/alloyPublish", h.PutPublishAlloy) // ← тоже

	// Удаление и лайк
	api.DELETE("/alloyDelete/:id", h.DeleteAlloy)
	api.POST("/alloyLike", h.PostLikeAlloy)

	// Пользователь
	api.POST("/userRegister", h.PostRegisterUser)
	api.POST("/userLogin", h.PostLoginUser)
	api.POST("/userLogout", h.PostLogoutUser)
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
