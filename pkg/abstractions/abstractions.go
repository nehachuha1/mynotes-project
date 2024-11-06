package abstractions

type User struct {
	Username string
	Email    string
	Initials string
	Telegram string
}

type Session struct {
	SessionID string
	Username  string
}
