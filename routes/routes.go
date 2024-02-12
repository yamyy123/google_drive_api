package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"googledrive/service"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/open-pdf", OpenPDF)
}

func OpenPDF(c echo.Context) error {
	log.Println("google drive service called")
		srv, err := service.GetDriveService()
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to retrieve Drive service: %v", err))
		}

		// Replace "YOUR_FOLDER_ID" with the actual ID of the folder you want to search in
		folderID := "1kM5D_O2m4X99utS24K4TlYmXzNm9Mopf"

		files, err := service.ListPDFFiles(srv, folderID)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to list PDF files: %v", err))
		}

		latestPDF, err := service.GetLatestPDFFile(files)
		if err != nil {
			if err.Error() == "no PDF files found" {
				return c.String(http.StatusNotFound, "No PDF files found in the specified folder")
			}
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to get the latest PDF file: %v", err))
		}

		url := fmt.Sprintf("https://drive.google.com/file/d/%s/view", latestPDF.Id)
		return c.Redirect(http.StatusSeeOther, url)
	}

