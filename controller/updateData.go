package controller

import (
	"managedata/db/mysql"
	"managedata/db/redis"
	"managedata/model"
	"managedata/utils" 
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateEmployee(c *gin.Context) {
	logs := utils.GetLogger()

	var employee model.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		logs.Error().Err(err).Msg("Failed to bind JSON for UpdateEmployee")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	logs.Info().Str("EmployeeID", employee.ID).Msg("Starting employee update process")

	// Fetch the current employee from DB
	existingEmployee, err := mysql.FetchEmployeeByID(employee.ID)
	if err != nil {
		logs.Error().Err(err).Str("EmployeeID", employee.ID).Msg("Failed to fetch employee from DB")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch employee",
			"error":   err.Error(),
		})
		return
	}

	if existingEmployee.ID == "" {
		logs.Warn().Str("EmployeeID", employee.ID).Msg("Employee not found in DB")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Employee not found",
		})
		return
	}

	// Update only non-empty fields
	if employee.FirstName != "" {
		existingEmployee.FirstName = employee.FirstName
	}
	if employee.LastName != "" {
		existingEmployee.LastName = employee.LastName
	}
	if employee.CompanyName != "" {
		existingEmployee.CompanyName = employee.CompanyName
	}
	if employee.Address != "" {
		existingEmployee.Address = employee.Address
	}
	if employee.City != "" {
		existingEmployee.City = employee.City
	}
	if employee.County != "" {
		existingEmployee.County = employee.County
	}
	if employee.Postal != "" {
		existingEmployee.Postal = employee.Postal
	}
	if employee.Phone != "" {
		existingEmployee.Phone = employee.Phone
	}
	if employee.Email != "" {
		existingEmployee.Email = employee.Email
	}
	if employee.Web != "" {
		existingEmployee.Web = employee.Web
	}

	// Update in Redis
	if err := redis.SetSingleEmployeeInRedis(existingEmployee); err != nil {
		logs.Error().Err(err).Str("EmployeeID", employee.ID).Msg("Failed to update employee in Redis")
		// (continue to update MySQL even if Redis fails)
	} else {
		logs.Info().Str("EmployeeID", employee.ID).Msg("Successfully updated employee in Redis")
	}

	// Update in MySQL
	if err := mysql.UpdateEmployeeByID(existingEmployee); err != nil {
		logs.Error().Err(err).Str("EmployeeID", employee.ID).Msg("Failed to update employee in MySQL")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update employee",
			"error":   err.Error(),
		})
		return
	}

	logs.Info().Str("EmployeeID", employee.ID).Msg("Successfully updated employee in MySQL and Redis")

	c.JSON(http.StatusOK, gin.H{
		"message": "Employee updated successfully",
	})
}
