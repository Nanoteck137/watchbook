package types

import (
	"errors"
	"path"
)

type WorkDir string

func (d WorkDir) String() string {
	return string(d)
}

func (d WorkDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d WorkDir) SetupFile() string {
	return path.Join(d.String(), "setup")
}

type AdminStatus string

const (
	AdminStatusNotFixed AdminStatus = "not-fixed"
	AdminStatusFixed    AdminStatus = "fixed"
)

func IsValidAdminStatus(l AdminStatus) bool {
	switch l {
	case AdminStatusNotFixed,
		AdminStatusFixed:
		return true
	}

	return false
}

func ValidateAdminStatus(val any) error {
	if s, ok := val.(string); ok {
		if s == "" {
			return nil
		}

		t := AdminStatus(s)
		if !IsValidAdminStatus(t) {
			return errors.New("invalid admin status")
		}
	} else if p, ok := val.(*string); ok {
		if p == nil {
			return nil
		}

		s := *p
		if s == "" {
			return nil
		}

		t := AdminStatus(s)
		if !IsValidAdminStatus(t) {
			return errors.New("invalid admin status")
		}
	} else {
		return errors.New("expected string")
	}

	return nil
}
