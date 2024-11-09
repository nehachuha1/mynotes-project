package postgresDB

import "github.com/nehachuha1/mynotes-project/pkg/abstractions"

type IPostgresRepo interface {
	CreateUser(newUser *abstractions.User) error
	AuthorizeUser(user *abstractions.User) *abstractions.Session
	DeleteUser(user *abstractions.User)
}
