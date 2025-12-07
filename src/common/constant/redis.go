package constant

import "time"

const (
	UPLOAD_TTL = 15 * time.Minute
	LOCK_TTL   = 5 * time.Second

	DISTRIBUTED_LOCK_PREFIX = "lock:"
	DEAD_LETTER_QUEUE_NAME  = "dlq:upload-failures"
	UPLOAD_QUEUE_NAME       = "queue:upload_furniture_files"

	MAX_RETRY_NUMBER = 5
)
