package middleware

import (
	"os"
	"time"
	"errors"
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/Hailemari/auth_task_manager/models"
)

func GenerateToken(user *models.User) (string, error){
	claims := jwt.MapClaims {
		"id": user.ID,
		"username": user.Username,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return 
		}
		token, err := ValidateToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return 
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return 
		}

		ctx.Set("user", claims)
		ctx.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			ctx.Abort()
			return
		}

		claims, ok := user.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
			ctx.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			ctx.Abort()
			return 
		}

		ctx.Next()
	}
}