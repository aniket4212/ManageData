package controller

import (
	"managedata/db/mysql"
	"managedata/db/redis"
	"managedata/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func DeleteEmployeeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	logs := utils.GetLogger()

	employeeID := c.Query("id")

	if employeeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Employee ID is required",
			"error":   "BadRequest",
		})
		return
	}

	err := redis.DeleteEmployeeFromRedis(ctx, employeeID)

	if err != nil {
		logs.Error().Err(err).Str("employee_id", employeeID).Msg("Failed to delete employee from Redis")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete employee ID from redis",
		})
		return
	}

	err = mysql.DeleteEmployeeByID(ctx, employeeID)
	if err != nil {
		log.Info().Msgf("Failed Deleted employee ID %s", employeeID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete employee ID",
		})
		return
	}

	logs.Info().Msgf("Successfully Deleted employee ID %s", employeeID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully delete employee",
	})

}
