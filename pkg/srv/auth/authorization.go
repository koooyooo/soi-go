package auth

type AuthInfo struct{}

type AuthRepo interface {
	Register(AuthInfo) error
	Reset(email string) error
	Authorize(id string) error
}
