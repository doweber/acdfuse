package acdfs

import (
	"golang.org/x/net/context"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	T *TreeEntry
}

func (this Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = this.T.E.Inode
	a.Mode = this.T.Mode
	return nil
}

func (this Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "hello" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

func (this Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return dirDirs, nil
}
