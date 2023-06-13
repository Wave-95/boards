package test

import "github.com/Wave-95/boards/backend-core/internal/jwt"

// NewJWTService returns a JWT service initialized with a secret and 24 hr expiration
func NewJWTService() jwt.Service {
	jwtSecret := "fake_jwt_secret"
	jwtExp := 24 // 24 hrs
	return jwt.New(jwtSecret, jwtExp)
}
