package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"managedata/model"
)

func FetchEmployeeByID(id string) (model.Employee, error) {
	ctx := context.Background()

	query := `
		SELECT id, first_name, last_name, company_name, address, city, county, postal, phone, email, web
		FROM Employee
		WHERE id = ?
	`

	var employee model.Employee
	err = Db.GetContext(ctx, &employee, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Employee{}, nil
		}
		return model.Employee{}, fmt.Errorf("failed to fetch data from Employee table: %w", err)
	}

	return employee, nil
}
