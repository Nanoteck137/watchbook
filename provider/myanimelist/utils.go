package myanimelist

import (
	"fmt"
	"time"
)

const DateLayout = "Jan _2, 2006"

func ParseDate(s string) (time.Time, error) {
	return time.Parse(DateLayout, s)
}

func formatDate(t time.Time) string {
	// TODO(patrik): Change with t.Format()
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
