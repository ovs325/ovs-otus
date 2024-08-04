package storage

import "fmt"

var (
	ErrDateBusy = fmt.Errorf("данное время уже занято другим событием")
	//
	ErrLostID        = fmt.Errorf("не удалось получить id клиента")
	ErrBadID         = fmt.Errorf("id клиента не является типом int")
	ErrBadFormatTime = fmt.Errorf("неправильный формат даты")
)
