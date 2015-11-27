package acdfs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/sethgrid/multibar"
	"golang.org/x/oauth2"
)

type EndpointConfig struct {
	ContentUrl     string `json:"contentUrl"`
	CustomerExists bool   `json:"customerExists"`
	MetadataUrl    string `json:"metadataUrl"`
}

var configPath = "./config.json"
var tokenPath = "./token.json"
var token = &oauth2.Token{}
var conf = &oauth2.Config{}

func auth() {
	// load consumer config
	if err := LoadConsumerConfig(configPath, conf); err != nil {
		fmt.Println("config file not found at", configPath)
		return
	}

	// load access token exists
	if err := LoadAccessToken(tokenPath, token); err != nil {
		// no token or problem with it so go get one
		fmt.Println("access token file not found at", tokenPath)
		return
	}
}

func NewEndpointConfig(apiClient *http.Client) *EndpointConfig {
	resp, err := apiClient.Get("https://drive.amazonaws.com/drive/v1/account/endpoint")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	cfg := &EndpointConfig{}
	if err := json.Unmarshal(body, cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

// Load the consumer key and secret in from the config file
func LoadConsumerConfig(configPath string, config *oauth2.Config) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("config file not found: %s", configPath)
		return err
	}

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, config); err != nil {
		return err
	}

	return nil
}

func LoadAccessToken(tokenPath string, token *oauth2.Token) error {
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		fmt.Printf("token file not found: %s\n", tokenPath)
		fmt.Println(err)
		return err
	}

	b, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, token); err != nil {
		return err
	}

	return nil
}

func LoadMetadata(client *http.Client, cfg *EndpointConfig) (nodes []Metadata, err error) {
	nodesPath := "./nodes.json"
	if err = LoadMetadataFromFile(&nodes, nodesPath); err == nil {
		return
	}

	// if we couldn't load from file, load from API
	if len(nodes) == 0 {

		// now try the metadata url
		list, err := ListNodes("nodes", client, cfg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("List Count:", list.Count)

		nodes = append(nodes, list.Data...)

		var progressBars, _ = multibar.New()
		progressBars.Println("getting the pages")
		barProgress1 := progressBars.MakeBar(list.Count, "loading...")
		go progressBars.Listen()

		getPage(list.NextToken, client, cfg, func(newList *MetadataPage) {
			barProgress1(len(nodes))
			nodes = append(nodes, newList.Data...)
		})

		SaveMetadataToFile(&nodes, nodesPath)
	}

	return
}

func LoadMetadataFromFile(nodes *[]Metadata, nodesPath string) error {
	if _, err := os.Stat(nodesPath); os.IsNotExist(err) {
		fmt.Printf("nodes file not found: %s\n", nodesPath)
		fmt.Println(err)
		return err
	}

	b, err := ioutil.ReadFile(nodesPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, nodes); err != nil {
		return err
	}

	return nil
}

func SaveMetadataToFile(nodes *[]Metadata, nodesPath string) {
	fmt.Println("saving nodes")
	var b []byte
	b, _ = json.Marshal(nodes)
	ioutil.WriteFile(nodesPath, b, 0600)
	fmt.Println("saved nodes")
}

func getPage(nextToken string, client *http.Client, cfg *EndpointConfig, callback func(*MetadataPage)) {
	list, err := ListNodes(fmt.Sprintf("nodes?startToken=%s&count=800", nextToken), client, cfg)
	if err != nil {
		log.Println(err)
		return
	}
	callback(list)

	if len(list.Data) == 200 {
		getPage(list.NextToken, client, cfg, callback)
	}
}
