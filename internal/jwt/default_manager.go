package jwt

var _ Manager[int] = (*DefaultManager[int])(nil)

type DefaultManager[T any] struct {
}

func NewDefaultManager[T any]() *DefaultManager[T] {
	return &DefaultManager[T]{}
}

func (m *DefaultManager[T]) GenerateToken(data T) (string, error) {
	panic("not implemented")
}

func (m *DefaultManager[T]) VerifyToken(token string) (*JClaims[T], error) {
	panic("not implemented")
}
