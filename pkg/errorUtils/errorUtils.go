package errorUtils

import "errors"

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrSteamUserNotFound = errors.New("steam user not found")
)
