package mysql

import (
	"context"
	"fmt"
	"managedata/model"
)

// FetchAllEmployees fetches all employee records from MySQL.
func FetchAllEmployeesFromMysql() ([]model.Employee, error) {
	ctx := context.Background()

	const query = `
		SELECT id, first_name, last_name, company_name, address, city, county, postal, phone, email, web
		FROM Employee
	`

	var employees []model.Employee

	if err := Db.SelectContext(ctx, &employees, query); err != nil {
		return nil, fmt.Errorf("failed to fetch employees from database: %w", err)
	}

	return employees, nil
}


