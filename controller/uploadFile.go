package controller

import (
	"bytes"
	"fmt"
	"io"
	"managedata/db/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func UploadExcelHandler(c *gin.Context) {
	// Parse the file from the multipart request
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read uploaded file"})
		return
	}
	// Create a new reader from the bytes
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse Excel file"})
		return
	}
	defer f.Close()

	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No sheets found in the Excel file"})
		return
	}

	for _, sheetName := range sheetNames {
		rows, err := f.GetRows(sheetName)
		if err != nil || len(rows) == 0 {
			fmt.Printf("unable to read rows from sheet %s: %v", sheetName, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read rows from sheet"})
			return
		}

		header := rows[0]
		if !mysql.CompareHeaders(header, mysql.ExpectedHeaders) {
			fmt.Printf("header format is incorrect in sheet %s. Expected: %v, got: %v", sheetName, mysql.ExpectedHeaders, header)
			c.JSON(http.StatusBadRequest, gin.H{"error": "header format is incorrect in sheet"})
			return
		}

	}
	// Parse and insert the Excel data
	go mysql.ParseAndInsertExcelFile(bytes.NewReader(fileBytes))

	// If everything went fine, return success message
	c.JSON(http.StatusOK, gin.H{"message": "Excel data inserted into DB"})
	fmt.Println("File uploaded successfully")
}
