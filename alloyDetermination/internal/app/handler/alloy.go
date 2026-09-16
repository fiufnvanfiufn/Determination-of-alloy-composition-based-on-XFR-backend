package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"alloyDetermination/internal/app/ds"
)

const currentUserID uint = 1

type AlloyDTO struct {
	ds.Alloy
	LikesCount int64
}

func (h *Handler) GetCatalogAlloy(ctx *gin.Context) {
	logrus.Infof("GetCatalog RAW QUERY: %q", ctx.Request.URL.RawQuery)
	energyStr := ctx.Query("energy")
	var energy float64
	if energyStr != "" {
		energy, _ = strconv.ParseFloat(energyStr, 64)
	}
	alloys, err := h.Repository.GetAlloysCatalog(energy)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ids := make([]uint, 0, len(alloys))
	for _, a := range alloys {
		ids = append(ids, a.ID)
	}
	likesMap, err := h.Repository.GetLikesCountsForAlloys(ids)
	if err != nil {
		logrus.Error(err)
		likesMap = map[uint]int64{}
	}
	dtos := make([]AlloyDTO, 0, len(alloys))
	for _, a := range alloys {
		dtos = append(dtos, AlloyDTO{Alloy: a, LikesCount: likesMap[a.ID]})
	}
	ctx.HTML(http.StatusOK, "catalog.html", gin.H{
		"alloys": dtos,
		"energy": energy,
	})
}

func (h *Handler) GetFeedAlloy(ctx *gin.Context) {
	idStr := ctx.Query("alloy_id")
	nextStr := ctx.Query("next")
	id, _ := strconv.Atoi(idStr)
	next := nextStr == "true"

	alloy, err := h.Repository.GetAlloyFeed(uint(id), next)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	likes, _ := h.Repository.GetLikesCount(alloy.ID)

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"alloy": AlloyDTO{Alloy: *alloy, LikesCount: likes},
	})
}

func (h *Handler) GetDraftAlloy(ctx *gin.Context) {
	draft, err := h.Repository.GetDraftByCreator(currentUserID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "draft.html", gin.H{
		"alloy": draft,
	})
}

func (h *Handler) PostDraftAlloy(ctx *gin.Context) {
	if existing, _ := h.Repository.GetDraftByCreator(currentUserID); existing != nil {
		ctx.Redirect(http.StatusFound, "/draft")
		return
	}

	name := ctx.PostForm("name")
	imgURL := ctx.PostForm("img_url")
	videoURL := ctx.PostForm("video_url")

	if _, err := h.Repository.CreateDraft(currentUserID, name, imgURL, videoURL); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/draft")
}

func (h *Handler) PostPublishAlloy(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	description := ctx.PostForm("description")
	energyKev, _ := strconv.ParseFloat(ctx.PostForm("energy_kev"), 64)
	intensityCps, _ := strconv.ParseFloat(ctx.PostForm("intensity_cps"), 64)

	if err := h.Repository.PublishAlloy(uint(id), description, energyKev, intensityCps); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/catalog")
}

func (h *Handler) PostDeleteAlloy(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("alloy_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.Repository.DeleteAlloySQL(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/catalog")
}
