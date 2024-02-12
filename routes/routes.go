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

	files, err := service.ListPDFFiles(srv)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to list PDF files: %v", err))
	}

	latestPDF, err := service.GetLatestPDFFile(files)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to get the latest PDF file: %v", err))
	}

	url := fmt.Sprintf("https://drive.google.com/file/d/%s/view", latestPDF.Id)
	return c.Redirect(http.StatusSeeOther, url)
}
