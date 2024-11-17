package session

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nehachuha1/mynotes-project/internal/abstractions"
	"github.com/nehachuha1/mynotes-project/internal/config"
	"github.com/nehachuha1/mynotes-project/pkg/database/redisDB"
	"go.uber.org/zap"
	"strings"
	"time"
)

type SessionManager struct {
	sessionKey   string
	jwtSecretKey []byte
	RedisDB      *redisDB.RedisDatabase
}

func NewSessionManager(cfg *config.Config, logger *zap.SugaredLogger) *SessionManager {
	return &SessionManager{
		sessionKey:   cfg.SessionConfig.SessionKey,
		jwtSecretKey: []byte(cfg.SessionConfig.SessionKey),
		RedisDB:      redisDB.NewRedisDatabase(cfg, logger),
	}
}

func (sm *SessionManager) CreateNewToken(user *abstractions.User, sessionID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"SessionID": sessionID,
			"Username":  user.Username,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
	})
	newToken, err := token.SignedString(sm.jwtSecretKey)
	return "Bearer: " + newToken, err
}

func (sm *SessionManager) CheckSession(c echo.Context) (*abstractions.Session, error) {
	tokenWithCookie, err := c.Cookie("SESSION")
	if err != nil {
		return nil, NoSessionInCookie
	}
	_, tokenString, ok := strings.Cut(tokenWithCookie.Value, "Bearer ")
	if !ok {
		return nil, fmt.Errorf("there's no prefix 'Bearer: ' in current token")
	}
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, WrongJWTMethod
		}
		return sm.jwtSecretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("wrong signing method")
	}
	user, isOk := (claims["user"]).(map[string]interface{})
	if !isOk {
		return nil, fmt.Errorf("impossible to parse map-field 'User'")
	}
	sessionID, isOk := (user["SessionID"]).(string)
	if !isOk {
		return nil, fmt.Errorf("impossible to parse field 'SessionID'")
	}
	username, isOk := (user["Username"]).(string)
	if !isOk {
		return nil, fmt.Errorf("impossible to parse field 'Username'")
	}
	userSession := &abstractions.Session{SessionID: sessionID, Username: username}
	validSession, err := sm.RedisDB.CheckSession(userSession)
	if err != nil {
		return nil, fmt.Errorf("current session is not valid: %v", err)
	}
	return validSession, nil
}

func (sm *SessionManager) CreateSession(username string) (*abstractions.Session, error) {
	newSession := sm.newSession(username)
	sessionWithID, err := sm.RedisDB.CreateSession(newSession)
	if err != nil {
		return nil, fmt.Errorf("unable to create session: %v", err)
	}
	checkedSession, err := sm.RedisDB.CheckSession(sessionWithID)
	if err != nil {
		return nil, fmt.Errorf("can't check created session: %v", err)
	}
	return checkedSession, nil
}

func (sm *SessionManager) DeleteSession(c echo.Context) error {
	sessionFromContext, err := sm.SessionFromContext(c)
	if err != nil {
		return fmt.Errorf("can't get session from context: %v", err)
	}
	currentSession, err := sm.RedisDB.CheckSession(sessionFromContext)
	if err != nil {
		return fmt.Errorf("can't check session from context: %v", err)
	}
	err = sm.RedisDB.DeleteSession(currentSession)
	if err != nil {
		return fmt.Errorf("can't delete session with sessionID %v: %v", currentSession.SessionID, err)
	}
	return nil
}
