package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ishare-task-api/internal/config"
	"ishare-task-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager handles JWT token operations
type JWTManager struct {
	config config.JWTConfig
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg config.JWTConfig) *JWTManager {
	return &JWTManager{
		config: cfg,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Scope  string    `json:"scope"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user
func (j *JWTManager) GenerateToken(user *models.User, scope string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Scope:  scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Audience:  []string{j.config.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.Expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}

// ValidateToken validates and parses a JWT token
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateJWS generates a JWS token (JWT with explicit JWS structure)
func (j *JWTManager) GenerateJWS(user *models.User, scope string) (string, error) {
	now := time.Now()
	
	// Create JWS header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	
	// Create JWS payload
	payload := map[string]interface{}{
		"sub": user.ID.String(),
		"email": user.Email,
		"scope": scope,
		"iss": j.config.Issuer,
		"aud": j.config.Audience,
		"exp": now.Add(j.config.Expiration).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}

	// Encode header and payload
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create signature
	signingInput := headerB64 + "." + payloadB64
	signature := j.sign(signingInput)

	// Combine to form JWS
	jws := signingInput + "." + signature
	return jws, nil
}

// ValidateJWS validates a JWS token
func (j *JWTManager) ValidateJWS(jws string) (*Claims, error) {
	parts := strings.Split(jws, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWS format")
	}

	headerB64, payloadB64, signatureB64 := parts[0], parts[1], parts[2]

	// Verify signature
	signingInput := headerB64 + "." + payloadB64
	expectedSignature := j.sign(signingInput)
	
	if signatureB64 != expectedSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	// Check expiration
	if exp, ok := payload["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, fmt.Errorf("token expired")
		}
	}

	// Extract claims
	userIDStr, ok := payload["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	email, _ := payload["email"].(string)
	scope, _ := payload["scope"].(string)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Scope:  scope,
	}

	return claims, nil
}

// sign creates HMAC-SHA256 signature
func (j *JWTManager) sign(input string) string {
	h := hmac.New(sha256.New, []byte(j.config.Secret))
	h.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// HasScope checks if the token has the required scope
func (j *JWTManager) HasScope(claims *Claims, requiredScope string) bool {
	if claims.Scope == "" {
		return false
	}
	
	scopes := strings.Split(claims.Scope, " ")
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	
	return false
} 