package discord

import (
	"database/sql"
	"time"
)

type Vehicles struct {
	Classname   string
	Image       string
	DisplayName string
}

type PlayerData struct {
	Id  int
	Uid string

	Name      string
	NickName  sql.NullString
	FirstName sql.NullString
	LastName  sql.NullString

	Cash uint32
	Bank uint32
	RC   uint32

	GroupName      sql.NullString
	GroupLevel     uint16
	GroupLevelName sql.NullString

	InsertTime    time.Time
	LastConnected time.Time
	TotalTime     uint64
}
