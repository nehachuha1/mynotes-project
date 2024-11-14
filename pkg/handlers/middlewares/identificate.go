package middlewares

import (
	"github.com/labstack/echo/v4"
	"math/rand"
)

func IdentifyRequest() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := generateRequestID()
			c.Set("RequestID", requestID)
			_ = next(c)
			return nil
		}
	}
}

func generateRequestID() string {
	b := make([]rune, 32)
	letterRunes := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
