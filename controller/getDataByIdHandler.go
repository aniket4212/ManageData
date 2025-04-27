package controller

import (
	"fmt"
	"managedata/db/mysql"
	"managedata/db/redis"
	"managedata/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDataByIdHandler(c *gin.Context) {

	ctx := c.Request.Context()

	employeeID := c.DefaultQuery("id", "")

	logs := utils.GetLogger()

	if employeeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Employee ID is required",
			"error":   "BadRequest",
		})
		return
	}

	// First, check Redis for the specific employee
	employee, err := redis.GetEmployeeByIDFromRedis(ctx, employeeID)
	if err == nil {
		logs.Info().Msgf("Successfully fetched employee ID %s from Redis", employeeID)
		c.JSON(http.StatusOK, gin.H{
			"employee": employee,
		})
		return
	}

	logs.Warn().Err(err).Msgf("Failed to fetch employee ID %s from Redis, trying MySQL...", employeeID)

	// Fetch from MySQL if not found in Redis
	employee, err = mysql.FetchEmployeeByID(ctx, employeeID)
	if err != nil {
		logs.Error().Err(err).Msgf("Failed to retrieve employee ID %s from MySQL", employeeID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to retrieve employee from MySQL: %v", err),
		})
		return
	}

	if employee.ID == "" {
		logs.Info().Msgf("Employee ID %s not found in MySQL", employeeID)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Employee not found",
		})
		return
	}

	// Cache the employee in Redis
	err = redis.SetEmployeeByIdInRedis(ctx, employee)
	if err != nil {
		logs.Warn().Err(err).Msgf("Failed to cache employee ID %s into Redis", employeeID)
	} else {
		logs.Info().Msgf("Successfully cached employee ID %s into Redis", employeeID)
	}

	logs.Info().Msgf("Successfully fetched employee ID %s from MySQL", employeeID)
	c.JSON(http.StatusOK, gin.H{
		"employee": employee,
	})
}
