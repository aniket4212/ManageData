package redis

import (
	"context"
	"fmt"
	"managedata/model"

	"github.com/redis/go-redis/v9"
)

func GetEmployeeByIDFromRedis(ctx context.Context, employeeID string) (model.Employee, error) {

	key := fmt.Sprintf("employee:%s", employeeID)

	result, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return model.Employee{}, fmt.Errorf("employee with ID %s not found in Redis", employeeID)
		}
		return model.Employee{}, fmt.Errorf("failed to fetch employee data from Redis: %w", err)
	}

	if len(result) == 0 {
		return model.Employee{}, fmt.Errorf("employee with ID %s not found in Redis", employeeID)
	}

	employee := model.Employee{
		ID:          employeeID,
		FirstName:   result["first_name"],
		LastName:    result["last_name"],
		CompanyName: result["company_name"],
		Address:     result["address"],
		City:        result["city"],
		County:      result["county"],
		Postal:      result["postal"],
		Phone:       result["phone"],
		Email:       result["email"],
		Web:         result["web"],
	}

	return employee, nil
}

func GetAllEmployeesFromRedis(ctx context.Context) ([]model.Employee, error) {

	// Get all keys matching pattern
	keys, err := RDB.Keys(ctx, "employee:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %v", err)
	}

	var employees []model.Employee

	for _, key := range keys {
		data, err := RDB.HGetAll(ctx, key).Result()
		if err != nil || len(data) == 0 {
			continue
		}

		emp := model.Employee{
			ID:          key[len("employee:"):],
			FirstName:   data["FirstName"],
			LastName:    data["LastName"],
			CompanyName: data["CompanyName"],
			Address:     data["Address"],
			City:        data["City"],
			County:      data["County"],
			Postal:      data["Postal"],
			Phone:       data["Phone"],
			Email:       data["Email"],
			Web:         data["Web"],
		}
		employees = append(employees, emp)
	}

	return employees, nil
}
