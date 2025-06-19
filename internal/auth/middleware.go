package auth

import (
	"net/http"
	"strings"

	"ishare-task-api/internal/config"
	"ishare-task-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthMiddleware provides authentication middleware
type AuthMiddleware struct {
	jwt *JWTManager
	db  *gorm.DB
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwt *JWTManager, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
		db:  db,
	}
}

// Authenticate middleware validates JWS tokens and sets user context
func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Use 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate JWS token
		claims, err := a.jwt.ValidateJWS(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Verify token exists in database
		var accessToken models.AccessToken
		if err := a.db.Where("token = ? AND expires_at > NOW()", tokenString).First(&accessToken).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token not found or expired",
			})
			c.Abort()
			return
		}

		// Get user from database
		var user models.User
		if err := a.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		// Set user and claims in context
		c.Set("user", &user)
		c.Set("claims", claims)
		c.Set("access_token", &accessToken)

		c.Next()
	}
}

// RequireScope middleware checks if the user has the required scope
func (a *AuthMiddleware) RequireScope(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsInterface, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid claims format",
			})
			c.Abort()
			return
		}

		if !a.jwt.HasScope(claims, requiredScope) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext gets the user from the Gin context
func GetUserFromContext(c *gin.Context) (*models.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	user, ok := userInterface.(*models.User)
	return user, ok
}

// GetClaimsFromContext gets the claims from the Gin context
func GetClaimsFromContext(c *gin.Context) (*Claims, bool) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	claims, ok := claimsInterface.(*Claims)
	return claims, ok
}

// GetAccessTokenFromContext gets the access token from the Gin context
func GetAccessTokenFromContext(c *gin.Context) (*models.AccessToken, bool) {
	tokenInterface, exists := c.Get("access_token")
	if !exists {
		return nil, false
	}

	token, ok := tokenInterface.(*models.AccessToken)
	return token, ok
}

// ValidateUserOwnership middleware ensures the user owns the resource
func (a *AuthMiddleware) ValidateUserOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetUserFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Get resource ID from URL parameter
		resourceID := c.Param("id")
		if resourceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Resource ID required",
			})
			c.Abort()
			return
		}

		// Parse UUID
		resourceUUID, err := uuid.Parse(resourceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid resource ID format",
			})
			c.Abort()
			return
		}

		// Check if user owns the resource (for tasks)
		var task models.Task
		if err := a.db.Where("id = ?", resourceUUID).First(&task).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Resource not found",
			})
			c.Abort()
			return
		}

		// For now, we'll allow all authenticated users to access all tasks
		// In a real application, you might want to add user_id to tasks table
		// and check ownership here

		c.Next()
	}
} 