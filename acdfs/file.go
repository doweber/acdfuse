package acdfs

import (
	"golang.org/x/net/context"

	"bazil.org/fuse"
	_ "bazil.org/fuse/fs/fstestutil"
)

// File implements both Node and Handle for the hello file.
type File struct {
}

const greeting = "hello, world\n"

func (File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(greeting))
	return nil
}

func (File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(greeting), nil
}
