package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

func RandStringRune(n int) string {
	letters := []rune("ABEIKMHOPCTXZ")
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandInt() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	return rand.Intn(max-min+1) + min
}

func GeneratePlateNumber() string {
	return fmt.Sprintf("DS %d %v", RandInt(), RandStringRune(2))
}

func MinutesToDate(minutes uint64) string {
	return fmt.Sprintf("%d дней %d часов %d минут", minutes/1440, minutes%1440/60, minutes%60)
}

func secondsToDate(seconds uint64) string {
	return fmt.Sprintf("%d дней %d часов %d минут", seconds/3600, (seconds%3600)/60, (seconds%3600)%60)
}
