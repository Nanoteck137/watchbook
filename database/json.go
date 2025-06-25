package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// TODO(patrik): Move
type JsonColumn[T any] struct {
	Data  T
	Valid bool
}

func (j *JsonColumn[T]) Scan(src any) error {
	var res T

	if src == nil {
		j.Data = res
		j.Valid = false
		return nil
	}

	switch value := src.(type) {
	case string:
		err := json.Unmarshal([]byte(value), &j.Data)
		if err != nil {
			return err
		}

		j.Valid = true
	case []byte:
		err := json.Unmarshal(value, &j.Data)
		if err != nil {
			return err
		}

		j.Valid = true
	default:
		return fmt.Errorf("unsupported type %T", src)
	}

	return nil
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	raw, err := json.Marshal(j.Data)
	return raw, err
}

// func (j *JsonColumn[T]) Get() *T {
// 	return j.Val
// }
