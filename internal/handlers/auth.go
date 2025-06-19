package handlers

import (
	"net/http"
	"strconv"

	"ishare-task-api/internal/auth"
	"ishare-task-api/internal/config"
	"ishare-task-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	oauth *auth.OAuthManager
	cfg   *config.Config
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(oauth *auth.OAuthManager, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		oauth: oauth,
		cfg:   cfg,
	}
}

// Authorize handles OAuth 2.0 authorization endpoint
// @Summary OAuth 2.0 Authorization
// @Description Initiates OAuth 2.0 authorization code flow
// @Tags OAuth
// @Accept x-www-form-urlencoded
// @Produce html
// @Param response_type formData string true "Must be 'code'" example(code)
// @Param client_id formData string true "OAuth client ID" example(test-client)
// @Param redirect_uri formData string true "Redirect URI" example(http://localhost:8080/oauth/callback)
// @Param scope formData string false "Requested scopes" example(tasks:read tasks:write)
// @Param state formData string false "State parameter for CSRF protection" example(random-state)
// @Success 200 {string} string "Authorization page"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /oauth/authorize [get]
func (h *AuthHandler) Authorize(c *gin.Context) {
	var req auth.AuthorizationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	// Validate response_type
	if req.ResponseType != "code" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "response_type must be 'code'",
		})
		return
	}

	// Validate client_id
	if req.ClientID != h.cfg.OAuth.ClientID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid client_id",
		})
		return
	}

	// Validate redirect_uri
	if req.RedirectURI != h.cfg.OAuth.RedirectURI {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid redirect_uri",
		})
		return
	}

	// For demo purposes, we'll show a simple login form
	// In a real application, you might redirect to a proper login page
	c.HTML(http.StatusOK, "authorize.html", gin.H{
		"client_id":    req.ClientID,
		"redirect_uri": req.RedirectURI,
		"scope":        req.Scope,
		"state":        req.State,
	})
}

// Login handles user login and creates authorization code
// @Summary User Login
// @Description Authenticates user and creates authorization code
// @Tags OAuth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email" example(user@example.com)
// @Param password formData string true "User password" example(password123)
// @Param client_id formData string true "OAuth client ID" example(test-client)
// @Param redirect_uri formData string true "Redirect URI" example(http://localhost:8080/oauth/callback)
// @Param scope formData string false "Requested scopes" example(tasks:read tasks:write)
// @Param state formData string false "State parameter" example(random-state)
// @Success 302 {string} string "Redirect to callback with authorization code"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /oauth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	clientID := c.PostForm("client_id")
	redirectURI := c.PostForm("redirect_uri")
	scope := c.PostForm("scope")
	state := c.PostForm("state")

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and password are required",
		})
		return
	}

	// Authenticate user
	user, err := h.oauth.AuthenticateUser(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Create authorization code
	authCode, err := h.oauth.CreateAuthorizationCode(user.ID, clientID, scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create authorization code",
		})
		return
	}

	// Redirect to callback with authorization code
	redirectURL := redirectURI + "?code=" + authCode.Code
	if state != "" {
		redirectURL += "&state=" + state
	}

	c.Redirect(http.StatusFound, redirectURL)
}

// Token handles OAuth 2.0 token endpoint
// @Summary OAuth 2.0 Token
// @Description Exchanges authorization code for access token
// @Tags OAuth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "Must be 'authorization_code'" example(authorization_code)
// @Param code formData string true "Authorization code" example(auth-code-here)
// @Param redirect_uri formData string true "Redirect URI" example(http://localhost:8080/oauth/callback)
// @Param client_id formData string true "OAuth client ID" example(test-client)
// @Param client_secret formData string true "OAuth client secret" example(test-secret)
// @Success 200 {object} auth.TokenResponse "Access token response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /oauth/token [post]
func (h *AuthHandler) Token(c *gin.Context) {
	var req auth.TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	// Validate grant_type
	if req.GrantType != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "grant_type must be 'authorization_code'",
		})
		return
	}

	// Validate authorization code
	authCode, err := h.oauth.ValidateAuthorizationCode(req.Code, req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create access token
	accessToken, err := h.oauth.CreateAccessToken(authCode.UserID, req.ClientID, authCode.Scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create access token",
		})
		return
	}

	// Return token response
	c.JSON(http.StatusOK, auth.TokenResponse{
		AccessToken: accessToken.Token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(24 * 60 * 60), // 24 hours in seconds
		Scope:       accessToken.Scope,
	})
}

// Callback handles OAuth callback
// @Summary OAuth Callback
// @Description Handles OAuth callback with authorization code
// @Tags OAuth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param code query string true "Authorization code" example(auth-code-here)
// @Param state query string false "State parameter" example(random-state)
// @Success 200 {object} map[string]interface{} "Callback response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /oauth/callback [get]
func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization code is required",
		})
		return
	}

	// For demo purposes, return the authorization code
	// In a real application, you might redirect to a frontend application
	c.JSON(http.StatusOK, gin.H{
		"message": "Authorization successful",
		"code":    code,
		"state":   state,
		"next_step": "Exchange this code for an access token using POST /oauth/token",
	})
}

// Register handles user registration
// @Summary User Registration
// @Description Creates a new user account
// @Tags OAuth
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User registration data"
// @Success 201 {object} models.UserResponse "User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Router /oauth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Create user
	user, err := h.oauth.CreateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Return user response (without password)
	c.JSON(http.StatusCreated, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// CleanupTokens handles cleanup of expired tokens
// @Summary Cleanup Expired Tokens
// @Description Removes expired authorization codes and access tokens
// @Tags OAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Cleanup completed"
// @Router /oauth/cleanup [post]
func (h *AuthHandler) CleanupTokens(c *gin.Context) {
	if err := h.oauth.CleanupExpiredTokens(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to cleanup expired tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expired tokens cleaned up successfully",
	})
} 