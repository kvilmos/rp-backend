package request

type MinioNotification struct {
	Records []Record `json:"Records"`
}

type Record struct {
	EventName string `json:"eventName"`
	S3        S3     `json:"s3"`
}

type S3 struct {
	Bucket Bucket `json:"bucket"`
	Object Object `json:"object"`
}

type Bucket struct {
	Name string `json:"name"`
}

type Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

type UploadNotification struct {
	BucketName string
	ObjectKey  string
	RetryCount int
}
