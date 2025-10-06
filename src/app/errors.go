package app

import "errors"

var (
	ErrValidationFailed   = errors.New("validation failed")
	ErrAlreadyExist       = errors.New("already exist")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrAuthHeaderMissing   = errors.New("authorization header is missing")
	ErrClaimsParsingFailed = errors.New("claims could not be parsed")
	ErrInvalidAuthHeader   = errors.New("invalid authorization header")
	ErrInvalidTokenClaims  = errors.New("invalid token claims")
	ErrInvalidTokenString  = errors.New("invalid token string method")
	ErrInvalidSession      = errors.New("invalid session")
	ErrSessionRevoked      = errors.New("session revoked")

	ErrMinioNotificationEmpty = errors.New("webhook notification contained no records")
)
