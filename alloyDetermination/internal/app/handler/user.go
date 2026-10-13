package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"alloyDetermination/internal/app/ds"
	"alloyDetermination/internal/app/role"
)

// @title Alloy Determination API
// @version 1.0
// @description API для определения состава сплавов
// @host localhost:3000
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <JWT>
type registerReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PostRegisterUser godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя с ролью Buyer
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body registerReq true "Логин и пароль"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/register [post]
func (h *Handler) PostRegisterUser(ctx *gin.Context) {
	var req registerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем, что логин не занят
	existing, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if existing != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "login already taken",
		})
		return
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Создаём пользователя с ролью Buyer
	user := &ds.User{
		Login:    req.Login,
		Password: string(hash),
		Role:     role.Buyer,
	}
	if err := h.Repository.Register(user); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"user_id": user.ID,
			"login":   user.Login,
			"role":    user.Role,
		},
	})
}

type loginReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResp struct {
	ExpiresIn   int64  `json:"expires_in"` // секунды
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// PostLoginUser godoc
// @Summary      Логин пользователя
// @Description  Проверяет логин/пароль, возвращает JWT access_token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body loginReq true "Логин и пароль"
// @Success      200  {object}  loginResp
// @Failure      400  {object}  map[string]interface{}
// @Failure      403
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/login [post]
func (h *Handler) PostLoginUser(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// 1. Ищем пользователя
	user, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// 2. Сравниваем bcrypt-хеш
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// 3. Генерируем JWT
	expiresIn := 24 * time.Hour // или из .env
	claims := &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "alloy-determination",
		},
		UserID: user.ID,
		Role:   user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, loginResp{
		ExpiresIn:   int64(expiresIn.Seconds()),
		AccessToken: tokenStr,
		TokenType:   "Bearer",
	})
}

// PostLogoutUser godoc
// @Summary      Выход из системы
// @Description  Кладёт токен в blacklist Redis
// @Tags         users
// @Produce      json
// @Param        Authorization header string true "Bearer <token>"
// @Success      200  {object}  map[string]interface{}
// @Failure      400
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /users/logout [post]
func (h *Handler) PostLogoutUser(ctx *gin.Context) {
	// 1. Достаём токен из заголовка
	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) < len(jwtPrefix) || authHeader[:len(jwtPrefix)] != jwtPrefix {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	tokenStr := authHeader[len(jwtPrefix):]

	// 2. Парсим, чтобы узнать TTL (сколько токену ещё жить)
	token, err := jwt.ParseWithClaims(tokenStr, &ds.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.JWTSecret), nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*ds.JWTClaims)
	if !ok || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TTL = сколько секунд токену ещё жить
	ttl := time.Until(time.Unix(claims.ExpiresAt, 0))
	if ttl <= 0 {
		ctx.Status(http.StatusOK)
		return
	}
	logrus.Infof("Logout: ttl=%v, ExpiresAt=%v, now=%v", ttl, claims.ExpiresAt, time.Now().Unix())
	// 3. Пишем в Redis
	if err := h.Redis.WriteJWTToBlacklist(ctx.Request.Context(), tokenStr, ttl); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	logrus.Infof("Logout: wrote token to blacklist, ttl=%v", ttl)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}
