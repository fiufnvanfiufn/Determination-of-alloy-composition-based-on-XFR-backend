package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"alloyDetermination/internal/app/ds"
)

const currentUserID uint = 1

type AlloyDTO struct {
	ds.Alloy
	LikesCount int64
}

func (h *Handler) GetCatalogAlloy(ctx *gin.Context) {
	energyStr := ctx.Query("energy")
	var energy float64
	if energyStr == "" {
		energy = 24.0
	} else {
		energy, _ = strconv.ParseFloat(energyStr, 64)
	}

	alloys, err := h.Repository.GetAlloysCatalog(energy)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	user := h.GetCurrentUser()

	type AlloyResponse struct {
		ds.Alloy
		IsMine bool `json:"is_mine"`
	}

	result := make([]AlloyResponse, 0, len(alloys))
	for _, a := range alloys {
		result = append(result, AlloyResponse{
			Alloy:  a,
			IsMine: a.CreatorID == user.ID,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
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

	user := h.GetCurrentUser()

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"alloy":   alloy,
			"likes":   likes,
			"is_mine": alloy.CreatorID == user.ID,
		},
	})
}

func (h *Handler) GetDraftAlloy(ctx *gin.Context) {
	user := h.GetCurrentUser()

	draft, err := h.Repository.GetDraftByCreator(user.ID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if draft == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   draft,
	})
}

func (h *Handler) PostDraftAlloy(ctx *gin.Context) {
	user := h.GetCurrentUser()

	// Проверяем, нет ли уже черновика
	if existing, _ := h.Repository.GetDraftByCreator(user.ID); existing != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "draft already exists",
		})
		return
	}

	// Читаем форму
	name := ctx.PostForm("name")
	if name == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("name is required"))
		return
	}

	// Создаём черновик
	draft, err := h.Repository.CreateDraft(user.ID, name, "", "")
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	imgName, videoName := "", ""

	// Картинка
	if file, err := ctx.FormFile("image"); err == nil {
		imgName = fmt.Sprintf("alloy_%d_image.png", draft.ID)
		if err := h.Repository.UploadFile(file, imgName); err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Видео
	if file, err := ctx.FormFile("video"); err == nil {
		videoName = fmt.Sprintf("alloy_%d_video.mp4", draft.ID)
		if err := h.Repository.UploadFile(file, videoName); err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Сохраняем имена файлов в БД
	if imgName != "" || videoName != "" {
		_ = h.Repository.UpdateAlloyFiles(draft.ID, imgName, videoName)
	}

	// Перечитываем с обновлёнными полями
	draft, _ = h.Repository.GetAlloyByID(draft.ID)

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   draft,
	})
}
func (h *Handler) PutPublishAlloy(ctx *gin.Context) {
	var req struct {
		AlloyID      uint    `json:"alloy_id" binding:"required"`
		Description  string  `json:"description"`
		EnergyKev    float64 `json:"energy_kev"`
		IntensityCps float64 `json:"intensity_cps"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем, что это черновик текущего пользователя
	user := h.GetCurrentUser()
	draft, err := h.Repository.GetDraftByCreator(user.ID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if draft == nil || draft.ID != req.AlloyID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "not your draft",
		})
		return
	}

	if err := h.Repository.PublishAlloy(req.AlloyID, req.Description, req.EnergyKev, req.IntensityCps); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "alloy published",
	})
}

func (h *Handler) DeleteAlloy(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user := h.GetCurrentUser()

	if err := h.Repository.SoftDeleteAlloy(uint(id), user.ID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "alloy deleted",
	})
}

func (h *Handler) PostLikeAlloy(ctx *gin.Context) {
	var req struct {
		AlloyID uint `json:"alloy_id" binding:"required"`
		Like    int  `json:"like"` // 0 или 1
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user := h.GetCurrentUser()

	if err := h.Repository.SetLike(user.ID, req.AlloyID, req.Like == 1); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	likes, _ := h.Repository.GetLikesCount(req.AlloyID)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"alloy_id":    req.AlloyID,
			"likes_count": likes,
		},
	})
}
