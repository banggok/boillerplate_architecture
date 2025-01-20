package authentication_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/authentication"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/recovery"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	app.AppConfig.SecretKey.Access = "access_secret"
	app.AppConfig.SecretKey.Refresh = "refresh_secret"
	app.AppConfig.SecretKey.AccessExpired = time.Minute * 15
	app.AppConfig.SecretKey.RefreshExpired = time.Hour * 24
}

func TestValidateToken(t *testing.T) {
	validAccessToken, _, _ := authentication.GenerateTokens(1)

	t.Run("Valid Access Token", func(t *testing.T) {
		claims, err := authentication.ValidateToken(validAccessToken, true)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, uint(1), claims.UserID)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := authentication.ValidateToken("invalid_token", true)
		assert.Error(t, err)
	})

	t.Run("Empty Token", func(t *testing.T) {
		_, err := authentication.ValidateToken("", true)
		assert.Error(t, err)
	})
}

func TestGenerateTokens(t *testing.T) {
	accessToken, refreshToken, err := authentication.GenerateTokens(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}

func TestAuthMiddleware(t *testing.T) {
	validAccessToken, _, _ := authentication.GenerateTokens(1)
	router := gin.New()
	router.Use(recovery.CustomRecoveryMiddleware())
	router.Use(authentication.AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validAccessToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRefreshAuthMiddleware(t *testing.T) {
	_, validRefreshToken, _ := authentication.GenerateTokens(1)
	router := gin.New()
	router.Use(recovery.CustomRecoveryMiddleware())
	router.Use(authentication.RefreshAuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validRefreshToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGetUserId(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set(authentication.USERID, uint(1))

	t.Run("Valid User ID", func(t *testing.T) {
		userID, err := authentication.GetUserId(ctx)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), *userID)
	})

	t.Run("Missing User ID", func(t *testing.T) {
		newCtx := &gin.Context{}
		_, err := authentication.GetUserId(newCtx)
		assert.Error(t, err)
	})
}
