package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sosedoff/pgweb/pkg/statements"
)

var (
	errNoPrimaryKey = errors.New("table has no primary key")
	errNoValues     = errors.New("no values provided")
)

// TablePrimaryKeys returns the list of primary key column names for a table.
func (client *Client) TablePrimaryKeys(table string) ([]string, error) {
	schema, tbl := getSchemaAndTable(table)

	res, err := client.query(statements.TablePrimaryKeys, schema, tbl)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, row := range res.Rows {
		if len(row) > 0 && row[0] != nil {
			keys = append(keys, fmt.Sprintf("%v", row[0]))
		}
	}

	return keys, nil
}

// InsertRow inserts a new row with the provided column values. Bound parameters
// are used so PostgreSQL infers and validates types per the target column.
func (client *Client) InsertRow(table string, values map[string]interface{}) (*Result, error) {
	schema, tbl := getSchemaAndTable(table)

	cols := []string{}
	placeholders := []string{}
	args := []interface{}{}

	i := 1
	for col, val := range values {
		cols = append(cols, fmt.Sprintf(`"%s"`, col))
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		args = append(args, val)
		i++
	}

	if len(cols) == 0 {
		sql := fmt.Sprintf(`INSERT INTO "%s"."%s" DEFAULT VALUES`, schema, tbl)
		return client.query(sql)
	}

	sql := fmt.Sprintf(
		`INSERT INTO "%s"."%s" (%s) VALUES (%s)`,
		schema, tbl,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	return client.query(sql, args...)
}

// UpdateRow updates a single row identified by its primary key values.
func (client *Client) UpdateRow(table string, primaryKey map[string]interface{}, values map[string]interface{}) (*Result, error) {
	schema, tbl := getSchemaAndTable(table)

	if len(primaryKey) == 0 {
		return nil, errNoPrimaryKey
	}
	if len(values) == 0 {
		return nil, errNoValues
	}

	set := []string{}
	args := []interface{}{}

	i := 1
	for col, val := range values {
		set = append(set, fmt.Sprintf(`"%s" = $%d`, col, i))
		args = append(args, val)
		i++
	}

	where := []string{}
	for col, val := range primaryKey {
		where = append(where, fmt.Sprintf(`"%s" = $%d`, col, i))
		args = append(args, val)
		i++
	}

	sql := fmt.Sprintf(
		`UPDATE "%s"."%s" SET %s WHERE %s`,
		schema, tbl,
		strings.Join(set, ", "),
		strings.Join(where, " AND "),
	)

	return client.query(sql, args...)
}

// DeleteRow deletes a single row identified by its primary key values.
func (client *Client) DeleteRow(table string, primaryKey map[string]interface{}) (*Result, error) {
	schema, tbl := getSchemaAndTable(table)

	if len(primaryKey) == 0 {
		return nil, errNoPrimaryKey
	}

	where := []string{}
	args := []interface{}{}

	i := 1
	for col, val := range primaryKey {
		where = append(where, fmt.Sprintf(`"%s" = $%d`, col, i))
		args = append(args, val)
		i++
	}

	sql := fmt.Sprintf(
		`DELETE FROM "%s"."%s" WHERE %s`,
		schema, tbl,
		strings.Join(where, " AND "),
	)

	return client.query(sql, args...)
}
