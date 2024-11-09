package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/nehachuha1/mynotes-project/pkg/services/session"
)

func Auth(sm *session.SessionManager) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			currentSession, err := sm.CheckSession(c)
			if err == nil || currentSession != nil {
				sm.SessionWithContext(c, currentSession)
				_ = next(c)
			} else {
				_ = next(c)
			}
			return nil
		}
	}
}
