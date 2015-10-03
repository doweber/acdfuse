package acdfs

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type TreeEntry struct {
	E    fuse.Dirent
	Mode os.FileMode
	Kids []*TreeEntry
}

func NewDirEntry(inode uint64, name string, kids []*TreeEntry) *TreeEntry {
	return &TreeEntry{
		E:    fuse.Dirent{Inode: inode, Name: name, Type: fuse.DT_Dir},
		Mode: os.ModeDir | 0555,
		Kids: kids,
	}
}
func NewFileEntry(inode uint64, name string) *TreeEntry {
	return &TreeEntry{
		E:    fuse.Dirent{Inode: inode, Name: name, Type: fuse.DT_File},
		Mode: 0444,
		Kids: []*TreeEntry{},
	}
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
