package acdfs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
