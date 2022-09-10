package models

import "time"

type Error struct {
	Err string `json:"error"`
}

type Government struct {
	Gov struct {
		Rule struct {
			President string
			Poor      string
			Bank      string
			Tax       int
			Credit    int
			Legal     bool
			Slavery   bool
		}
		Info struct {
			Cop int
			Ems int
			Reb int
			Civ int
			All int
		}
		Prem struct {
			Police int
			Ems    int
			Taxi   int
			Press  int
		}
		DropGifts struct {
			TotalGifts  int
			ActiveGifts int
		}
	}
}

type Economy struct {
	ResourceName     string
	Localize         string
	Price            int
	MaxPrice         int
	MinPrice         int
	DownPricePerItem float64
	RandomDownPrice  bool
	RandomMax        int
	RandomMin        int
	Illegal          bool
	Influenced       string
	LastUpdate       time.Time
}
