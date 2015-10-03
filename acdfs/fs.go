package acdfs

import "bazil.org/fuse/fs"

// FS implements the hello world file system.
type FS struct{}

func (this FS) Root() (fs.Node, error) {
	return NewDirEntry(1, "root", []*TreeEntry{
		NewDirEntry(2, "foodir", []*TreeEntry{
			NewFileEntry(3, "fooHello", getContentSize, getContent),
		}),
		NewFileEntry(3, "fooHello", getContentSize, getContent),
		NewFileEntry(4, "barHello", getContentSize, getContent),
	}), nil
}

var greeting = "hello, world\n"

func getContentSize() uint64 {
	return uint64(len(greeting))
}
func getContent() ([]byte, error) {
	return []byte(greeting), nil
}
