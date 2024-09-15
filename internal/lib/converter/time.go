package converter

import "time"

func Time(t time.Time) string {
	return t.Format(time.RFC3339)
}
