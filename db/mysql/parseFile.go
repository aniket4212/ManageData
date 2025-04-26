package mysql

import (
	"fmt"
	"io"
	"managedata/db/redis"
	"managedata/model"
	"managedata/utils"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

var ExpectedHeaders = []string{
	"first_name", "last_name", "company_name", "address", "city",
	"county", "postal", "phone", "email", "web",
}

func CompareHeaders(headers []string, expected []string) bool {
	if len(headers) != len(expected) {
		return false
	}
	for i := range expected {
		if strings.ToLower(headers[i]) != strings.ToLower(expected[i]) {
			return false
		}
	}
	return true
}

func ParseAndInsertExcelFile(file io.Reader) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Printf("unable to parse Excel file: %v\n", err)
		return
	}
	defer f.Close()

	logs := utils.GetLogger()

	sheetNames := f.GetSheetList()

	var employees []model.Employee
	const batchSize = 100
	var wg sync.WaitGroup
	var totalInserted int64

	for _, sheetName := range sheetNames {
		rows, err := f.GetRows(sheetName)
		if err != nil || len(rows) == 0 {
			fmt.Printf("unable to read rows from sheet %s: %v", sheetName, err)
			return
		}

		for i, row := range rows {
			if i == 0 {
				continue
			}
			if len(row) < 10 {
				continue
			}

			eachRow := make([]string, 10)
			copy(eachRow, row)

			if !utils.IsValidEmail(eachRow[8]) {
				logs.Warn().Msgf("Skipping row %d: Invalid email '%s'", i, eachRow[8])
				continue
			}
			if !utils.IsValidPhone(eachRow[7]) {
				logs.Warn().Msgf("Skipping row %d: Invalid phone '%s'", i, eachRow[7])
				continue
			}

			emp := model.Employee{
				ID:          fmt.Sprintf("%s-%d", uuid.New().String(), time.Now().UnixMicro()),
				FirstName:   eachRow[0],
				LastName:    eachRow[1],
				CompanyName: eachRow[2],
				Address:     eachRow[3],
				City:        eachRow[4],
				County:      eachRow[5],
				Postal:      eachRow[6],
				Phone:       eachRow[7],
				Email:       eachRow[8],
				Web:         eachRow[9],
			}
			employees = append(employees, emp)

			if len(employees) == batchSize {
				wg.Add(1)
				batch := make([]model.Employee, len(employees))
				copy(batch, employees)
				go UpdateDatabaseAndCache(batch, &wg, &totalInserted)
				employees = employees[:0]
			}
		}
	}

	if len(employees) > 0 {
		wg.Add(1)
		batch := make([]model.Employee, len(employees))
		copy(batch, employees)
		go UpdateDatabaseAndCache(batch, &wg, &totalInserted)
	}

	wg.Wait()
	logs.Info().Msgf("Total records inserted into MySQL: %d", totalInserted)

}

func UpdateDatabaseAndCache(employees []model.Employee, wg *sync.WaitGroup, totalInserted *int64) {
	defer wg.Done()

	if len(employees) == 0 {
		return
	}

	InsertBatch(employees, totalInserted)

	redis.SetEmployee(employees)
}
