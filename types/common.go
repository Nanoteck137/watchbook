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

func (d WorkDir) AnimesDir() AnimeDir {
	return AnimeDir(path.Join(d.String(), "animes"))
}

type AnimeDir string

func (d AnimeDir) String() string {
	return string(d)
}

func (d AnimeDir) ImagesDir() string {
	return path.Join(d.String(), "images")
}

func (d AnimeDir) AnimeImageDir(id string) string {
	return path.Join(d.ImagesDir(), id)
}

type Change[T any] struct {
	Value   T
	Changed bool
}
