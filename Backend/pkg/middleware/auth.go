package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your_jwt_secret_key")

// RequestLogger logs incoming requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		logMessage := fmt.Sprintf("[%s] %s %d %s %s %s %s",
			timestamp.Format("2006-01-02 15:04:05"),
			method,
			statusCode,
			path,
			latency,
			clientIP,
			errorMessage,
		)

		if statusCode >= http.StatusInternalServerError {
			log.Println(logMessage)
		} else {
			log.Println(logMessage)
		}
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID := uint(claims["user_id"].(float64))
		c.Set("userID", userID)
		c.Next()
	}
}

// GenerateToken creates a new JWT token
func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve userID from the context (set by JWTAuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in token"})
			c.Abort()
			return
		}

		// Fetch the user from the database
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: User not found"})
			c.Abort()
			return
		}

		// Check if the user is an admin
		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Admin privileges required"})
			c.Abort()
			return
		}

		// Set isAdmin in the context for further use
		c.Set("isAdmin", true)
		c.Next()
	}
}
