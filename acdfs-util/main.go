package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"

	//"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/Reisender/acdfuse/acdfs"
)

var configPath = "./config.json"
var tokenPath = "./token.json"
var token = &oauth2.Token{}
var conf = &oauth2.Config{
//ClientID:     "",
//ClientSecret: "",
//Scopes:       []string{"clouddrive:read_all"},
//Endpoint: oauth2.Endpoint{
//AuthURL:  "https://www.amazon.com/ap/oa",
//TokenURL: "https://api.amazon.com/auth/o2/token",
//},
//RedirectURL: "https://www.google.com/",
}

func main() {
	app := cli.NewApp()
	app.Name = "acdfs-util"
	app.Usage = "utility to help configure acdfs"
	app.Commands = []cli.Command{
		{
			Name:    "auth",
			Aliases: []string{"a"},
			Usage:   "authorize acdfs",
			Action:  auth,
		},
		{
			Name:   "save-config",
			Usage:  "save out the config",
			Action: SaveConfig,
		},
		{
			Name:   "test",
			Usage:  "test the config",
			Action: TestConfig,
		},
	}

	app.Run(os.Args)
}

func TestConfig(c *cli.Context) {
	auth(c)

	client := conf.Client(oauth2.NoContext, token)
	cfg := acdfs.NewEndpointConfig(client)
	fmt.Println(cfg)

	// now try the metadata url
	list := acdfs.ListNodes("nodes?filters=isRoot:true", client, cfg)

	if len(list.Data) != 1 {
		log.Fatal("no root node")
	}

	topLevelList := acdfs.ListNodes(fmt.Sprintf("nodes/%s/children", list.Data[0].Id), client, cfg)
	fmt.Println("list length:", len(topLevelList.Data))

}

func auth(c *cli.Context) {

	err := acdfs.LoadConsumerConfig(configPath, conf)
	if err != nil {
		fmt.Println("config file not found at", configPath)
		return
	}

	// see if token exists
	if err := acdfs.LoadAccessToken(tokenPath, token); err != nil {
		// no token or problem with it so go get one

		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v", url)

		print("\nenter code: ")
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}

		token, err = conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Fatal(err)
		}

		// save the token for later use
		SaveToken()
	}

	fmt.Println("authenticated")
	//client := conf.Client(oauth2.NoContext, token)
	//client.Get("...")
}

func SaveToken() {
	var b []byte
	b, _ = json.Marshal(token)
	ioutil.WriteFile(tokenPath, b, 0600)
}

// Save the consumer key and secret in from the config file
func SaveConfig(c *cli.Context) {
	var b []byte
	b, _ = json.Marshal(conf)
	ioutil.WriteFile(configPath, b, 0600)
}
