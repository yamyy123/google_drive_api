//  package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"sort"

// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// 	"google.golang.org/api/drive/v3"
// )

// func getDriveService() (*drive.Service, error) {
// 	b, err := ioutil.ReadFile("credentials.json")
// 	if err != nil {
// 		fmt.Printf("Unable to read credentials.json file. Err: %v\n", err)
// 		return nil, err
// 	}

// 	config, err := google.ConfigFromJSON(b, drive.DriveScope)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client := getClient(config)

// 	service, err := drive.New(client)

// 	if err != nil {
// 		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
// 		return nil, err
// 	}

// 	return service, err
// }

// func getClient(config *oauth2.Config) *http.Client {
// 	tokFile := "token.json"
// 	tok, err := tokenFromFile(tokFile)
// 	if err != nil {
// 		tok = getTokenFromWeb(config)
// 		saveToken(tokFile, tok)
// 	}
// 	return config.Client(context.Background(), tok)
// }

// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	tok := &oauth2.Token{}
// 	err = json.NewDecoder(f).Decode(tok)
// 	return tok, err
// }

// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
// 	fmt.Printf("Go to the following link in your browser then type the "+
// 		"authorization code: \n%v\n", authURL)

// 	fmt.Println("Paste Authorization code here:")
// 	var authCode string
// 	if _, err := fmt.Scan(&authCode); err != nil {
// 		log.Fatalf("Unable to read authorization code %v", err)
// 	}

// 	tok, err := config.Exchange(context.Background(), authCode)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve token from web %v", err)
// 	}
// 	return tok
// }

// func listPDFFiles(service *drive.Service) ([]*drive.File, error) {
// 	q := "mimeType='application/pdf'"
// 	files, err := service.Files.List().Q(q).Do()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return files.Files, nil
// }

// func getLatestPDFFile(files []*drive.File) (*drive.File, error) {
// 	if len(files) == 0 {
// 		return nil, fmt.Errorf("no PDF files found")
// 	}

// 	// Sort files by modified time in descending order
// 	sort.Slice(files, func(i, j int) bool {
// 		return files[i].ModifiedTime > files[j].ModifiedTime
// 	})

// 	return files[0], nil
// }
// func saveToken(path string, token *oauth2.Token) {
// 	fmt.Printf("Saving credential file to: %s\n", path)
// 	tokenJSON, err := json.Marshal(token)
// 	if err != nil {
// 		log.Fatalf("Unable to marshal token: %v", err)
// 	}

// 	err = ioutil.WriteFile(path, tokenJSON, 0600)
// 	if err != nil {
// 		log.Fatalf("Unable to write token to file: %v", err)
// 	}
// }
// func downloadPDFFile(service *drive.Service, file *drive.File) error {
// 	resp, err := service.Files.Get(file.Id).Download()
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	out, err := os.Create(file.Name)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func main() {
// 	// Get Drive service
// 	srv, err := getDriveService()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve Drive service: %v", err)
// 	}

// 	// List PDF files
// 	files, err := listPDFFiles(srv)
// 	if err != nil {
// 		log.Fatalf("Unable to list PDF files: %v", err)
// 	}

// 	// Get the latest PDF file
// 	latestPDF, err := getLatestPDFFile(files)
// 	if err != nil {
// 		log.Fatalf("Unable to get the latest PDF file: %v", err)
// 	}

// 	// Download the latest PDF file
// 	err = downloadPDFFile(srv, latestPDF)
// 	if err != nil {
// 		log.Fatalf("Unable to download PDF file: %v", err)
// 	}

// 	fmt.Printf("Latest PDF file downloaded: %s\n", latestPDF.Name)
// }
