package acdfs

import (
	"log"

	"bazil.org/fuse/fs"
)

// FS implements the hello world file system.
type FS struct{}

func (this FS) Root() (fs.Node, error) {
	log.Println("FS-Root", this)
	return foo, nil
}

var greeting = "hello, world\n"

var foo = NewDirEntry(1, "root", []*TreeEntry{
	NewDirEntry(2, "foodir", []*TreeEntry{
		NewFileEntry(3, "fooHello"),
	}),
	NewFileEntry(3, "fooHello"),
	NewFileEntry(4, "barHello"),
})
