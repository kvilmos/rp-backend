package constant

import "time"

const (
	UPLOAD_STATE_TTL          = 30 * time.Minute
	LOCK_TTL                  = 5 * time.Second
	DISTRIBUTED_LOCK_PREFIX   = "lock:"
	UPLOAD_CHANNEL            = "upload-updated"
	DEAD_LETTER_QUEUE_CHANNEL = "dlq:upload-failures"
)
