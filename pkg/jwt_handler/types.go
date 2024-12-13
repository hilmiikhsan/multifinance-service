package jwt_handler

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserId   int    `json:"user_id"`
	Nik      string `json:"nik"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

type CostumClaimsPayload struct {
	UserId    int    `json:"user_id"`
	Nik       string `json:"nik"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	TokenType string `json:"token_type"`
}
