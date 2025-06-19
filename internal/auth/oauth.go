package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"ishare-task-api/internal/config"
	"ishare-task-api/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// OAuthManager handles OAuth 2.0 operations
type OAuthManager struct {
	config config.OAuthConfig
	db     *gorm.DB
	jwt    *JWTManager
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager(cfg config.OAuthConfig, db *gorm.DB, jwt *JWTManager) *OAuthManager {
	return &OAuthManager{
		config: cfg,
		db:     db,
		jwt:    jwt,
	}
}

// AuthorizationRequest represents an OAuth authorization request
type AuthorizationRequest struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
}

// TokenRequest represents an OAuth token request
type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
}

// TokenResponse represents an OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

// CreateAuthorizationCode creates a new authorization code for OAuth flow
func (o *OAuthManager) CreateAuthorizationCode(userID uuid.UUID, clientID, scope string) (*models.AuthorizationCode, error) {
	// Generate random authorization code
	codeBytes := make([]byte, 32)
	if _, err := rand.Read(codeBytes); err != nil {
		return nil, err
	}
	code := base64.URLEncoding.EncodeToString(codeBytes)

	// Create authorization code record
	authCode := &models.AuthorizationCode{
		Code:      code,
		UserID:    userID,
		ClientID:  clientID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(10 * time.Minute), // Authorization codes expire in 10 minutes
	}

	if err := o.db.Create(authCode).Error; err != nil {
		return nil, err
	}

	return authCode, nil
}

// ValidateAuthorizationCode validates and consumes an authorization code
func (o *OAuthManager) ValidateAuthorizationCode(code, clientID, clientSecret string) (*models.AuthorizationCode, error) {
	var authCode models.AuthorizationCode
	
	if err := o.db.Where("code = ? AND client_id = ? AND expires_at > ?", 
		code, clientID, time.Now()).First(&authCode).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired authorization code")
	}

	// Validate client secret
	if clientSecret != o.config.ClientSecret {
		return nil, fmt.Errorf("invalid client secret")
	}

	// Delete the used authorization code
	o.db.Delete(&authCode)

	return &authCode, nil
}

// CreateAccessToken creates a new access token
func (o *OAuthManager) CreateAccessToken(userID uuid.UUID, clientID, scope string) (*models.AccessToken, error) {
	// Generate JWS token
	user := &models.User{ID: userID}
	tokenString, err := o.jwt.GenerateJWS(user, scope)
	if err != nil {
		return nil, err
	}

	// Create access token record
	accessToken := &models.AccessToken{
		Token:     tokenString,
		UserID:    userID,
		ClientID:  clientID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Access tokens expire in 24 hours
	}

	if err := o.db.Create(accessToken).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

// ValidateAccessToken validates an access token
func (o *OAuthManager) ValidateAccessToken(tokenString string) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	
	if err := o.db.Where("token = ? AND expires_at > ?", 
		tokenString, time.Now()).First(&accessToken).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired access token")
	}

	return &accessToken, nil
}

// AuthenticateUser authenticates a user with email and password
func (o *OAuthManager) AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User
	
	if err := o.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

// CreateUser creates a new user with hashed password
func (o *OAuthManager) CreateUser(email, password string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := o.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// CleanupExpiredTokens removes expired tokens from the database
func (o *OAuthManager) CleanupExpiredTokens() error {
	// Delete expired authorization codes
	if err := o.db.Where("expires_at < ?", time.Now()).Delete(&models.AuthorizationCode{}).Error; err != nil {
		return err
	}

	// Delete expired access tokens
	if err := o.db.Where("expires_at < ?", time.Now()).Delete(&models.AccessToken{}).Error; err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (o *OAuthManager) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	
	if err := o.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
} 