package acdfs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetRootNode(client *http.Client, cfg *EndpointConfig) *Metadata {
	list := ListNodes("nodes?filters=isRoot:true", client, cfg)
	if len(list.Data) != 1 {
		log.Fatal("no root node")
	}
	return &list.Data[0]
}

func ListNodes(urlRequest string, client *http.Client, cfg *EndpointConfig) *MetadataList {
	resp, err := client.Get(fmt.Sprintf("%s/%s", cfg.MetadataUrl, urlRequest))
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	list := &MetadataList{}
	if err := json.Unmarshal(body, list); err != nil {
		log.Fatal(err)
	}

	fmt.Println("list length:", len(list.Data))

	return list
}
