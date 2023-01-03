package utils

import "time"

func TimePtr(d time.Time) *time.Time {
	return &d
}
