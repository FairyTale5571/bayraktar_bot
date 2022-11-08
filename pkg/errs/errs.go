package errs

import "errors"

var (
	ErrMongoClientNotInitialized = errors.New("mongo client not initialized")
)
