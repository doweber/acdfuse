package acdfs

import (
	"bazil.org/fuse/fs"
)

// FS implements the hello world file system.
type FS struct{}

func (FS) Root() (fs.Node, error) {
	return Dir{}, nil
}
