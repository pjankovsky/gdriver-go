package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const DriveClientTimeout = 30 * time.Minute

var driveClient *drive.Service
var driveClientEndtime time.Time

func getDrive() *drive.Service {
	if driveClientEndtime.After(time.Now()) {
		driveClient = nil
	}

	if driveClient != nil {
		return driveClient
	}
	buffer, err := ioutil.ReadFile("etc/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(buffer, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	driveClient, err = drive.New(getClient(config))
	if err != nil {
		log.Fatalf("Unable to retrieve drive client: %v", err)
	}

	driveClientEndtime = time.Now().Add(time.Duration(DriveClientTimeout))

	return driveClient
}

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "etc/token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authUrl)
	var authCode string
	_, err := fmt.Scan(&authCode)
	if err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(file string, tok *oauth2.Token) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to save oauth token: %v", err)
	}
	err = json.NewEncoder(f).Encode(tok)
	if err != nil {
		log.Fatalf("Unable to save oauth token: %v", err)
	}
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
