package auth

type Auth interface {
	RedirectAuth()
	RequestAccessToken()
	RequestUser()
}
