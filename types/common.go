package types

import "path"

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

func (d WorkDir) ImagesDir() string {
	return path.Join(d.String(), "images")
}

func (d WorkDir) ImagesEntriesDir() string {
	return path.Join(d.ImagesDir(), "entries")
}

type Change[T any] struct {
	Value   T
	Changed bool
}
