package redis

import (
	"context"
	"fmt"
	"managedata/model"
	"managedata/utils"
	"time"
)

func SetEmployee(employees []model.Employee) error {
	logs := utils.GetLogger()
	ctx := context.Background()

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

func SetEmployeeByIdInRedis(ctx context.Context, employee model.Employee) error {
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
