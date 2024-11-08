package session

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
)

var (
	NoSessionInContext = errors.New("no session in context")
	NoSessionInCookie  = errors.New("no session in cookie")
	WrongJWTMethod     = errors.New("JWT token method isn't HS256")
	NoAuthError        = errors.New("wrong token or impossible to parse field")
)

func (sm *SessionManager) NewSession(username string) *abstractions.Session {
	return &abstractions.Session{
		Username: username,
	}
}

func (sm *SessionManager) SessionFromContext(ctx echo.Context) (*abstractions.Session, error) {
	currentSession, ok := ctx.Get(sm.sessionKey).(*abstractions.Session)
	if !ok {
		return nil, NoSessionInContext
	}
	return currentSession, nil
}

func (sm *SessionManager) SessionWithContext(ctx echo.Context, session *abstractions.Session) {
	ctx.Set(sm.sessionKey, session)
}
