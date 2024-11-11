package postgresDB

import "github.com/nehachuha1/mynotes-project/pkg/abstractions"

type IPostgresRepo interface {
	CreateUser(newUser *abstractions.User) error
	RegisterUser(newRegistration *abstractions.Registration) error
	AuthorizeUser(user *abstractions.Registration) (*abstractions.Registration, error)
	DeleteUser(user *abstractions.User) error
}
