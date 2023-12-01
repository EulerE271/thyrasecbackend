package helpers

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go" // Import "jwt-go" here
)

// Rest of the code in auth.go

func ExtractTokenFromCookie(r *http.Request) (string, error) {
	// Get the token from the "token" cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	// The entire cookie value is treated as the token
	tokenString := cookie.Value

	// Return the token value
	return tokenString, nil
}

func ValidateToken(tokenString string, secretKey string) (*jwt.Token, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return token, nil
}
