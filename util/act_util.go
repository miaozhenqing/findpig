package util

import (
	"fmt"
	"time"
)

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func SafeGet[T any](s []T, index int) T {
	var defaultValue T
	if index >= 0 && index < len(s) {
		return s[index]
	}
	return defaultValue
}

func DaysBetween(t1, t2 time.Time) int {
	t1 = t1.UTC()
	t2 = t2.UTC()

	days := int(t1.Sub(t2).Hours() / 24)
	if days < 0 {
		days = -days
	}
	return days
}

const (
	timeBase = "2006-01-02 15:04:05"
)

func CurrentTime() time.Time {
	return time.Now()
}
func IsTimeBetween(start time.Time, end time.Time) bool {
	currentTime := CurrentTime()
	return currentTime.Before(end) && currentTime.After(start)
}
func ParseTime(timeStr string) time.Time {
	parse, err := time.Parse(timeBase, timeStr)
	if err != nil {
		panic(fmt.Sprintf("unexpected timeStr: %v", timeStr))
	}
	return parse
}
