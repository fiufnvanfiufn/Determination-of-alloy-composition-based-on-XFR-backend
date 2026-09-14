package handler

import (
	"alloyDetermination/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// DTO (Data Transfer Object) для передачи длины массива лайков в шаблон
type AlloyDTO struct {
	repository.Alloy
	LikesCount int
}

func (h *Handler) GetCatalog(ctx *gin.Context) {
	energyQuery := ctx.Query("energy") // GET параметр фильтрации

	alloys, err := h.Repository.GetAlloysCatalog(energyQuery)
	if err != nil {
		logrus.Error(err)
	}

	var dtos []AlloyDTO
	for _, a := range alloys {
		dtos = append(dtos, AlloyDTO{Alloy: a, LikesCount: len(a.Likes)})
	}

	ctx.HTML(http.StatusOK, "catalog.html", gin.H{
		"alloys": dtos,
		"query":  energyQuery, // Передаем обратно, чтобы сохранить в input
	})
}

func (h *Handler) GetDraft(ctx *gin.Context) {
	draft, err := h.Repository.GetAlloyDraft()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "draft.html", gin.H{
		"alloy": draft,
	})
}

func (h *Handler) GetFeed(ctx *gin.Context) {
	idStr := ctx.Query("id")
	nextStr := ctx.Query("next")

	id, _ := strconv.Atoi(idStr)
	next := nextStr == "true"

	alloy, err := h.Repository.GetAlloyFeed(id, next)
	if err != nil {
		logrus.Error(err)
	}

	dto := AlloyDTO{Alloy: alloy, LikesCount: len(alloy.Likes)}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"alloy": dto,
	})
}
