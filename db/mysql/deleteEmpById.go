package mysql

import (
	"context"
	"fmt"
)

func DeleteEmployeeByID(ctx context.Context, id string) error {

	query := `DELETE FROM Employee WHERE id = ?`

	_, err = Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}
