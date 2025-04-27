package controller

import (
	"bytes"
	"io"
	"managedata/services"
	"managedata/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func UploadExcelHandler(c *gin.Context) {
	logs := utils.GetLogger()
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logs.Info().Msgf("Failed to get file from request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logs.Info().Msgf("Unable to read uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read uploaded file"})
		return
	}
	// Create a new reader from the bytes
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		logs.Info().Msgf("Unable to parse Excel file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse Excel file"})
		return
	}
	defer f.Close()

	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		logs.Info().Msgf("No sheets found in the Excel file")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No sheets found in the Excel file"})
		return
	}

	for _, sheetName := range sheetNames {
		rows, err := f.GetRows(sheetName)
		if err != nil || len(rows) == 0 {
			logs.Info().Msgf("unable to read rows from sheet %s: %v", sheetName, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read rows from sheet"})
			return
		}

		header := rows[0]
		if !services.CompareHeaders(header, services.ExpectedHeaders) {
			logs.Info().Msgf("header format is incorrect in sheet %s. Expected: %v, got: %v", sheetName, services.ExpectedHeaders, header)
			c.JSON(http.StatusBadRequest, gin.H{"error": "header format is incorrect in sheet"})
			return
		}

	}
	go func() {
		logs.Info().Msgf("Starting async parsing and DB insertion")
		services.ParseAndInsertExcelFile(bytes.NewReader(fileBytes))
	}()

	logs.Info().Msgf("File uploaded successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Excel data inserted into DB"})
}
