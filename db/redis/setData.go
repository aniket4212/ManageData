package redis

import (
	"context"
	"fmt"
	"managedata/model"
	"managedata/utils"
	"time"
)

func SetEmployee(employees []model.Employee) error {
	ctx := context.Background()
	logs := utils.GetLogger()

	ttl := 300

	for _, employee := range employees {

		key := fmt.Sprintf("employee:%s", employee.ID)

		err := RDB.HSet(ctx, key, map[string]interface{}{
			"first_name":   employee.FirstName,
			"last_name":    employee.LastName,
			"company_name": employee.CompanyName,
			"address":      employee.Address,
			"city":         employee.City,
			"county":       employee.County,
			"postal":       employee.Postal,
			"phone":        employee.Phone,
			"email":        employee.Email,
			"web":          employee.Web,
		}).Err()
		if err != nil {
			logs.Error().Err(err).Msgf("Failed to store employee data with ID: %s", employee.ID)
			continue
		}
		err = RDB.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
		if err != nil {
			logs.Error().Err(err).Msgf("Failed to set TTL for employee data with key: %s", key)
			continue
		}

		// logs.Info().Msgf("Successfully stored employee data in Redis with key: %s", key)
	}

	return nil
}

func SetSingleEmployeeInRedis(employee model.Employee) error {
	ctx := context.Background()
	key := fmt.Sprintf("employee:%s", employee.ID)

	data := map[string]interface{}{
		"first_name":   employee.FirstName,
		"last_name":    employee.LastName,
		"company_name": employee.CompanyName,
		"address":      employee.Address,
		"city":         employee.City,
		"county":       employee.County,
		"postal":       employee.Postal,
		"phone":        employee.Phone,
		"email":        employee.Email,
		"web":          employee.Web,
	}

	return RDB.HSet(ctx, key, data).Err()
}

// func GetEmployeeFromRedis(employeeID string) (*model.Employee, error) {
// 	ctx := context.Background()
// 	logs := utils.GetLogger()

// 	key := fmt.Sprintf("employee:%s", employeeID)

// 	// Attempt to get employee data from Redis
// 	data, err := RDB.HGetAll(ctx, key).Result()
// 	if err != nil {
// 		logs.Error().Err(err).Msgf("Failed to retrieve employee data for ID: %s from Redis", employeeID)
// 		return nil, fmt.Errorf("failed to retrieve employee data from Redis: %v", err)
// 	}

// 	// If the data doesn't exist, return nil
// 	if len(data) == 0 {
// 		logs.Warn().Msgf("No data found for employee ID: %s in Redis", employeeID)
// 		return nil, nil
// 	}

// 	// Map Redis data to Employee struct
// 	employee := &model.Employee{
// 		ID:          employeeID,
// 		FirstName:   data["FirstName"],
// 		LastName:    data["LastName"],
// 		CompanyName: data["CompanyName"],
// 		Address:     data["Address"],
// 		City:        data["City"],
// 		County:      data["County"],
// 		Postal:      data["Postal"],
// 		Phone:       data["Phone"],
// 		Email:       data["Email"],
// 		Web:         data["Web"],
// 	}

// 	return employee, nil
// }

// func FetchImportedListFromRedis() ([]model.Employee, error) {
// 	ctx := context.Background()

// 	// Fetch data from Redis
// 	employees, err := RDB.HGetAll(ctx, "employee:list").Result()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get data from Redis: %v", err)
// 	}

// 	if len(employees) == 0 {
// 		return nil, nil // Indicating data is not found
// 	}

// 	// Convert the Redis response into the employee model (adapt based on how you store the data)
// 	var employeeList []model.Employee
// 	for key, value := range employees {
// 		emp := model.Employee{
// 			ID:        key,
// 			FirstName: value, // Example, adjust based on your data format
// 		}
// 		employeeList = append(employeeList, emp)
// 	}

// 	return employeeList, nil
// }
