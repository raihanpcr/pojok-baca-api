package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"log"
	"strings"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			// Pastikan Bearer token valid
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Println("Invalid Authorization Header Format")
				return echo.ErrUnauthorized
			}
			tokenString := strings.TrimSpace(parts[1])

			// Parse token
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				// Validasi metode signing
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				log.Printf("JWT Parsing Error : %v\n", err)
				return echo.ErrUnauthorized
			}

			c.Set("user", token)
			return next(c)
		}
	}
}
