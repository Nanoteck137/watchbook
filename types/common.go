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

func (d WorkDir) MediaDir() MediaDir {
	return MediaDir(path.Join(d.String(), "media"))
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

type Change[T any] struct {
	Value   T
	Changed bool
}
