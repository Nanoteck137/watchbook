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

func (d WorkDir) MediaDir() MediaDir {
	return MediaDir(path.Join(d.String(), "media"))
}

func (d WorkDir) CacheDir() string {
	return path.Join(d.String(), "cache")
}

func (d WorkDir) CacheProvidersDir() string {
	return path.Join(d.CacheDir(), "providers")
}

func (d WorkDir) CacheProviderDir(providerName string) string {
	return path.Join(d.CacheProvidersDir(), providerName)
}

type MediaDir string

func (d MediaDir) String() string {
	return string(d)
}

func (d MediaDir) ImagesDir() string {
	return path.Join(d.String(), "images")
}

func (d MediaDir) MediaImageDir(mediaId string) string {
	return path.Join(d.ImagesDir(), mediaId)
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
