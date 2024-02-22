package converter

import (
	"strings"
)

func MakeTimeStamp(date string, time string) string {

	date = strings.Replace(date, "/", "-", -1)
	time = strings.TrimSuffix(time, ".000")
	return date + " " + time
}
