package mysql

import (
	"context"
	"fmt"
	"managedata/model"
)

func FetchAllEmployeesFromMysql(ctx context.Context) ([]model.Employee, error) {

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
