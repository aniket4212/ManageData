package mysql

import (
	"context"
	"fmt"
	"managedata/model"
)

func UpdateEmployeeByID(employee model.Employee) error {
	ctx := context.Background()

	query := `
		UPDATE Employee
		SET first_name = ?, last_name = ?, company_name = ?, address = ?, city = ?, county = ?, postal = ?, phone = ?, email = ?, web = ?
		WHERE id = ?
	`

	_, err = Db.ExecContext(ctx, query, employee.FirstName, employee.LastName, employee.CompanyName, employee.Address,
		employee.City, employee.County, employee.Postal, employee.Phone, employee.Email, employee.Web, employee.ID)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}
