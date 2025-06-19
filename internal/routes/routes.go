package routes

import (
	"ishare-task-api/internal/auth"
	"ishare-task-api/internal/config"
	"ishare-task-api/internal/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Setup configures all routes and middleware
func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Initialize auth components
	jwtManager := auth.NewJWTManager(cfg.JWT)
	oauthManager := auth.NewOAuthManager(cfg.OAuth, db, jwtManager)
	authMiddleware := auth.NewAuthMiddleware(jwtManager, db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(oauthManager, cfg)
	taskHandler := handlers.NewTaskHandler(db)

	// Load HTML templates for OAuth flow
	router.LoadHTMLGlob("templates/*")

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "iSHARE Task API is running",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// OAuth 2.0 routes (no authentication required)
	oauth := router.Group("/oauth")
	{
		oauth.GET("/authorize", authHandler.Authorize)
		oauth.POST("/login", authHandler.Login)
		oauth.POST("/token", authHandler.Token)
		oauth.GET("/callback", authHandler.Callback)
		oauth.POST("/register", authHandler.Register)
		oauth.POST("/cleanup", authHandler.CleanupTokens)
	}

	// Task management routes (authentication required)
	tasks := router.Group("/tasks")
	tasks.Use(authMiddleware.Authenticate())
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.GET("", taskHandler.ListTasks)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}

	// API documentation endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": "iSHARE Task Management API",
			"version": "1.0.0",
			"description": "A secure REST API for task management with OAuth 2.0 and JWS token signing",
			"documentation": "/swagger/index.html",
			"endpoints": gin.H{
				"oauth": gin.H{
					"authorize": "GET /oauth/authorize - OAuth 2.0 authorization endpoint",
					"token": "POST /oauth/token - OAuth 2.0 token endpoint",
					"callback": "GET /oauth/callback - OAuth callback endpoint",
					"register": "POST /oauth/register - User registration",
				},
				"tasks": gin.H{
					"create": "POST /tasks - Create a new task",
					"list": "GET /tasks - List all tasks",
					"get": "GET /tasks/{id} - Get a specific task",
					"update": "PUT /tasks/{id} - Update a task",
					"delete": "DELETE /tasks/{id} - Delete a task",
				},
			},
			"authentication": "All task endpoints require Bearer token authentication",
		})
	})

	return router
} 