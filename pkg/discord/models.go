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
	LastConnected  time.Time
	InsertTime     time.Time
	Uid            string
	Name           string
	NickName       sql.NullString
	FirstName      sql.NullString
	LastName       sql.NullString
	GroupName      sql.NullString
	GroupLevelName sql.NullString
	TotalTime      uint64
	Id             int
	GroupID        int
	RC             int64
	Cash           int64
	Bank           int64
	GroupLevel     int16
}

type gov struct {
	Gov struct {
		Info struct {
			All int `json:"all"`
		} `json:"info"`
	} `json:"gov"`
}
