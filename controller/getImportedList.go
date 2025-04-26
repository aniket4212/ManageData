package controller

import (
	"fmt"
	"managedata/db/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ViewImportedList(c *gin.Context) {
	// Fetch from MySQL if Redis was empty
	employees, err := mysql.FetchAllEmployeesFromMysql()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to retrieve data from MySQL: %v", err),
		})
		return
	}

	if len(employees) > 0 {
		log.Info().Msgf("Successfully fetched employee data")
		c.JSON(http.StatusOK, gin.H{
			"count":     len(employees),
			"employees": employees,
		})
	} else {
		log.Info().Msg("No data found in MySQL")
		c.JSON(http.StatusOK, gin.H{
			"message": "No data found in Redis or MySQL",
		})
	}
}

// func ViewEmployeeByID(c *gin.Context) {
// 	var req model.GetEmployeeRequest

// 	if err := c.ShouldBindJSON(req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid request payload",
// 		})
// 		return
// 	}

// 	// First, check Redis for the specific employee
// 	employee, err := redis.GetEmployeeByIDFromRedis(req.EmployeeID)
// 	if err == nil {
// 		log.Info().Msgf("✅ Successfully fetched employee ID %s from Redis", req.EmployeeID)
// 		c.JSON(http.StatusOK, gin.H{
// 			"employee": employee,
// 		})
// 		return
// 	}

// 	log.Info().Msgf("🔄 Employee ID %s not found in Redis, fetching from MySQL", req.EmployeeID)

// 	// Fetch from MySQL if not found in Redis
// 	employee, err = mysql.FetchEmployeeByID(req.EmployeeID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": fmt.Sprintf("Failed to retrieve employee from MySQL: %v", err),
// 		})
// 		return
// 	}

// 	if employee.ID == "" {
// 		log.Info().Msgf("❌ Employee ID %s not found in MySQL", req.EmployeeID)
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "Employee not found",
// 		})
// 		return
// 	}

// 	// Cache the employee in Redis
// 	err = redis.SetSingleEmployeeInRedis(employee)
// 	if err != nil {
// 		log.Warn().Msgf("⚠️ Failed to cache employee ID %s into Redis: %v", req.EmployeeID, err)
// 	} else {
// 		log.Info().Msgf("✅ Successfully cached employee ID %s into Redis", req.EmployeeID)
// 	}

// 	log.Info().Msgf("✅ Successfully fetched employee ID %s from MySQL", req.EmployeeID)
// 	c.JSON(http.StatusOK, gin.H{
// 		"employee": employee,
// 	})
// }
