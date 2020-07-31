package auth

type Auth interface {
	RedirectAuth()
	RequestAccessToken(string) (string, error)
	RequestUser()
}

type AuthUser struct {
	ID        int64
	Username  string
	Avatar    string
}