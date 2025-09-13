package main

import (
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// Test JWT token from the curl response
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiIiLCJyb2xlcyI6WyJhZG1pbiJdLCJwZXJtaXNzaW9ucyI6W10sInRva2VuX3R5cGUiOiJhY2Nlc3MiLCJpc3MiOiJhdXRoLXNlcnZpY2UiLCJzdWIiOiJhZG1pbiIsImV4cCI6MTc1NDYyNTUwMiwibmJmIjoxNzU0NTM5MTAyLCJpYXQiOjE3NTQ1MzkxMDJ9.jctiOJ6efJKMh39VOx3Zsl9z-gzYPXcOdOEJ_WDUOjA"
	jwtSecret := "your-secret-key" // This should match the actual secret

	// Parse with ChatJWTClaims
	token, err := jwt.ParseWithClaims(tokenString, &interfaces.ChatJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return
	}

	claims, ok := token.Claims.(*interfaces.ChatJWTClaims)
	if !ok || !token.Valid {
		log.Printf("Invalid token claims")
		return
	}

	fmt.Printf("UserID: '%s'\n", claims.UserID)
	fmt.Printf("Username: '%s'\n", claims.Username)
	fmt.Printf("Roles: %v\n", claims.Roles)
	fmt.Printf("TokenType: '%s'\n", claims.TokenType)
}
