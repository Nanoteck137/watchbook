package utils

import (
	"database/sql"
	"log"
	"math"
	"strings"

	"github.com/gosimple/slug"
	"github.com/nrednav/cuid2"
)

var CreateId = createIdGenerator(32)
var CreateSmallId = createIdGenerator(8)

var CreateAnimeId = createIdGenerator(8)

var CreateApiTokenId = createIdGenerator(32)

func createIdGenerator(length int) func() string {
	res, err := cuid2.Init(cuid2.WithLength(length))
	if err != nil {
		log.Fatal("Failed to create id generator", "err", err)
	}

	return res
}

func ParseAuthHeader(authHeader string) string {
	splits := strings.Split(authHeader, " ")
	if len(splits) != 2 {
		return ""
	}

	if splits[0] != "Bearer" {
		return ""
	}

	return splits[1]
}

func Slug(s string) string {
	return slug.Make(s)
}

func SplitString(s string) []string {
	tags := []string{}
	if s != "" {
		tags = strings.Split(s, ",")
	}

	return tags
}

func TotalPages(perPage, totalItems int) int {
	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}

func FixSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func Int64PtrToSqlNull(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: *i,
		Valid: true,
	}
}

func Float64PtrToSqlNull(i *float64) sql.NullFloat64 {
	if i == nil {
		return sql.NullFloat64{}
	}

	return sql.NullFloat64{
		Float64: *i,
		Valid:   true,
	}
}

func StringPtrToSqlNull(i *string) sql.NullString {
	if i == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: *i,
		Valid:  true,
	}
}

func NullToDefault[T any](p *T) T {
	var res T

	if p != nil {
		res = *p
	}

	return res
}

func SqlNullToStringPtr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}

	return nil
}

func SqlNullToInt64Ptr(value sql.NullInt64) *int64 {
	if value.Valid {
		return &value.Int64
	}

	return nil
}
