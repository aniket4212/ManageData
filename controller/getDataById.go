package controller

import (
	"fmt"
	"managedata/db/mysql"
	"managedata/db/redis"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ViewEmployeeByID(c *gin.Context) {
	employeeID := c.DefaultQuery("id", "")

	if employeeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Employee ID is required",
			"error":   "BadRequest",
		})
		return
	}

	// First, check Redis for the specific employee
	employee, err := redis.GetEmployeeByIDFromRedis(employeeID)
	if err == nil {
		log.Info().Msgf("Successfully fetched employee ID %s from Redis", employeeID)
		c.JSON(http.StatusOK, gin.H{
			"employee": employee,
		})
		return
	}

	log.Warn().Err(err).Msgf("Failed to fetch employee ID %s from Redis, trying MySQL...", employeeID)

	// Fetch from MySQL if not found in Redis
	employee, err = mysql.FetchEmployeeByID(employeeID)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to retrieve employee ID %s from MySQL", employeeID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to retrieve employee from MySQL: %v", err),
		})
		return
	}

	if employee.ID == "" {
		log.Info().Msgf("Employee ID %s not found in MySQL", employeeID)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Employee not found",
		})
		return
	}

	// Cache the employee in Redis
	err = redis.SetSingleEmployeeInRedis(employee)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to cache employee ID %s into Redis", employeeID)
	} else {
		log.Info().Msgf("Successfully cached employee ID %s into Redis", employeeID)
	}

	log.Info().Msgf("Successfully fetched employee ID %s from MySQL", employeeID)
	c.JSON(http.StatusOK, gin.H{
		"employee": employee,
	})
}
