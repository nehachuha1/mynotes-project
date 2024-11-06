package session

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"time"
)

type SessionManager struct {
	// TODO: добавить привязку к базе данных
}

func CreateNewToken(user abstractions.User, sessionID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"SessionID": sessionID,
			"Username":  user.Username,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
	})
	newToken, err := token.SignedString(jwtSecretKey)
	return "Bearer: " + newToken, err
}

func (sm *SessionManager) CheckSession(c echo.Context) (*abstractions.Session, error) {
	tokenWithCookie, err := c.Cookie("SESSION")
	if err != nil {
		return nil, NoSessionInCookie
	}
	tokenString := tokenWithCookie.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, WrongJWTMethod
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, NoAuthError
	}
	user, isOk := (claims["user"]).(map[string]interface{})
	if !isOk {
		return nil, NoAuthError
	}
	sessionID, isOk := (user["SessionID"]).(string)
	if !isOk {
		return nil, NoAuthError
	}
	//	TODO: поход в Redis чтобы проверить актуален токен или нет
	//  return currentSession, nil
}
