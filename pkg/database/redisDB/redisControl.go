package redisDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/nehachuha1/mynotes-project/internal/abstractions"
	"github.com/nehachuha1/mynotes-project/internal/config"
	"go.uber.org/zap"
	"math/rand"
	"time"
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
	logger          *zap.SugaredLogger
}

func NewRedisDatabase(cfg *config.Config, logger *zap.SugaredLogger) *RedisDatabase {
	return &RedisDatabase{
		RedisConnection: redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(cfg.RedisConfig.RedisURL)
			},
			MaxIdle:     8,
			MaxActive:   0,
			IdleTimeout: 100,
		},
		logger: logger,
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
		rdb.logger.Warnw("can't create session in redis database", "type", "redis", "output", err.Error(), "time", time.Now().String())
		return nil, fmt.Errorf("cannot create session: %v", err)
	}
	if result != "OK" {
		rdb.logger.Warnw("redis.DO result is not okay", "type", "redis", "output", result, "time", time.Now().String())
		return nil, fmt.Errorf("redis.DO result is not okay: %v", err)
	}
	rdb.logger.Infow("successfully created session", "type", "redis",
		"output", "CREATED SESSION IN REDIS", "time", time.Now().String())
	return session, nil
}

func (rdb *RedisDatabase) CheckSession(session *abstractions.Session) (*abstractions.Session, error) {
	mKey := "SESSION: " + session.SessionID
	currentConnection := rdb.RedisConnection.Get()
	data, err := redis.Bytes(currentConnection.Do("GET", mKey))
	if err != nil {
		rdb.logger.Warnw("can't get session with key", "type", "redis", "output", err.Error(), "time", time.Now().String())
		return nil, fmt.Errorf("can't get session with key: %v", err)
	}
	currentSession := &abstractions.Session{}
	err = json.Unmarshal(data, currentSession)
	if err != nil {
		rdb.logger.Warnw("can't unmarshal to session struct", "type", "redis", "output", err.Error(), "time", time.Now().String())
		return nil, fmt.Errorf("can't unmarshal to session struct: %v", err)
	}
	rdb.logger.Infow("successfully checked session", "type", "redis",
		"output", "CHECKED SESSION IN REDIS", "time", time.Now().String())
	return currentSession, nil
}

func (rdb *RedisDatabase) DeleteSession(session *abstractions.Session) error {
	mKey := "SESSION: " + session.SessionID
	currentConnection := rdb.RedisConnection.Get()
	_, err := redis.Int(currentConnection.Do("DEL", mKey))
	if err != nil {
		rdb.logger.Warnw("can't delete session", "type", "redis", "output", err.Error(), "time", time.Now().String())
		return fmt.Errorf("can't delete session: %v", err)
	}
	rdb.logger.Infow("successfully deleted session", "type", "redis",
		"output", "DELETE SESSION IN REDIS", "time", time.Now().String())
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
