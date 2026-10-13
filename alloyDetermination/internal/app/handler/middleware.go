package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"

	"alloyDetermination/internal/app/ds"
	"alloyDetermination/internal/app/role"
)

const jwtPrefix = "Bearer "

// WithAuthCheck — middleware для проверки JWT.
//
// Принимает список ролей, которым разрешён доступ.
// Если список пуст — доступ только проверяется (без ролевой модели).
func (h *Handler) WithAuthCheck(assignedRoles ...role.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Достаём заголовок Authorization
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, jwtPrefix) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 2. Отрезаем префикс "Bearer "
		tokenStr := authHeader[len(jwtPrefix):]
		logrus.Infof("WithAuthCheck: tokenLen=%d", len(tokenStr))

		// 2.5. Проверяем, не в блэклисте ли токен (logout)
		err := h.Redis.CheckJWTInBlacklist(ctx.Request.Context(), tokenStr)
		logrus.Infof("CheckJWTInBlacklist: err=%v, isNil=%v", err, errors.Is(err, redis.Nil))
		if err == nil {
			// токен найден в блэклисте → разлогинен
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !errors.Is(err, redis.Nil) {
			// не nil-ошибка = проблема с Redis
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// 3. Парсим и валидируем JWT
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&ds.JWTClaims{},
			func(t *jwt.Token) (interface{}, error) {
				// Проверяем, что алгоритм — HMAC
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(h.JWTSecret), nil
			},
		)
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		claims, ok := token.Claims.(*ds.JWTClaims)
		if !ok || !token.Valid {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 4. Проверяем роль (если роли переданы)
		if len(assignedRoles) > 0 {
			roleOK := false
			for _, r := range assignedRoles {
				if claims.Role == r {
					roleOK = true
					break
				}
			}
			if !roleOK {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		// 5. Кладём данные из токена в контекст — их потом достанет GetCurrentUser
		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}
