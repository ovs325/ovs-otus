package pkg

import "time"

func GetDayInterval(date time.Time) (first, last time.Time) {
	first = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	last = first.AddDate(0, 0, 1).Add(-time.Nanosecond)
	return
}

func GetWeekInterval(date time.Time) (first, last time.Time) {
	first = date.AddDate(0, 0, -int(date.Weekday()))
	last = first.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return
}

func GetMonthInterval(date time.Time) (first, last time.Time) {
	first = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	last = first.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return
}
