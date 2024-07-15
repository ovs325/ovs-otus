package storage

import "fmt"

var (
	ErrDateBusy error = fmt.Errorf("данное время уже занято другим событием")
)
