package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "alloyDetermination/docs"

	"alloyDetermination/internal/app/redis"
	"alloyDetermination/internal/app/repository"
	"alloyDetermination/internal/app/role"
)

type Handler struct {
	Repository *repository.Repository
	JWTSecret  string
	Redis      *redis.Client
}

func NewHandler(r *repository.Repository, jwtSecret string, redisClient *redis.Client) *Handler {
	return &Handler{Repository: r, JWTSecret: jwtSecret, Redis: redisClient}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")

	api.GET("/alloys", h.GetCatalogAlloy)
	api.GET("/alloys/feed", h.GetFeedAlloy)

	// Публичные (без токена):
	api.POST("/users/register", h.PostRegisterUser)
	api.POST("/users/login", h.PostLoginUser)

	// Защищённые (нужен JWT):
	auth := api.Group("")
	auth.Use(h.WithAuthCheck(role.Buyer, role.Manager, role.Admin))
	{
		auth.GET("/alloys/draft", h.GetDraftAlloy)
		auth.POST("/alloys", h.PostDraftAlloy)
		auth.PUT("/alloys/publish", h.PutPublishAlloy)
		auth.DELETE("/alloys/:id", h.DeleteAlloy)
		auth.POST("/alloys/like", h.PostLikeAlloy)
		auth.POST("/users/logout", h.PostLogoutUser)
	}

	// Только для админа (пример):
	admin := api.Group("")
	admin.Use(h.WithAuthCheck(role.Admin))
	{
		// admin-only endpoints
	}
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
		"description": err.Error(),
	})
}
