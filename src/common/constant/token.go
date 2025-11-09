package constant

import "time"

const (
	ACCESS_TOKEN_TTL  = 10 * time.Minute
	REFRESH_TOKEN_TTL = 24 * time.Hour

	REFRESH_TOKEN_COOKIE = "refresh-token"
	USER_CLAIMS          = "user_claims"
)
