package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type Response struct {
	Message string `json:"message"`
}

func GetDriveService() (*drive.Service, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, err
	}

	client := GetClient(config)

	service, err := drive.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func GetClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		tok = GetTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(oauth2.NoContext, tok)
}

func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	fmt.Println("Paste Authorization code here:")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func ListPDFFiles(service *drive.Service, folderID string) ([]*drive.File, error) {
	q := fmt.Sprintf("'%s' in parents and mimeType='application/pdf'", folderID)
	files, err := service.Files.List().Q(q).Do()
	if err != nil {
		return nil, err
	}
	return files.Files, nil
}

func GetLatestPDFFile(files []*drive.File) (*drive.File, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no PDF files found")
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModifiedTime > files[j].ModifiedTime
	})

	return files[0], nil
}

func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("Unable to marshal token: %v", err)
	}

	err = os.WriteFile(path, tokenJSON, 0600)
	if err != nil {
		log.Fatalf("Unable to write token to file: %v", err)
	}
}