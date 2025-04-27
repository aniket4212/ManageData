package redis

import (
	"context"
	"fmt"
)

func DeleteEmployeeFromRedis(ctx context.Context, employeeID string) error {

	Key := fmt.Sprintf("employee:%s", employeeID)

	_, err := RDB.Del(ctx, Key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete employee from Redis: %w", err)
	}

	return nil
}
