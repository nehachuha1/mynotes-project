package redisDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"github.com/nehachuha1/mynotes-project/pkg/services/config"
	"math/rand"
)

var (
	CannotCreateSession          = errors.New("cannot create session")
	ResultIsNotOK                = errors.New("redis.DO result is not okay")
	CantGetSessionWithKey        = errors.New("get session with key")
	CantUnmarshalToSessionStruct = errors.New("can't unmarshal to session struct")
	CantRemoveSession            = errors.New("can't delete session")
)

type RedisDatabase struct {
	RedisConnection redis.Pool
}

func NewRedisDatabase(cfg *config.Config) *RedisDatabase {
	return &RedisDatabase{
		RedisConnection: redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(cfg.RedisConfig.RedisURL)
			},
			MaxIdle:     8,
			MaxActive:   0,
			IdleTimeout: 100,
		},
	}
}

var _ IRedisControl = &RedisDatabase{}

func (rdb *RedisDatabase) CreateSession(session *abstractions.Session) (*abstractions.Session, error) {
	session.SessionID = RandStringRunes(32)
	dataSerialized, _ := json.Marshal(session)
	mKey := "SESSION: " + session.SessionID
	currentConnection := rdb.RedisConnection.Get()
	result, err := redis.String(currentConnection.Do("SET", mKey, dataSerialized, "EX", 259200))

	if err != nil {
		return nil, fmt.Errorf("cannot create session: %v", err)
	}
	if result != "OK" {
		return nil, fmt.Errorf("redis.DO result is not okay: %v", err)
	}
	return session, nil
}

func (rdb *RedisDatabase) CheckSession(session *abstractions.Session) (*abstractions.Session, error) {
	mKey := "SESSION: " + session.SessionID
	currentConnection := rdb.RedisConnection.Get()
	data, err := redis.Bytes(currentConnection.Do("GET", mKey))
	if err != nil {
		return nil, fmt.Errorf("can't get session with key: %v", err)
	}
	currentSession := &abstractions.Session{}
	err = json.Unmarshal(data, currentSession)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal to session struct: %v", err)
	}
	return currentSession, nil
}

func (rdb *RedisDatabase) DeleteSession(session *abstractions.Session) error {
	mKey := "SESSION: " + session.SessionID
	currentConnection := rdb.RedisConnection.Get()
	_, err := redis.Int(currentConnection.Do("DEL", mKey))
	if err != nil {
		return fmt.Errorf("can't delete session: %v", err)
	}
	return nil
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
