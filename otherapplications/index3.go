package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func getDriveService() (*drive.Service, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, err
	}

	client := getClient(config)

	service, err := drive.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(oauth2.NoContext, tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
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

func listPDFFiles(service *drive.Service, folderID string) ([]*drive.File, error) {
	q := fmt.Sprintf("'%s' in parents and mimeType='application/pdf'", folderID)
	files, err := service.Files.List().Q(q).Do()
	if err != nil {
		return nil, err
	}
	return files.Files, nil
}

func getLatestPDFFile(files []*drive.File) (*drive.File, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no PDF files found")
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModifiedTime > files[j].ModifiedTime
	})

	return files[0], nil
}

func saveToken(path string, token *oauth2.Token) {
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

func downloadFile(service *drive.Service, fileID, downloadPath string) error {
	resp, err := service.Files.Get(fileID).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	e := echo.New()

	e.GET("/download-latest-pdf", func(c echo.Context) error {
		log.Println("google drive service called")
		srv, err := getDriveService()
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to retrieve Drive service: %v", err))
		}

		// Replace "YOUR_FOLDER_ID" with the actual ID of the folder you want to search in
		folderID := "1kM5D_O2m4X99utS24K4TlYmXzNm9Mopf"

		files, err := listPDFFiles(srv, folderID)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to list PDF files: %v", err))
		}

		latestPDF, err := getLatestPDFFile(files)
		if err != nil {
			if err.Error() == "no PDF files found" {
				return c.String(http.StatusNotFound, "No PDF files found in the specified folder")
			}
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to get the latest PDF file: %v", err))
		}

		// Specify the directory where you want to save the downloaded PDF file
		downloadDir := "/home/mohammedyameen/Music/drivedownloads"
		// Ensure the directory exists
		if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to create download directory: %v", err))
		}

		// Construct the download path
		downloadPath := filepath.Join(downloadDir, latestPDF.Name)

		err = downloadFile(srv, latestPDF.Id, downloadPath)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to download the latest PDF file: %v", err))
		}

		return c.String(http.StatusOK, fmt.Sprintf("Latest PDF file downloaded successfully to: %s", downloadPath))
	})

	e.Logger.Fatal(e.Start(":8080"))
}