package acdfs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetRootNode(client *http.Client, cfg *EndpointConfig) *Metadata {
	list, err := ListNodes("nodes?filters=isRoot:true", client, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if len(list.Data) != 1 {
		log.Fatal("no root node")
	}
	return &list.Data[0]
}

func ListNodes(urlRequest string, client *http.Client, cfg *EndpointConfig) (list *MetadataPage, err error) {
	resp, err := client.Get(fmt.Sprintf("%s/%s", cfg.MetadataUrl, urlRequest))
	if err != nil {
		return
	}
	if resp.StatusCode == 500 {
		log.Println("500 error, retry")
		return ListNodes(urlRequest, client, cfg)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	list = &MetadataPage{}
	err = json.Unmarshal(body, list)

	return
}

func DownloadContent(nodeId string, client *http.Client, cfg *EndpointConfig) (*http.Response, error) {
	resp, err := client.Get(fmt.Sprintf("%s/nodes/%s/content", cfg.ContentUrl, nodeId))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Error", resp.Status, string(body))
	}

	return resp, nil
}
