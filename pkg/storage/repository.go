package storage

type Bucket string

const (
	UsersMuted   Bucket = "muted_users"
	LastWsUpdate Bucket = "last_ws_update"
	Cache        Bucket = "cache"
)
