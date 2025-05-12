package mal

import (
	"fmt"
	"time"
)

const DateLayout = "Jan _2, 2006"

func parseDate(s string) (time.Time, error) {
	return time.Parse(DateLayout, s)
}

func formatDate(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
