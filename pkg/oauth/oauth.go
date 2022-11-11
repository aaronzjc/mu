package oauth

type OAuth interface {
	Type() string
	RedirectAuth() string
	RequestAccessToken(string) (string, error)
	RequestUser(string) (User, error)
}

type User struct {
	ID       int64
	Username string
	Nickname string
	Avatar   string
}
