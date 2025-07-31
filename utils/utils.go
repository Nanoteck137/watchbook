package utils

import (
	"cmp"
	"database/sql"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"

	"slices"

	"github.com/gosimple/slug"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nrednav/cuid2"
)

var CreateId = createIdGenerator(32)
var CreateSmallId = createIdGenerator(8)

var CreateMediaId = createIdGenerator(12)
var CreateCollectionId = createIdGenerator(8)

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

func MediaUserListPtrToSqlNull(i *types.MediaUserList) sql.NullString {
	if i == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: string(*i),
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

func SqlNullToFloat64Ptr(value sql.NullFloat64) *float64 {
	if value.Valid {
		return &value.Float64
	}

	return nil
}

func Min[T cmp.Ordered](value T, min T) T {
	if value < min {
		return min
	}

	return value
}

func Clamp[T cmp.Ordered](value T, min T, max T) T {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func TransformStringSlug(s string) string {
	s = strings.TrimSpace(s)
	return Slug(s)
}

func TransformSlugArray(arr []string) []string {
	res := make([]string, 0, len(arr))

	for _, t := range arr {
		t = strings.TrimSpace(t)
		s := Slug(t)
		if s != "" {
			res = append(res, s)
		}
	}

	return res
}

func FixNilArrayToEmpty[T any](a []T) []T {
	if a == nil {
		return []T{}
	}

	return a
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func RoundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}

func CreateUrlBase(addr, path string, query url.Values) (*url.URL, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	u.Path = path

	if query != nil {
		params := u.Query()
		for k, v := range query {
			params[k] = v
		}
		u.RawQuery = params.Encode()
	}

	return u, nil
}

func ExtractNumber(s string) int {
	n := ""
	for _, c := range s {
		if unicode.IsDigit(c) {
			n += string(c)
		} else {
			break
		}
	}

	if len(n) == 0 {
		return 0
	}

	i, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return 0
	}

	return int(i)
}

// NOTE(patrik): From https://stackoverflow.com/questions/19374219/how-to-find-the-difference-between-two-slices-of-strings
func SliceDifference[S ~[]E, E comparable](slice1 S, slice2 S) S {
	var diff S

	for i := range 2 {
		for _, s1 := range slice1 {
			found := slices.Contains(slice2, s1)
			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func DownloadImage(url, outDir, name string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to send http get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download unsuccessful: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("failed to parse Content-Type: %w", err)
	}

	ext := ""
	switch mediaType {
	case "image/png":
		ext = ".png"
	case "image/jpeg":
		ext = ".jpeg"
	default:
		return "", fmt.Errorf("unsupported media type: %s", mediaType)
	}

	out := path.Join(outDir, name+ext)

	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy response body to file: %w", err)
	}

	return out, nil
}
