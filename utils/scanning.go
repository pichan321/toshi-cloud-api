package utils

import (
	"github.com/jmoiron/sqlx"
)

func ScanToMap(columnNames []string, rows *sqlx.Row) map[string]string {
	data := make(map[string]string, len(columnNames))
	columns := make([]string, len(columnNames))
	columnPointers := make([]interface{}, len(columnNames))

	for i, _ := range columnNames {
		columnPointers[i] = &columns[i]
	}

	rows.Scan(columnPointers...)

	for i, v := range columnNames {
		data[v] = columns[i]
	}

	return data
}

func ScanToMapRows(columnNames []string, rows *sqlx.Rows) map[string]interface{} {
	data := make(map[string]interface{}, len(columnNames))
	columns := make([]string, len(columnNames))
	columnPointers := make([]interface{}, len(columnNames))

	for i, _ := range columnNames {
		columnPointers[i] = &columns[i]
	}

	rows.Scan(columnPointers...)

	for i, v := range columnNames {
		data[v] = columns[i]
	}

	return data
}
