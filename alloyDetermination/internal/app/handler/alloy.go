package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"alloyDetermination/internal/app/ds"
)

// Хардкод пользователя, пока нет авторизации (JWT появится позже).
const currentUserID uint = 1

// DTO для передачи в шаблон — расширяет модель вычисляемыми полями.
type AlloyDTO struct {
	ds.Alloy
	LikesCount int64
}

// ---------- GET ----------

// /catalog — список опубликованных услуг с поиском.
func (h *Handler) GetCatalog(ctx *gin.Context) {
	search := ctx.Query("search")

	alloys, err := h.Repository.GetPublishedAlloys(search)
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
		"search": search,
	})
}

// /feed — лента (плитка с "next"/"prev").
func (h *Handler) GetFeed(ctx *gin.Context) {
	idStr := ctx.Query("id")
	nextStr := ctx.Query("next")
	id, _ := strconv.Atoi(idStr)
	next := nextStr == "true"

	alloy, err := h.Repository.GetAlloyFeed(id, next)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	likes, _ := h.Repository.GetLikesCount(alloy.ID)

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"alloy": AlloyDTO{Alloy: *alloy, LikesCount: likes},
	})
}

// /draft — страница добавления/публикации.
// Если у пользователя есть черновик — рендерим с кнопкой "Опубликовать",
// иначе — с пустыми полями и кнопкой "Далее".
func (h *Handler) GetDraft(ctx *gin.Context) {
	draft, err := h.Repository.GetDraftByCreator(currentUserID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "draft.html", gin.H{
		"alloy": draft, // nil, если черновика нет
	})
}

// ---------- POST ----------

// POST /draft — создание черновика (ORM, кнопка "Далее").
func (h *Handler) PostDraft(ctx *gin.Context) {
	// Защита: у пользователя уже есть черновик — открываем его же.
	if existing, _ := h.Repository.GetDraftByCreator(currentUserID); existing != nil {
		ctx.Redirect(http.StatusFound, "/draft")
		return
	}

	name := ctx.PostForm("name")
	imgURL := ctx.PostForm("img_url") // файлы в этой лабе не сохраняем
	videoURL := ctx.PostForm("video_url")

	if _, err := h.Repository.CreateDraft(currentUserID, name, imgURL, videoURL); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/draft")
}

// POST /publish — публикация (ORM, кнопка "Опубликовать").
func (h *Handler) PostPublish(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	description := ctx.PostForm("description")
	energyKev, _ := strconv.ParseFloat(ctx.PostForm("energy_kev"), 64)
	intensityCps, _ := strconv.ParseFloat(ctx.PostForm("intensity_cps"), 64)

	if err := h.Repository.PublishDraft(uint(id), description, energyKev, intensityCps); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/catalog")
}

// POST /delete — логическое удаление (raw SQL, БЕЗ ORM).
func (h *Handler) PostDelete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("alloy_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.SoftDeleteAlloy(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/catalog")
}
