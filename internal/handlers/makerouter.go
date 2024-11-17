package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nehachuha1/mynotes-project/internal/config"
	"github.com/nehachuha1/mynotes-project/internal/handlers/middlewares"
	"github.com/nehachuha1/mynotes-project/internal/services/session"
)

func GenerateRoutesWithMiddlewares(cfg *config.Config, sm *session.SessionManager) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middlewares.IdentifyRequest())
	e.Use(middlewares.Auth(sm))
	return e
}
