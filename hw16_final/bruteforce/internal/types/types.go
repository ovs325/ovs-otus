package types

import "time"

type IsBlack bool

type RequestParams map[string]string

type AllBucketParams struct {
	ID           string        `json:"id"`           // Имя корзины
	Capacity     int64         `json:"capacity"`     // Емкость в каплях
	FreeCapacity int64         `json:"freeCapacity"` // Свободная емкость корзины в каплях.
	DropsSum     int64         `json:"dropsSum"`     // Объем недокапанного запаса в каплях.
	LeakageRate  float64       `json:"leakageRate"`  // Скорость утечки (капли.сек)
	Expiry       time.Time     `json:"expiry"`       // время опустошения корзины.
	DeadLine     time.Duration `json:"deadLine"`     // Промежуток времени до опустошения корзины.
	QPriority    int           `json:"priority"`     // Приоритет корзины в очереди кучи
	IsReset      bool          `json:"isReset"`      // Корзина должна быть обнулена?
	IsAllowed    bool          `json:"isAllowed"`    // Разрешить запрос?
}

type AllGroupData struct {
	LenMap        int                        `json:"lenMap"`        // Свободная емкость корзины в каплях.
	FirstDrops    string                     `json:"firstDrops"`    // Объем недокапанного запаса в каплях.
	BacketsParams map[string]AllBucketParams `json:"backetsParams"` // Разрешить запрос?
}

type AllGroupsData struct {
	Elapsed      string                        `json:"elapsed"`
	LenGroupsMap int                           `json:"lenGroupsMap"` // Свободная емкость корзины в каплях.
	LenSList     int                           `json:"lenSList"`
	SpecList     map[string]map[string]IsBlack `json:"specList"`
	GroupsParams map[string]AllGroupData       `json:"groupsParams"`
}
