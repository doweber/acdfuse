package acdfs

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Tree struct {
	Root *TreeEntry
}

type TreeEntry struct {
	E    fuse.Dirent
	Mode os.FileMode
	Kids []*TreeEntry
}

func (this *TreeEntry) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = this.E.Inode
	a.Mode = this.Mode

	if this.E.Type == fuse.DT_File {
		a.Size = uint64(len(greeting))
	}

	return nil
}

func (this *TreeEntry) Lookup(ctx context.Context, name string) (fs.Node, error) {
	for _, k := range this.Kids {
		if k.E.Name == name {
			return k, nil
		}
	}
	return nil, fuse.ENOENT
}

func (this *TreeEntry) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirDirs := []fuse.Dirent{}

	for _, k := range this.Kids {
		dirDirs = append(dirDirs, k.E)
	}

	return dirDirs, nil
}

func (this *TreeEntry) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(greeting), nil
}
