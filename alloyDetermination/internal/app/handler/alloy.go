package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"alloyDetermination/internal/app/ds"
	"alloyDetermination/internal/app/repository"
)

const currentUserID uint = 1

type AlloyDTO struct {
	ds.Alloy
	LikesCount int64
}

// GetCatalogAlloy godoc
// @Summary      Каталог сплавов
// @Description  Опубликованные сплавы с фильтром по энергии
// @Tags         alloys
// @Produce      json
// @Param        energy query float64 false "Максимальная энергия, keV"
// @Success      200  {object}  map[string]interface{}
// @Failure      403
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /alloys [get]
func (h *Handler) GetCatalogAlloy(ctx *gin.Context) {
	energyStr := ctx.Query("energy")
	var energy float64
	if energyStr == "" {
		energy = 25.0
	} else {
		energy, _ = strconv.ParseFloat(energyStr, 64)
	}

	alloys, err := h.Repository.GetAlloysCatalog(energy)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	user := h.GetCurrentUser(ctx)

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
		"data": result,
	})
}

// GetFeedAlloy godoc
// @Summary      Лента сплавов
// @Tags         alloys
// @Produce      json
// @Param        alloy_id query int false "ID текущего сплава"
// @Param        next     query bool false "true — следующий, false — предыдущий"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /alloys/feed [get]
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

	user := h.GetCurrentUser(ctx)
	liked, _ := h.Repository.IsLikedByUser(user.ID, alloy.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"alloy":    alloy,
			"likes":    likes,
			"is_mine":  boolToInt(alloy.CreatorID == user.ID),
			"is_liked": boolToInt(liked),
		},
	})
}

// GetDraftAlloy godoc
// @Summary      Получить черновик текущего пользователя
// @Description  Возвращает черновик сплава, созданный текущим пользователем. У одного пользователя не может быть больше одного черновика.
// @Tags         alloys
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /alloys/draft [get]
func (h *Handler) GetDraftAlloy(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)

	draft, err := h.Repository.GetDraftByCreator(user.ID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if draft == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": draft,
	})
}

// PostDraftAlloy godoc
// @Summary      Создать черновик сплава
// @Tags         alloys
// @Accept       multipart/form-data
// @Produce      json
// @Param        name  formData string true  "Название сплава"
// @Param        image formData file   false "Картинка"
// @Param        video formData file   false "Видео"
// @Success      201  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /alloys [post]
func (h *Handler) PostDraftAlloy(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)

	// Проверяем, нет ли уже черновика
	if existing, _ := h.Repository.GetDraftByCreator(user.ID); existing != nil {
		ctx.JSON(http.StatusConflict, gin.H{
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
		"data": draft,
	})
}

// PutPublishAlloy godoc
// @Summary      Опубликовать сплав
// @Tags         alloys
// @Accept       json
// @Produce      json
// @Param        body body object true "Описание, energy_kev, intensity_cps"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /alloys/publish [put]
func (h *Handler) PutPublishAlloy(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req struct {
		Description  string  `json:"description"`
		EnergyKev    float64 `json:"energy_kev"`
		IntensityCps float64 `json:"intensity_cps"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user := h.GetCurrentUser(ctx)
	draft, err := h.Repository.GetDraftByCreator(user.ID)
	// ...
	if draft == nil || draft.ID != uint(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "not your draft"})
		return
	}

	if err := h.Repository.PublishAlloy(uint(id), req.Description, req.EnergyKev, req.IntensityCps); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "alloy published"})
}

// DeleteAlloy godoc
// @Summary      Удалить сплав (soft delete)
// @Tags         alloys
// @Produce      json
// @Param        id path int true "ID сплава"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /alloys/{id} [delete]
func (h *Handler) DeleteAlloy(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user := h.GetCurrentUser(ctx)

	if err := h.Repository.SoftDeleteAlloy(uint(id), user.ID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "alloy deleted",
	})
}

// PostLikeAlloy godoc
// @Summary      Поставить/снять лайк
// @Tags         alloys
// @Accept       json
// @Produce      json
// @Param        body body object true "alloy_id и like (0 или 1)"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /alloys/like [post]
func (h *Handler) PostLikeAlloy(ctx *gin.Context) {
	var req struct {
		AlloyID uint `json:"alloy_id" binding:"required"`
		Like    int  `json:"like"` // 0 или 1
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user := h.GetCurrentUser(ctx)

	likes, err := h.Repository.SetLike(req.AlloyID, user.ID, req.Like)
	if err != nil {
		if errors.Is(err, repository.ErrAlloyNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"alloy_id":    req.AlloyID,
			"likes_count": likes,
		},
	})
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
