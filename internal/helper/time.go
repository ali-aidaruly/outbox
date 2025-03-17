package helper

import "time"

func MustParseTime(layout, t string) time.Time {
	res, _ := time.Parse(layout, t)

	return res
}
