package errorUtils

import "errors"

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrSteamUserNotFound = errors.New("steam user not found")

	ErrTicketOpened = errors.New("ticket already opened")

	ErrorNotCached        = errors.New("not cached")
	ErrorCantCacheRedis   = errors.New("can't cache redis")
	ErrorCantCacheMemory  = errors.New("can't cache memory")
	ErrorCantDeleteMemory = errors.New("can't delete memory")
	ErrorCantDeleteRedis  = errors.New("can't delete redis")
)
