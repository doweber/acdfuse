package acdfs

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// FS implements the hello world file system.
type FS struct{}

func (this FS) Root() (fs.Node, error) {
	log.Println("FS-Root", this)
	return &foo, nil
	return Dir{
		//Mode:  os.ModeDir | 0555,
		//Entry: fuse.Dirent{Inode: 1, Name: "root", Type: fuse.DT_Dir},
		T: &foo,
	}, nil
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "hello", Type: fuse.DT_File},
}

var foo = TreeEntry{
	E:    fuse.Dirent{Inode: 1, Name: "root", Type: fuse.DT_Dir},
	Mode: os.ModeDir | 0555,
	Kids: []*TreeEntry{
		&TreeEntry{
			E:    fuse.Dirent{Inode: 2, Name: "fooHello", Type: fuse.DT_File},
			Mode: 0444,
			Kids: []*TreeEntry{},
		},
		&TreeEntry{
			E:    fuse.Dirent{Inode: 3, Name: "barHello", Type: fuse.DT_File},
			Mode: 0444,
			Kids: []*TreeEntry{},
		},
	},
}
