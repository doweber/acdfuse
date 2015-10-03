package acdfs

import (
	"bazil.org/fuse"
)

type DirEntry interface {
	Inode() uint64
	Name() string
	Type() fuse.DirentType
}
