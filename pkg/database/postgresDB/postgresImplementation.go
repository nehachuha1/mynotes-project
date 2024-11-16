package postgresDB

import "github.com/nehachuha1/mynotes-project/pkg/abstractions"

type IPostgresRepo interface {
	CreateUser(newUser *abstractions.User) error
	RegisterUser(newRegistration *abstractions.Registration) error
	AuthorizeUser(user *abstractions.Registration) (*abstractions.Registration, error)
	DeleteUser(user *abstractions.User) error

	GetUserWorkspaces(user *abstractions.User) ([]*abstractions.Workspace, error)
	CreateWorkspace(newWorkspace *abstractions.Workspace) (*abstractions.Workspace, error)
	DeleteWorkspace(workspace *abstractions.Workspace) error
	EditWorkspacePrivacy(workspace *abstractions.Workspace) error

	GetUserNotes(user *abstractions.User) ([]*abstractions.Note, error)
	GetUserNote(noteID *abstractions.NoteID) (*abstractions.Note, error)
	CreateNote(note *abstractions.Note) error
	DeleteNote(note *abstractions.Note) error
	EditNote(note *abstractions.Note) error
}
