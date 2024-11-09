package user

import "github.com/nehachuha1/mynotes-project/pkg/services/session"

type UserRepository struct {
	SessionControl *session.SessionManager
	//PostgresDB
}
