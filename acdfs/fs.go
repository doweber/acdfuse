package acdfs

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"bazil.org/fuse/fs"
)

var inodes = make(map[uint64]*TreeEntry)
var inodeKids = make(map[uint64][]*TreeEntry)

// FS implements the hello world file system.
type FS struct{}

var apiClient *http.Client
var endpointCfg *EndpointConfig

func init() {
	auth()
	apiClient = conf.Client(oauth2.NoContext, token)

	endpointCfg = NewEndpointConfig(apiClient)
}

var inodeCnt uint64 = 0

func genInode() uint64 {
	inodeCnt += 1
	return inodeCnt
}

func (this FS) Root() (fs.Node, error) {

	root := GetRootNode(apiClient, endpointCfg)

	rootNode := NewDirEntry(genInode(), "root", getKidsFunc)
	rootNode.CustomId = root.Id

	return rootNode, nil

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

func getKidsFunc(node *TreeEntry) []*TreeEntry {
	if kids, ok := inodeKids[node.E.Inode]; ok {
		// if exists in cache, use that
		return kids
	} else {
		fmt.Println("building cached copy of kids")
		// otherwise load it
		apiList := ListNodes(fmt.Sprintf("nodes/%s/children", node.CustomId), apiClient, endpointCfg)

		kids := []*TreeEntry{}

		for _, v := range apiList.Data {
			var newNode *TreeEntry
			switch v.Kind {
			case "FILE":
				newNode = NewFileEntry(genInode(), v.Name, getContentSize, getContent)
				kids = append(kids, newNode)
				break
			case "FOLDER":
				newNode = NewDirEntry(genInode(), v.Name, getKidsFunc)
				kids = append(kids, newNode)
				break
			}

			newNode.CustomId = v.Id
			inodes[newNode.E.Inode] = newNode
		}

		inodeKids[node.E.Inode] = kids
	}
	return []*TreeEntry{}
}

func getContentSize(node *TreeEntry) uint64 {
	return uint64(len(greeting))
}
func getContent(node *TreeEntry) ([]byte, error) {
	return []byte(greeting), nil
}
