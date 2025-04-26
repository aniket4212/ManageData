package mysql

import (
	"context"
	"strings"
	"sync/atomic"

	"managedata/model"
	"managedata/utils"
)

func InsertBatch(employees []model.Employee, totalInserted *int64) {

	logs := utils.GetLogger()
	ctx := context.Background()

	if len(employees) == 0 {
		logs.Warn().Msg(" No employees to insert in this batch")
		return
	}

	// Begin transaction
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		logs.Error().Err(err).Msg("Failed to begin transaction")
		return
	}

	// Construct bulk INSERT query
	var builder strings.Builder
	builder.WriteString(`
		INSERT INTO Employee (
			id, first_name, last_name, company_name, address, city, county, postal, phone, email, web
		) VALUES 
	`)

	values := []interface{}{}
	for i, emp := range employees {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		values = append(values,
			emp.ID, emp.FirstName, emp.LastName, emp.CompanyName,
			emp.Address, emp.City, emp.County, emp.Postal,
			emp.Phone, emp.Email, emp.Web,
		)
	}

	stmt := builder.String()
	result, err := tx.ExecContext(ctx, stmt, values...)
	if err != nil {
		tx.Rollback()
		logs.Error().Err(err).Msg("Failed to execute batch insert")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if err := tx.Commit(); err != nil {
		logs.Error().Err(err).Msg("Failed to commit transaction")
		return
	}
	atomic.AddInt64(totalInserted, rowsAffected)

	logs.Info().Msgf("Inserted %d employees into MySQL", rowsAffected)
}
