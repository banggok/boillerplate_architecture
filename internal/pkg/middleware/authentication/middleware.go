package authentication

import (
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	accessTokenSecret  = []byte(app.AppConfig.SecretKey.Access)
	refreshTokenSecret = []byte(app.AppConfig.SecretKey.Refresh)
)

const USERID = "userID"

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Function to validate JWT and extract claims
func ValidateToken(tokenString string, isAccessToken bool) (*Claims, error) {
	if tokenString == "" {
		return nil, custom_errors.New(nil, custom_errors.Unauthorized, "missing token")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if isAccessToken {
			return accessTokenSecret, nil
		}
		return refreshTokenSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, custom_errors.New(nil, custom_errors.Unauthorized, "invalid token")
	}

	return claims, nil
}

// Generate JWT Tokens
func GenerateTokens(userID uint) (string, string, error) {
	// Create Access Token
	accessTokenClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(app.AppConfig.SecretKey.AccessExpired)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(accessTokenSecret)
	if err != nil {
		return "", "", err
	}

	// Create Refresh Token
	refreshTokenClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(app.AppConfig.SecretKey.RefreshExpired)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshTokenSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// JWT Authentication Middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:] // Remove "Bearer " prefix
		}

		// Validate the token
		claims, err := ValidateToken(tokenString, true)
		if err != nil {
			c.Error(custom_errors.New(
				err,
				custom_errors.Unauthorized,
				"Invalid or missing authorization token",
			))
			c.Abort() // Stop further processing
			return
		}

		// Store the user ID in the context
		c.Set(USERID, claims.UserID)

		// Proceed to the next handler
		c.Next()
	}
}

func GetUserId(c *gin.Context) (*uint, error) {
	// Retrieve the value from the context
	userID, exists := c.Get(USERID)
	if !exists {
		return nil, custom_errors.New(nil, custom_errors.InternalServerError, "failed to get userID from context: key not found")
	}

	id, ok := userID.(uint)
	if !ok {
		return nil, custom_errors.New(nil, custom_errors.InternalServerError, "failed to get userID from context: invalid type")
	}

	return &id, nil
}
