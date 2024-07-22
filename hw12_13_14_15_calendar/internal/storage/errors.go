package storage

import "fmt"

var ErrDateBusy = fmt.Errorf("данное время уже занято другим событием")
