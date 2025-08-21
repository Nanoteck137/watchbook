package kvstore

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Store map[string]string

func (kv Store) Serialize() (string, error) {
    b, err := json.Marshal(kv)
    if err != nil {
        return "", err
    }

    return string(b), nil
}

func Deserialize(data string) (Store, error) {
    kv := make(Store)
    if data == "" {
        return kv, nil
    }

    err := json.Unmarshal([]byte(data), &kv)
	if err != nil {
		return nil, err
	}

    return kv, nil
}

func (kv Store) Value() (driver.Value, error) {
	return kv.Serialize()
}

func (kv *Store) Scan(src any) error {
	if src == nil {
		*kv = make(Store)
		return nil
	}

	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("KVStore: cannot scan type %T", src)
	}

	r, err := Deserialize(str)
	if err != nil {
		return fmt.Errorf("KVStore: failed to deserialize store: %w", err)
	}

	*kv = r

	return nil
}
