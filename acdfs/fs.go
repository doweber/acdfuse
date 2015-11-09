package acdfs

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"bazil.org/fuse/fs"
)

// FS implements the hello world file system.
type FS struct{}

var apiClient *http.Client
var endpointCfg *EndpointConfig

func init() {
	auth()
	apiClient = conf.Client(oauth2.NoContext, token)

	endpointCfg = NewEndpointConfig(apiClient)
}

func (this FS) Root() (fs.Node, error) {

	root := GetRootNode(apiClient, endpointCfg)
	topLevelList := ListNodes(fmt.Sprintf("nodes/%s/children", root.Id), apiClient, endpointCfg)

	kids := []*TreeEntry{}

	getKidsFunc := func() []*TreeEntry {
		return []*TreeEntry{}
	}

	for i, v := range topLevelList.Data {
		switch v.Kind {
		case "FILE":
			kids = append(kids, NewFileEntry(uint64(i+1), v.Name, getContentSize, getContent))
			break
		case "FOLDER":
			kids = append(kids, NewDirEntry(uint64(i+1), v.Name, getKidsFunc))
			break
		}
	}

	return NewDirEntry(1, "root", func() []*TreeEntry {
		return kids
	}), nil

	/*
		return NewDirEntry(1, "root", func() []*TreeEntry {
			return []*TreeEntry{
				NewDirEntry(2, "foodir", func() []*TreeEntry {
					return []*TreeEntry{
						NewFileEntry(3, "fooHello", getContentSize, getContent),
					}
				}),
				NewFileEntry(3, "fooHello", getContentSize, getContent),
				NewFileEntry(4, "barHello", getContentSize, getContent),
			}
		}), nil
	*/
}

var greeting = "hello, world\n"

func getContentSize() uint64 {
	return uint64(len(greeting))
}
func getContent() ([]byte, error) {
	return []byte(greeting), nil
}
