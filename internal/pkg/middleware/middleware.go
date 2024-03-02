package middleware

type authMiddleware struct {
}

func NewAuthMiddleware() *authMiddleware {
	return &authMiddleware{}
}
