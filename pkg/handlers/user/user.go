package user

import "github.com/nehachuha1/mynotes-project/pkg/abstractions"

type IUserRepo interface {
	CreateUser(newUser *abstractions.User) error
	AuthorizeUser(user *abstractions.User) *abstractions.Session
	DeleteUser(user *abstractions.User)
}
