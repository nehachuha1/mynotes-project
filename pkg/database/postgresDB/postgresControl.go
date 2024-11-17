package postgresDB

import (
	"errors"
	"fmt"
	"github.com/nehachuha1/mynotes-project/internal/abstractions"
	"github.com/nehachuha1/mynotes-project/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type PostgresDatabase struct {
	database *gorm.DB
	logger   *zap.SugaredLogger
}

func makeDsn(cfg *config.Config) string {
	dsn := "postgres://" + cfg.PostgresConfig.PostgresUser + ":" + cfg.PostgresConfig.PostgresPassword + "@" +
		cfg.PostgresConfig.PostgresAddress + ":" + cfg.PostgresConfig.PostgresPort + "/" +
		cfg.PostgresConfig.PostgresDB
	return dsn
}

func NewPostgresDB(cfg *config.Config, logger *zap.SugaredLogger) *PostgresDatabase {
	newDsn := makeDsn(cfg)
	dbConn, err := gorm.Open(postgres.Open(newDsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("can't open gorm.Open: %v", err))
	}
	return &PostgresDatabase{
		database: dbConn,
		logger:   logger,
	}
}

func (pgdb *PostgresDatabase) RegisterUser(newRegistration *abstractions.Registration) error {
	result := &abstractions.Registration{}
	pgdb.database.Table("relation_registrations").Where("username = ?", newRegistration.Username).First(result)
	if result.GetUsername() != "" {
		pgdb.logger.Warnw("failed register user in postgres control", "type", "postgres",
			"output", "user with this username is already registered", "time", time.Now().String())
		return fmt.Errorf("user with username %v already exists", newRegistration.GetUsername())
	}
	newUser := &RelationRegistration{
		Username: newRegistration.GetUsername(),
		Password: newRegistration.GetPassword(),
	}
	resultFromDB := pgdb.database.Create(newUser)
	if resultFromDB.Error != nil {
		pgdb.logger.Warnw("failed register user in postgres control", "type", "postgres",
			"output", resultFromDB.Error, "time", time.Now().String())
		return fmt.Errorf("failed register user in postgres control: %v", resultFromDB.Error)
	}
	pgdb.logger.Infow("successfully registered user", "type", "postgres",
		"output", "REGISTERED USER IN POSTGRES", "time", time.Now().String())
	return nil
}

