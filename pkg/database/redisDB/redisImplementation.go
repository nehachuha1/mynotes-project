package redisDB

import "github.com/nehachuha1/mynotes-project/pkg/abstractions"

type IRedisControl interface {
	CreateSession(session *abstractions.Session) (*abstractions.Session, error)
	DeleteSession(session *abstractions.Session) error
	CheckSession(session *abstractions.Session) (*abstractions.Session, error)
}
