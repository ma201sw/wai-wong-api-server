package golib

type Service interface {
	Base64Client
	IOClient
	JSONClient
	JwtClient
}

// verify interface compliance
var _ Service = (*goLibImpl)(nil)

type goLibImpl struct{}

func New() goLibImpl {
	return goLibImpl{}
}
