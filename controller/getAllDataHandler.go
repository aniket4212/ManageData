package controller

import (
	"fmt"
	"managedata/db/mysql"
	"managedata/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllDataHandler(c *gin.Context) {
	ctx := c.Request.Context() 
	logs := utils.GetLogger()
	
	employees, err := mysql.FetchAllEmployeesFromMysql(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to retrieve data from MySQL: %v", err),
		})
		return
	}

	if len(employees) > 0 {
		logs.Info().Msgf("Successfully fetched employee data")
		c.JSON(http.StatusOK, gin.H{
			"count":     len(employees),
			"employees": employees,
		})
	} else {
		logs.Info().Msg("No data found in MySQL")
		c.JSON(http.StatusOK, gin.H{
			"message": "No data found in Redis or MySQL",
		})
	}
}
