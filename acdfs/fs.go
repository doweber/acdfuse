package acdfs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

var inodes = make(map[uint64]*TreeEntry)
var inodeKids = make(map[uint64][]*TreeEntry)

// FS implements the hello world file system.
type FS struct{}

var apiClient *http.Client
var endpointCfg *EndpointConfig

func initApi() {
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
	if apiClient == nil {
		initApi()
	}

	root := GetRootNode(apiClient, endpointCfg)
	rootNode := NewDirEntry(genInode(), "root", KidsFunc)
	rootNode.Kids = []*TreeEntry{}
	rootNode.CustomId = root.Id

	nodes, _ := LoadMetadata(apiClient, endpointCfg)

	// create entire tree structure from metadata nodes
	tNodes := make(map[string]*TreeEntry)
	for _, n := range nodes {
		switch n.Kind {
		case FILE:
			tNodes[n.Id] = NewFileEntry(genInode(), n.Name, getContentSize, readContent)
			tNodes[n.Id].CustomId = n.Id
			tNodes[n.Id].Metadata = n
		case FOLDER:
			if n.Id == root.Id {
				tNodes[n.Id] = rootNode
				tNodes[n.Id].Metadata = *root
			} else {
				tNodes[n.Id] = NewDirEntry(genInode(), n.Name, KidsFunc)
				tNodes[n.Id].Kids = []*TreeEntry{}
				tNodes[n.Id].CustomId = n.Id
				tNodes[n.Id].Metadata = n
			}
		}
	}

	// fill in the links to the kids
	for _, n := range nodes {
		if n.Kind == FOLDER || n.Kind == FILE {
			for _, p := range n.Parents {
				tNodes[p].Kids = append(tNodes[p].Kids, tNodes[n.Id])
			}
		}
	}

	fmt.Println("root kids:", len(rootNode.Kids))

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

func getContentSize(node *TreeEntry) uint64 {
	return node.Metadata.ContentProperties.Size
}

func readBytes(r io.Reader, p []byte, cnt int) (readCount int, err error) {
	readCount, err = r.Read(p)

	if readCount < cnt && err != io.EOF {
		newCount := cnt - readCount
		n := make([]byte, newCount)
		newCount, err = readBytes(r, n, newCount)
		p = append(p[:readCount], n[:newCount]...)
		readCount += newCount
	}

	return
}

func readContent(node *TreeEntry, ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) (err error) {
	var p []byte = make([]byte, req.Size)

	// see if we already have an open body
	if node.HttpResponse == nil {
		fmt.Printf("starting to download %s\n", node.E.Name)
		b, _ := json.Marshal(node.Metadata)
		fmt.Println(string(b))
		// get the reader
		node.HttpResponse, err = DownloadContent(node.CustomId, apiClient, endpointCfg)
		if err != nil {
			return err
		}
	}

	cnt, err := readBytes(node.HttpResponse.Body, p, req.Size)
	if err == io.EOF {
		fmt.Println("GOT EOF!!!")
		err = nil
		defer func() {
			node.HttpResponse.Body.Close()
			node.HttpResponse = nil
		}()
	}

	resp.Data = p[:cnt]

	return
}
func getContent(node *TreeEntry) ([]byte, error) {
	resp, err := DownloadContent(node.CustomId, apiClient, endpointCfg)
	defer resp.Body.Close()
	var p []byte
	cnt, err := resp.Body.Read(p)
	fmt.Printf("read %d bytes\n", cnt)
	//body, err := ioutil.ReadAll(resp.Body)
	//return body, err
	return p, err
}