func (pgdb *PostgresDatabase) CreateUser(user *abstractions.User) error {
	foundRegistration := &RelationRegistration{}
	result := pgdb.database.Table("relation_registrations").Where("username = ?", user.GetUsername()).First(
		foundRegistration)
	if !errors.Is(result.Error, nil) {
		return fmt.Errorf("can't find user in registrations")
	}
	currentUser := &RelationUser{
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Initials: user.GetInitials(),
		Telegram: user.GetTelegram(),
	}
	foundUser := &RelationUser{}
	result = pgdb.database.Table("relation_users").Where("username = ?", user.GetUsername()).First(
		foundUser)
	if errors.Is(result.Error, nil) {
		pgdb.logger.Warnw("failed creating user in postgres control", "type", "postgres",
			"output", fmt.Errorf("user already exists in 'relation_users' database"), "time", time.Now().String())
		return fmt.Errorf("user already exists in 'relation_users' database")
	}
	result = pgdb.database.Create(currentUser)
	if result.Error != nil {
		pgdb.logger.Warnw("failed creating user in postgres control", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("failed creating user in postgres control: %v", result.Error)
	}
	pgdb.logger.Infow("successfully created user", "type", "postgres",
		"output", "CREATED USER IN POSTGRES", "time", time.Now().String())
	return nil
}

func (pgdb *PostgresDatabase) AuthorizeUser(user *abstractions.Registration) (*abstractions.Registration, error) {
	userToAuthorize := &RelationRegistration{
		Username: user.GetUsername(),
	}
	result := pgdb.database.Table("relation_registrations").Where("username = ?", userToAuthorize.Username).First(userToAuthorize)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		pgdb.logger.Warnw("failed on authorizing user", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return nil, fmt.Errorf("failed on authorizing user: %v", result.Error)
	}

	checkedUser := &abstractions.Registration{
		Id:       userToAuthorize.Id,
		Username: userToAuthorize.Username,
		Password: userToAuthorize.Password,
	}
	pgdb.logger.Infow("successfully got user", "type", "postgres",
		"output", "GREP USER FROM POSTGRES", "time", time.Now().String())
	return checkedUser, nil
}

func (pgdb *PostgresDatabase) DeleteUser(user *abstractions.User) error {
	userToDelete := &RelationUser{
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Initials: user.GetInitials(),
		Telegram: user.GetTelegram(),
	}
	checkedUser := &RelationUser{}
	result := pgdb.database.Table("relation_users").Where("username = ?", userToDelete.Username).First(checkedUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		pgdb.logger.Warnw("can't find user in database that could be deleted", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't find user in database that could be deleted: %v", result.Error)
	}
	result = pgdb.database.Table("relation_users").Delete(&RelationUser{}, checkedUser.Id)
	if result.Error != nil {
		pgdb.logger.Warnw("can't delete user in relation_user table", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't delete user in relation_user table: %v", result.Error)
	}
	registeredUser := &RelationRegistration{}
	recordInRegistrations := pgdb.database.Table(
		"relation_registrations").Where("username = ?", checkedUser.Username).First(registeredUser)
	if recordInRegistrations.Error != nil {
		return fmt.Errorf(
			"user that should be deleted from registrations and users cannot be found in table 'registrations'")
	}
	result = pgdb.database.Table("relation_registrations").Delete(&RelationRegistration{}, registeredUser.Id)
	if result.Error != nil {
		pgdb.logger.Warnw("can't delete user in relation_registration table", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't delete user in relation_registration table: %v", result.Error)
	}
	pgdb.logger.Infow("successfully deleted user", "type", "postgres",
		"output", "DELETED USER FROM POSTGRES", "time", time.Now().String())
	return nil
}

func (pgdb *PostgresDatabase) GetUserWorkspaces(user *abstractions.User) ([]*abstractions.Workspace, error) {
	var allWorkspaces []RelationWorkspace
	result := pgdb.database.Table("relation_workspaces").Where(
		"owner_username = ?", user.GetUsername()).Find(&allWorkspaces)
	if result.Error != nil {
		pgdb.logger.Warnw("can't find user workspaces in 'relation_workspaces'", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return []*abstractions.Workspace{}, fmt.Errorf("can't find user workspaces in 'relation_workspaces': %v", result.Error)
	}
	toReturnWorkspaces := make([]*abstractions.Workspace, len(allWorkspaces))
	for ind, ws := range allWorkspaces {
		newWorkspace := &abstractions.Workspace{
			Id:            ws.Id,
			OwnerUsername: ws.OwnerUsername,
			IsPrivate:     ws.IsPrivate,
			NotesID:       ws.NotesID,
		}
		toReturnWorkspaces[ind] = newWorkspace
	}
	return toReturnWorkspaces, nil
}

func (pgdb *PostgresDatabase) CreateWorkspace(newWorkspace *abstractions.Workspace) (*abstractions.Workspace, error) {
	validationErr := validateNewWorkspace(newWorkspace)
	if validationErr != nil {
		pgdb.logger.Warnw("no input owner username of workspace", "type", "postgres",
			"output", validationErr.Error(), "time", time.Now().String())
		return &abstractions.Workspace{}, validationErr
	}
	toCreateWorkspace := &RelationWorkspace{
		OwnerUsername: newWorkspace.GetOwnerUsername(),
		IsPrivate:     true,
		NotesID:       make([]int64, 0),
	}
	result := pgdb.database.Table("relation_workspaces").Create(toCreateWorkspace)
	if result.Error != nil {
		pgdb.logger.Warnw("error by creating new workspace", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return &abstractions.Workspace{}, result.Error
	}
	currentWorkspace := &RelationWorkspace{}
	result = pgdb.database.Table("relation_workspaces").Where("owner_username = ?",
		newWorkspace.GetOwnerUsername()).Order("id desc").First(currentWorkspace)
	if result.Error != nil {
		pgdb.logger.Warnw("can't find created workspace", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return &abstractions.Workspace{}, fmt.Errorf("can't find created workspace: %v", result.Error)
	}
	createdWorkspace := &abstractions.Workspace{
		Id:            currentWorkspace.Id,
		OwnerUsername: currentWorkspace.OwnerUsername,
		IsPrivate:     currentWorkspace.IsPrivate,
		NotesID:       currentWorkspace.NotesID,
	}
	return createdWorkspace, nil
}

func (pgdb *PostgresDatabase) DeleteWorkspace(workspace *abstractions.Workspace) error {
	err := validateExistingWorkspace(workspace)
	if err != nil {
		pgdb.logger.Warnw("error in validation workspace to delete", "type", "postgres",
			"output", err.Error(), "time", time.Now().String())
		return err
	}
	foundedWorkspace := &RelationWorkspace{}
	result := pgdb.database.Table("relation_workspaces").Where(
		"id = ?", workspace.Id).First(foundedWorkspace)
	if result.Error != nil {
		pgdb.logger.Warnw("can't find current workspace in the 'relation_workspaces'", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't find current workspace in the 'relation_workspaces: %v", result.Error)
	}
	result = pgdb.database.Table("relation_workspaces").Delete(&RelationWorkspace{}, foundedWorkspace.Id)
	if result.Error != nil {
		pgdb.logger.Warnw("can't delete current workspace from database 'relation_workspaces'",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't delete current workspace from database 'relation_workspaces': %v",
			result.Error)
	}
	return nil
}

func (pgdb *PostgresDatabase) EditWorkspacePrivacy(workspace *abstractions.Workspace) error {
	err := validateExistingWorkspace(workspace)
	if err != nil {
		pgdb.logger.Warnw("error in validation workspace to edit privacy policy", "type", "postgres",
			"output", err.Error(), "time", time.Now().String())
		return err
	}
	foundedWorkspace := &RelationWorkspace{}
	result := pgdb.database.Table("relation_workspaces").Where(
		"id = ?", workspace.Id).First(foundedWorkspace)
	if result.Error != nil {
		pgdb.logger.Warnw("can't update workspace privacy policy in 'relation_workspaces'",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't find current workspace in the 'relation_workspaces: %v", result.Error)
	}
	result = pgdb.database.Model(&RelationUser{}).Where(
		"id = ?", foundedWorkspace.Id).Update("is_private", !foundedWorkspace.IsPrivate)
	if result.Error != nil {
		pgdb.logger.Warnw("can't update workspace privacy policy", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("error by updating privacy policy in 'relation_workspaces': %v",
			result.Error)
	}
	return nil
}

func (pgdb *PostgresDatabase) GetUserNotes(user *abstractions.User) ([]*abstractions.Note, error) {
	currentUser := &RelationUser{}
	result := pgdb.database.Table("relation_users").Where(
		"username = ?", user.GetUsername()).Find(currentUser)
	if result.Error != nil {
		pgdb.logger.Warnw("can't get user in 'relation_users'",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return []*abstractions.Note{}, fmt.Errorf("can't get user in 'relation_users': %v", result.Error)
	}
	var allNotes []RelationNote
	result = pgdb.database.Table("relation_notes").Where(
		"owner_username = ?", currentUser.Username).Find(&allNotes)
	if result.Error != nil {
		pgdb.logger.Warnw("can't get user's notes in 'relation_notes'",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return []*abstractions.Note{}, fmt.Errorf("can't get user's notes in 'relation_notes': %v", result.Error)
	}
	toReturnNotes := make([]*abstractions.Note, len(allNotes))
	for ind, note := range allNotes {
		toAddNote := &abstractions.Note{
			Id:            note.Id,
			WorkspaceID:   note.WorkspaceID,
			OwnerUsername: note.OwnerUsername,
			NoteText:      note.NoteText,
			IsPrivate:     note.IsPrivate,
			Tags:          note.Tags,
			CreatedAt:     note.CreatedAt,
			LastEditedAt:  note.LastEditedAt,
		}
		toReturnNotes[ind] = toAddNote
	}
	return toReturnNotes, nil
}

func (pgdb *PostgresDatabase) GetUserNote(noteID *abstractions.NoteID) (*abstractions.Note, error) {
	currentNote := &RelationNote{}
	result := pgdb.database.Table("relation_notes").Where(
		"id = ?", noteID.GetNoteID()).First(&currentNote)
	if result.Error != nil {
		pgdb.logger.Warnw("can't get user note in 'relation_notes",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return &abstractions.Note{}, fmt.Errorf("can't get user note in 'relation_notes: %v", result.Error)
	}
	toReturnNote := &abstractions.Note{
		Id:            currentNote.Id,
		WorkspaceID:   currentNote.WorkspaceID,
		OwnerUsername: currentNote.OwnerUsername,
		NoteText:      currentNote.NoteText,
		IsPrivate:     currentNote.IsPrivate,
		Tags:          currentNote.Tags,
		CreatedAt:     currentNote.CreatedAt,
		LastEditedAt:  currentNote.LastEditedAt,
	}
	return toReturnNote, nil
}

func (pgdb *PostgresDatabase) CreateNote(note *abstractions.Note) error {
	err := validateNewNote(note)
	if err != nil {
		pgdb.logger.Warnw("error by validating new note",
			"type", "postgres", "output", err.Error(), "time", time.Now().String())
		return fmt.Errorf("error by validating new note: %v", err.Error())
	}
	newNote := &RelationNote{
		WorkspaceID:   note.GetWorkspaceID(),
		OwnerUsername: note.GetOwnerUsername(),
		NoteText:      note.GetNoteText(),
		IsPrivate:     note.GetIsPrivate(),
		Tags:          note.GetTags(),
		CreatedAt:     note.GetCreatedAt(),
		LastEditedAt:  note.GetLastEditedAt(),
	}
	result := pgdb.database.Table("relation_notes").Create(newNote)
	if result.Error != nil {
		pgdb.logger.Warnw("failed on creating new note",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return fmt.Errorf("failed on creating new note: %v", result.Error)
	}
	return nil
}

func (pgdb *PostgresDatabase) DeleteNote(note *abstractions.Note) error {
	err := validateExistingNote(note)
	if err != nil {
		pgdb.logger.Warnw("failed on validating note to delete",
			"type", "postgres", "output", err.Error(), "time", time.Now().String())
		return fmt.Errorf("failed on validating note to delete: %v", err)
	}
	result := pgdb.database.Table("relation_notes").Delete(&RelationNote{}, note.Id)
	if result.Error != nil {
		pgdb.logger.Warnw("failed on deleting note",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return fmt.Errorf("failed on deleting note: %v", result.Error)
	}
	return nil
}

func (pgdb *PostgresDatabase) EditNote(note *abstractions.Note) error {
	existingNote := &RelationNote{}
	result := pgdb.database.Table("relation_notes").Where("id = ?", note.GetId()).First(existingNote)
	if result.Error != nil {
		pgdb.logger.Warnw("failed on validating note to delete",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return fmt.Errorf("failed on validating note to delete: %v", result.Error)
	}
	existingNote.NoteText = note.GetNoteText()
	existingNote.IsPrivate = note.GetIsPrivate()
	existingNote.Tags = note.GetTags()
	existingNote.LastEditedAt = note.GetLastEditedAt()
	result = pgdb.database.Table("relation_notes").Save(existingNote)
	if result.Error != nil {
		pgdb.logger.Warnw("failed on save edited note",
			"type", "postgres", "output", result.Error, "time", time.Now().String())
		return fmt.Errorf("failed on save edited note: %v", result.Error)
	}
	return nil
}

func validateExistingNote(note *abstractions.Note) error {
	if note.GetId() == 0 || note.GetWorkspaceID() == 0 ||
		note.GetOwnerUsername() == "" {
		return fmt.Errorf("wrong field in input note")
	}
	return nil
}

func validateNewNote(note *abstractions.Note) error {
	if note.GetWorkspaceID() == 0 || note.GetNoteText() == "" ||
		note.GetOwnerUsername() == "" || note.GetCreatedAt() == "" ||
		note.GetLastEditedAt() == "" {
		return fmt.Errorf("wrong field of new note")
	}
	return nil
}

func validateNewWorkspace(ws *abstractions.Workspace) error {
	if ws.GetOwnerUsername() == "" {
		return fmt.Errorf("empty ownerUsername of workspace")
	}
	return nil
}

func validateExistingWorkspace(ws *abstractions.Workspace) error {
	if ws.GetOwnerUsername() == "" || ws.GetId() == 0 {
		return fmt.Errorf("empty ownerUsername or empty workspaceID")
	}
	return nil
}
