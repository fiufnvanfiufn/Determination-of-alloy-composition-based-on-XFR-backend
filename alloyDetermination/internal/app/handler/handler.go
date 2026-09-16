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
	// GET
	router.GET("/catalog", h.GetCatalogAlloy)
	router.GET("/feed", h.GetFeedAlloy)
	router.GET("/draft", h.GetDraftAlloy)

	// POST
	router.POST("/draft", h.PostDraftAlloy)
	router.POST("/publish", h.PostPublishAlloy)
	router.POST("/delete", h.PostDeleteAlloy)

	// редирект с корня на каталог — удобно
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/catalog")
	})
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
		"alloy_status":      "error",
		"alloy_description": err.Error(),
	})
}
