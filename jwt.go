// sg-common/jwt/jwt.go
package jwt

import (
	"errors"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

// A custom type for the context key to avoid potential collisions with other context keys.
type contextKey string

const (
	// ClaimsKey is the context key for storing JWT claims in the request context.
	ClaimsKey contextKey = "claims"
	// TokenKey is the context key for storing the JWT token in the request context.
	TokenKey contextKey = "token"
)

type JWTConfig struct {
	VerifyIssuer    bool
	Issuer          string
	SigningMethod   golangjwt.SigningMethod
	SigningKey      []byte
	ValidationDelay time.Duration
	Leeway          int64
}

// Add a custom Claims structure to include internal role, external roles, and permissions
type Claims struct {
	golangjwt.RegisteredClaims
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	Type        string   `json:"type"`
}

// IsValidJWT is a custom validator to check if the provided string is a valid JWT token.
func IsValidJWT(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid JWT token")
	}

	parser := golangjwt.Parser{}
	_, _, err := parser.ParseUnverified(str, &Claims{})
	if err != nil {
		return errors.New("invalid JWT token")
	}

	return nil
}

// GetSubject extracts the subject (user ID) from a JWT token.
func GetSubject(tokenString string) (string, error) {
	parser := golangjwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New("invalid JWT claims")
	}

	return claims.Subject, nil
}

// Parse function validates and extracts the JWT claims.
func Parse(tokenString string, signingKey []byte) (*Claims, error) {
	token, err := golangjwt.ParseWithClaims(tokenString, &Claims{}, func(token *golangjwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid JWT claims")
	}

	return claims, nil
}


// ExtractToken removes "Bearer " prefix from Authorization header
func ExtractToken(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:] // Extract the token part
	}
	return authHeader // Return as is if it's not in "Bearer <token>" format
}
