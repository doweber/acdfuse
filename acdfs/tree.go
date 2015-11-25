package acdfs

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type KidsCallbackFunc func(*TreeEntry) []*TreeEntry
type SizeCallbackFunc func(*TreeEntry) uint64
type ContentCallbackFunc func(*TreeEntry) ([]byte, error)

func KidsFunc(t *TreeEntry) []*TreeEntry {
	return t.Kids
}

type TreeEntry struct {
	E               fuse.Dirent
	Mode            os.FileMode
	CustomId        string
	Kids            []*TreeEntry
	kidsCallback    KidsCallbackFunc
	sizeCallback    SizeCallbackFunc
	contentCallback ContentCallbackFunc
}

func NewDirEntry(inode uint64, name string, kidsCallback KidsCallbackFunc) *TreeEntry {
	return &TreeEntry{
		E:            fuse.Dirent{Inode: inode, Name: name, Type: fuse.DT_Dir},
		Mode:         os.ModeDir | 0555,
		kidsCallback: kidsCallback,
	}
}
func NewFileEntry(inode uint64, name string, size SizeCallbackFunc, content ContentCallbackFunc) *TreeEntry {
	return &TreeEntry{
		E:               fuse.Dirent{Inode: inode, Name: name, Type: fuse.DT_File},
		Mode:            0444,
		Kids:            []*TreeEntry{},
		sizeCallback:    size,
		contentCallback: content,
	}
}

func (this *TreeEntry) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = this.E.Inode
	a.Mode = this.Mode

	if this.E.Type == fuse.DT_File {
		a.Size = this.sizeCallback(this)
	}

	return nil
}

func (this *TreeEntry) Lookup(ctx context.Context, name string) (fs.Node, error) {
	this.Kids = this.kidsCallback(this)
	for _, k := range this.Kids {
		if k.E.Name == name {
			return k, nil
		}
	}
	return nil, fuse.ENOENT
}

func (this *TreeEntry) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirDirs := []fuse.Dirent{}
	this.Kids = this.kidsCallback(this)

	for _, k := range this.Kids {
		dirDirs = append(dirDirs, k.E)
	}

	return dirDirs, nil
}

func (this *TreeEntry) ReadAll(ctx context.Context) ([]byte, error) {
	return this.contentCallback(this)
}
