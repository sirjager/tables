package core_repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	TABLE_NOT_FOUND = "table not found"
)

type Row struct {
	Value map[string]interface{} `json:"row"`
}

// This is where magic happens
func sqlRowsToJson(r *sql.Rows, k []string) ([]any, error) {
	var rs []any = []any{}
	for r.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		c := make([]interface{}, len(k))
		ptr := make([]interface{}, len(k))
		for i := range c {
			ptr[i] = &c[i]
		}
		// Scan the result into the column pointers...
		if err := r.Scan(ptr...); err != nil {
			return nil, err
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, cn := range k {
			val := ptr[i].(*interface{})
			m[cn] = *val
		}
		f := make(map[string]interface{})
		for k, v := range m {
			f[k] = v
		}
		rs = append(rs, f)
	}
	return rs, nil
}

func (s *TableSchema) insertRowString(row Row, index int) (string, error) {
	columns := ""
	values := ""
	i := 0
	err := s.ValidateRequiredColumns([]map[string]interface{}{row.Value})
	if err != nil {
		return "", err
	}
	for k, v := range row.Value {
		isLast := i == len(row.Value)-1
		for _, c := range s.Columns {
			if k == c.Name {
				if isLast {
					columns += c.Name
				} else {
					columns += c.Name + ","
				}
				if c.Type == "varchar" || c.Type == "text" {
					if c.Type == "varchar" {
						if c.Length > 0 {
							if v != nil && len(fmt.Sprintf("%v", v)) > int(c.Length) {
								return "", fmt.Errorf("row [%d] column [%v] can not have characters more than [%d]",
									index+1, k, c.Length)
							}
						}
					}
					if isLast {
						values += fmt.Sprintf("'%v'", v)
					} else {
						values += fmt.Sprintf("'%v',", v)
					}
				} else if c.Type == "boolean" {
					bv, isb := v.(bool)
					if !isb { // if it is not boolean then we will check whether if it is valid in string or not
						if strings.ToLower(fmt.Sprintf("%v", v)) == "true" || strings.ToLower(fmt.Sprintf("%v", v)) == "false" {
							if isLast {
								values += fmt.Sprintf("%v", v)
							} else {
								values += fmt.Sprintf("%v,", v)
							}
						} else {
							return "", fmt.Errorf("invalid value [%v] for boolean type column [%s]", v, c.Name)
						}
					} else {
						if isLast {
							values += fmt.Sprintf("%v", bv)
						} else {
							values += fmt.Sprintf("%v,", bv)
						}
					}
				} else {
					//! For any other type of column we will just add value
					if isLast {
						values += fmt.Sprintf("%v", v)
					} else {
						values += fmt.Sprintf("%v,", v)
					}
				}
				break
			}
		}
		i++
	}
	if columns == "" || values == "" {
		return "", fmt.Errorf("empty row provided")
	}
	return fmt.Sprintf(`INSERT INTO "public"."%v" (%v) VALUES (%v);`, s.Name, columns, values), err
}

func (s *TableSchema) updateRowString(row Row, primaryColumn Column, index int) (string, error) {
	var err error
	//END goal to make like this:
	// title = 'First Title', author = 'First Author' WHERE id = 1;
	query := ""
	// We dont want to update primary column value so we will remove it

	rowIndex := row.Value[primaryColumn.Name]
	if rowIndex == nil {
		return query, fmt.Errorf("which row to update ?. provide a primary column and  value in row [%d]", index+1)
	}
	println(fmt.Sprintf("Update row #%v", rowIndex))
	i := 0
	for k, v := range row.Value {
		isLast := i == len(row.Value)-1
		for _, c := range s.Columns {
			if k == c.Name && k != primaryColumn.Name {
				query += c.Name + " = "
				if c.Type == "varchar" || c.Type == "text" {
					if c.Type == "varchar" {
						if c.Length > 0 {
							if v != nil && len(fmt.Sprintf("%v", v)) > int(c.Length) {
								return "", fmt.Errorf("row [%d] column [%v] can not have characters more than [%d]",
									index+1, k, c.Length)
							}
						}
					}
					if isLast {
						query += fmt.Sprintf("'%v'", v)
					} else {
						query += fmt.Sprintf("'%v',", v)
					}
				} else if c.Type == "boolean" {
					bv, isb := v.(bool)
					if !isb { // if it is not boolean then we will check whether if it is valid in string or not
						if strings.ToLower(fmt.Sprintf("%v", v)) == "true" || strings.ToLower(fmt.Sprintf("%v", v)) == "false" {
							if isLast {
								query += fmt.Sprintf("%v", v)
							} else {
								query += fmt.Sprintf("%v,", v)
							}
						} else {
							return "", fmt.Errorf("invalid value [%v] for boolean type column [%s]", v, c.Name)
						}
					} else {
						if isLast {
							query += fmt.Sprintf("%v", bv)
						} else {
							query += fmt.Sprintf("%v,", bv)
						}
					}
				} else {
					//! For any other type of column we will just add value
					if isLast {
						query += fmt.Sprintf("%v", v)
					} else {
						query += fmt.Sprintf("%v,", v)
					}
				}
				break
			}
		}
		i++
	}
	if query == "" {
		return "", fmt.Errorf("empty row provided")
	}
	if primaryColumn.Type == "text" || primaryColumn.Type == "varchar" {
		query += fmt.Sprintf(" WHERE %v = '%v'", primaryColumn.Name, rowIndex)
	} else {
		query += fmt.Sprintf(" WHERE %v = %v", primaryColumn.Name, rowIndex)
	}
	return fmt.Sprintf(`UPDATE "public"."%v" SET %v;`, s.Name, query), err
}

type InsertRowsParams struct {
	Table  string                   `json:"table" validate:"required,alphanum,gte=3,lte=60"`
	UserID int64                    `json:"user_id" validate:"required,numeric,min=1"`
	Rows   []map[string]interface{} `json:"rows" validate:"required"`
}

func (store *SQLStore) InsertRows(ctx context.Context, arg InsertRowsParams) error {
	validate := validator.New()
	err := validate.Struct(arg)
	if err != nil {
		return err
	}

	// First we will validate so that we dont make interactions in database with invalid data
	// Requirements: Table Schema
	dbtable, err := store.GetTableByUserIdAndTableName(ctx, GetTableByUserIdAndTableNameParams{Name: arg.Table, UserID: arg.UserID})
	// If schema is not found then probably table does not exits or does not belongs to user
	// In both cases we will send table not found.
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(TABLE_NOT_FOUND)
		}
		return err
	}
	schema, err := dbtable.Schema()
	if err != nil {
		return err
	}
	colsThatDontExist := schema.ColumnsThatDontExists(arg.Rows)
	if len(colsThatDontExist) != 0 {
		return fmt.Errorf("columns %v dosen't exists", colsThatDontExist)
	}
	var insertStrings []string
	for i, r := range arg.Rows {
		row := Row{Value: r}
		istr, err := schema.insertRowString(row, i)
		if err != nil {
			return err
		}
		insertStrings = append(insertStrings, istr)
	}
	// to use and read colums properly we need to format table columns
	insertStatement := strings.Join(insertStrings, "\n")
	// If insert strings are safely build without erros then execute statements
	_, err = store.db.ExecContext(ctx, insertStatement)
	return err
}

type GetRowsParams struct {
	UserID  int64                  `json:"user_id" validate:"required,numeric,min=1"`
	Table   string                 `json:"table" validate:"required,alphanum,min=1"`
	Fields  []string               `json:"fields" validate:""`
	Filters map[string]interface{} `json:"filters" validate:""`
}

func (q *Queries) GetRows(ctx context.Context, arg GetRowsParams) ([]any, error) {
	var err error
	validate := validator.New()
	err = validate.Struct(arg)
	if err != nil {
		return nil, err
	}
	var all_columns []string
	orRows := make(map[string]interface{})
	var andRows map[string]interface{}

	for col, vals := range arg.Filters {
		if col != "&" {
			orRows[col] = vals
			all_columns = append(all_columns, col)
		} else {
			andRowMap, isMap := vals.(map[string]interface{})
			if !isMap {
				return nil, fmt.Errorf("values inside & should be a map")
			}
			andRows = andRowMap
			for ncol := range andRowMap {
				all_columns = append(all_columns, ncol)
			}
		}
	}

	table, err := q.GetTableByUserIdAndTableName(ctx, GetTableByUserIdAndTableNameParams{Name: arg.Table, UserID: int64(arg.UserID)})

	if err != nil {
		return nil, err
	}

	mytable, err := FormatTableEntryToTable(table)
	if err != nil {
		return nil, err
	}

	//Check if whatever fields that is request exits else return error
	// If no fields are requested then send all back

	fieldsString, err := GenerateFieldString(arg.Fields, mytable.Columns)
	if err != nil {
		return nil, err
	}

	// Now we First Check if any of the column doesnt exists
	var columnThatDontExists []string = []string{}
	for _, c := range all_columns {
		exists := false
		for _, col := range mytable.Columns {
			if c == col.Name {
				exists = true
				break
			}
		}
		if !exists {
			columnThatDontExists = append(columnThatDontExists, c)
		}
	}
	if len(columnThatDontExists) != 0 {
		if len(columnThatDontExists) == 1 {
			return nil, fmt.Errorf("column %v does not exits", columnThatDontExists)
		}
		return nil, fmt.Errorf("columns %v does not exits", columnThatDontExists)
	}

	// Example of Main string:
	// if and -> SELECT * FROM tablename WHERE name IN ('John','Elsa') AND verified = true
	// if or ->  SELECT * FROM tablename WHERE email IN ('user1@email.com','user2@email.com') OR username IN ('user1','user2') OR name IN ('user one','user two')
	// if or+and ->  SELECT * FROM tablename WHERE email IN ('user1@email.com','user2@email.com') OR username IN ('user1','user2') AND verified = false

	// First we loop over orColumns
	// ori = key name(column name) , orv = value of item at ori

	var all_or_strings []string = []string{}

	for ork, orv := range orRows {

		// Goal is to make string like for non boolean columns:  name IN ('John','Elsa') , id IN (1,34,123)
		// Goal is to make string like for boolean columns :  verified = false

		// We also need to get data type of column
		dtype := ""
		// for datatype we need to loop over mytable.columns
		for _, mycol := range mytable.Columns {
			// if it is a same column
			if ork == mycol.Name {
				dtype = mycol.Type
				break
			}
		}

		// Now we need to add values and build column string
		// for that we need to loop over orv
		// orv is a interface for now we need a list and it must be list

		orvList, isList := orv.([]any)
		if !isList {
			return nil, fmt.Errorf("values of %v must be an arrary", ork)
		}

		// At this point our columnString = name
		// In general columnString = columnname

		// For any boolean column our row can only match true or false
		// so a single row can either true or false

		// we will make sure our list contains only one item in list of any boolean column
		if dtype == "boolean" {
			if len(orvList) == 0 {
				return nil, fmt.Errorf("value for boolean column can not be empty")
			}
			if len(orvList) != 1 {
				return nil, fmt.Errorf("%v is a boolean column and must have one boolean value", ork)
			}
			// if we have boolean column suppose: verified. then we want verified = true
			// not just anything like : verified = yaga yaga or antthing
			// we need to make sure that v is a boolean not just anyting

			booleanValue, isBoolean := orvList[0].(bool)
			if !isBoolean {
				return nil, fmt.Errorf("invalid value=(%v).  %v is a boolean column and must have one boolean value", orvList[0], ork)
			}

			all_or_strings = append(all_or_strings, fmt.Sprintf("%v = %v", ork, booleanValue))

			// now for any boolean column we will have
			// columnname = false
			// example verified = true
		} else if dtype == "text" || dtype == "varchar" {
			// if column is of text/varchar data type we need to wrap value in single quote
			// sql string values needs to be wrapped in a single quote if it is not a single word

			// we will also need to make sure that list is not empty

			if len(orvList) == 0 {
				return nil, fmt.Errorf("%v can not have an empty array", ork)
			}

			valueString := ""
			// now we will add values
			// for any text column: valueString = 'any string 1', 'anystring 2'
			// goal = valueString = 'John', 'elsa','any name'

			// we will loop over orvlist

			for i, v := range orvList {
				isLast := i == len(orvList)-1
				if isLast {
					if v != nil {
						valueString += fmt.Sprintf("'%v'", v)
					} else {
						valueString += "NULL"
					}
				} else {
					if v != nil {
						valueString += fmt.Sprintf("'%v',", v)
					} else {
						valueString += "NULL,"
					}
				}
			}
			all_or_strings = append(all_or_strings, fmt.Sprintf("%v IN (%v)", ork, valueString))

		} else {
			if len(orvList) == 0 {
				return nil, fmt.Errorf("%v can not have an empty arrary", ork)
			}

			valueString := ""

			// we will loop over orvlist
			for i, v := range orvList {
				isLast := i == len(orvList)-1
				if isLast {
					if v != nil {
						valueString += fmt.Sprintf("%v", v)
					} else {
						valueString += "NULL"
					}
				} else {
					if v != nil {
						valueString += fmt.Sprintf("%v,", v)
					} else {
						valueString += "NULL,"
					}
				}
			}
			all_or_strings = append(all_or_strings, fmt.Sprintf(" %v IN (%v)", ork, valueString))
		}

	}

	// Now we will loop over rows inside andRows
	// nk = key(column name)  nval = value of nk

	var all_and_strings []string = []string{}

	for nk, nval := range andRows {
		// same as or  strings

		dtype := ""

		for _, mycol := range mytable.Columns {
			if nk == mycol.Name {
				dtype = mycol.Type
				break
			}
		}

		nvList, isList := nval.([]any)
		if !isList {
			return nil, fmt.Errorf("values of %v must be an arrary", nk)
		}

		if dtype == "boolean" {
			if len(nvList) == 0 {
				return nil, fmt.Errorf("value for boolean column can not be empty")
			}
			if len(nvList) != 1 {
				return nil, fmt.Errorf("%v is a boolean column and must have one boolean value", nk)
			}
			booleanValue, isBoolean := nvList[0].(bool)
			if !isBoolean {
				return nil, fmt.Errorf("invalid value=(%v).  %v is a boolean column and must have one boolean value", nvList[0], nk)
			}

			all_and_strings = append(all_and_strings, fmt.Sprintf("%v = %v", nk, booleanValue))
		} else if dtype == "text" || dtype == "varchar" {
			if len(nvList) == 0 {
				return nil, fmt.Errorf("%v can not have an empty arrary", nk)
			}

			valueString := ""
			for i, v := range nvList {
				isLast := i == len(nvList)-1
				if isLast {
					if v != nil {
						valueString += fmt.Sprintf("'%v'", v)
					} else {
						valueString += "NULL"
					}
				} else {
					if v != nil {
						valueString += fmt.Sprintf("'%v',", v)
					} else {
						valueString += "NULL,"
					}
				}
			}
			all_and_strings = append(all_and_strings, fmt.Sprintf("%v IN (%v)", nk, valueString))
		} else {
			if len(nvList) == 0 {
				return nil, fmt.Errorf("%v can not have an empty arrary", nk)
			}

			valueString := ""

			// we will loop over orvlist
			for i, v := range nvList {
				isLast := i == len(nvList)-1
				if isLast {
					if v != nil {
						valueString += fmt.Sprintf("%v", v)
					} else {
						valueString += "NULL"
					}
				} else {
					if v != nil {
						valueString += fmt.Sprintf("%v,", v)
					} else {
						valueString += "NULL,"
					}
				}
			}
			all_and_strings = append(all_and_strings, fmt.Sprintf(" %v IN (%v)", nk, valueString))
		}
	}

	mainString := ""
	var all_string []string = []string{}

	for _, str := range all_or_strings {
		if len(all_string) != 0 {
			mainString += fmt.Sprintf("OR %v", str)
			all_string = append(all_string, str)
		} else {
			mainString += fmt.Sprintf("%v", str)
			all_string = append(all_string, str)
		}
	}
	for _, str := range all_and_strings {
		if len(all_string) != 0 {
			mainString += fmt.Sprintf(" AND %v", str)
			all_string = append(all_string, str)
		} else {
			mainString += fmt.Sprintf(" %v", str)
			all_string = append(all_string, str)
		}
	}

	query := fmt.Sprintf(`SELECT %v FROM  "public"."%v"`, fieldsString, arg.Table)
	if len(arg.Filters) != 0 {
		query += fmt.Sprintf("WHERE %v", mainString)
	}
	query += ";"
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		println(query)
		println(err.Error())
		return nil, err
	}
	cols, _ := rows.Columns()
	result, err := sqlRowsToJson(rows, cols)
	return result, err
}

type DeleteRowsParams struct {
	UserID  int64                  `json:"user_id" validate:"required,numeric,min=1"`
	Table   string                 `json:"table" validate:"required,alphanum,min=1"`
	Filters map[string]interface{} `json:"filters" validate:""`
}

func (q *Queries) DeleteRows(ctx context.Context, arg DeleteRowsParams) ([]any, error) {
	var err error
	validate := validator.New()
	err = validate.Struct(arg)
	if err != nil {
		return nil, err
	}
	// Column Statements
	// allColStrings := ""
	var all_columns []string
	orRows := make(map[string]interface{})
	var andRows map[string]interface{}

	for col, vals := range arg.Filters {
		if col != "&" {
			orRows[col] = vals
			all_columns = append(all_columns, col)
		} else {
			andRowMap, isMap := vals.(map[string]interface{})
			if !isMap {
				return nil, fmt.Errorf("values inside & should be a map")
			}
			andRows = andRowMap
			for ncol := range andRowMap {
				all_columns = append(all_columns, ncol)
			}
		}
	}

	table, err := q.GetTableByUserIdAndTableName(ctx, GetTableByUserIdAndTableNameParams{Name: arg.Table, UserID: int64(arg.UserID)})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("table '%s' not found", arg.Table)
		}
		return nil, err
	}
	mytable, err := FormatTableEntryToTable(table)
	if err != nil {
		return nil, err
	}

	//Check if whatever fields that is request exits else return error
	// If no fields are requested then send all back
	// Now we First Check if any of the column doesnt exists
	var columnThatDontExists []string = []string{}
	for _, c := range all_columns {
		exists := false
		for _, col := range mytable.Columns {
			if c == col.Name {
				exists = true
				break
			}
		}
		if !exists {
			columnThatDontExists = append(columnThatDontExists, c)
		}
	}
	if len(columnThatDontExists) != 0 {
		if len(columnThatDontExists) == 1 {
			return nil, fmt.Errorf("column %v does not exits", columnThatDontExists)
		}
		return nil, fmt.Errorf("columns %v does not exits", columnThatDontExists)
	}

	// Example of Main string:
	// if and -> SELECT * FROM tablename WHERE name IN ('John','Elsa') AND verified = true
	// if or ->  SELECT * FROM tablename WHERE email IN ('user1@email.com','user2@email.com') OR username IN ('user1','user2') OR name IN ('user one','user two')
	// if or+and ->  SELECT * FROM tablename WHERE email IN ('user1@email.com','user2@email.com') OR username IN ('user1','user2') AND verified = false

	// First we loop over orColumns
	// ori = key name(column name) , orv = value of item at ori

	var all_or_strings []string = []string{}

	for ork, orv := range orRows {

		// Goal is to make string like for non boolean columns:  name IN ('John','Elsa') , id IN (1,34,123)
		// Goal is to make string like for boolean columns :  verified = false

		// We also need to get data type of column
		dtype := ""
		// for datatype we need to loop over mytable.columns
		for _, mycol := range mytable.Columns {
			// if it is a same column
			if ork == mycol.Name {
				dtype = mycol.Type
				break
			}
		}

		// Now we need to add values and build column string
		// for that we need to loop over orv
		// orv is a interface for now we need a list and it must be list

		orvList, isList := orv.([]any)
		if !isList {
			return nil, fmt.Errorf("values of %v must be an arrary", ork)
		}

		// At this point our columnString = name
		// In general columnString = columnname

		// For any boolean column our row can only match true or false
		// so a single row can either true or false

		// we will make sure our list contains only one item in list of any boolean column
		if dtype == "boolean" {
			if len(orvList) == 0 {
				return nil, fmt.Errorf("value for boolean column can not be empty")
			}
			if len(orvList) != 1 {
				return nil, fmt.Errorf("%v is a boolean column and must have one boolean value", ork)
			}
			// if we have boolean column suppose: verified. then we want verified = true
			// not just anything like : verified = yaga yaga or antthing
			// we need to make sure that v is a boolean not just anyting

			booleanValue, isBoolean := orvList[0].(bool)
			if !isBoolean {
				return nil, fmt.Errorf("invalid value=(%v).  %v is a boolean column and must have one boolean value", orvList[0], ork)
			}

			all_or_strings = append(all_or_strings, fmt.Sprintf("%v = %v", ork, booleanValue))

			// now for any boolean column we will have
			// columnname = false
			// example verified = true
		} else if dtype == "text" || dtype == "varchar" {
			// if column is of text/varchar data type we need to wrap value in single quote
			// sql string values needs to be wrapped in a single quote if it is not a single word

			// we will also need to make sure that list is not empty

			if len(orvList) == 0 {
				return nil, fmt.Errorf("%v can not have an empty array", ork)
			}

			valueString := ""
			// now we will add values
			// for any text column: valueString = 'any string 1', 'anystring 2'
			// goal = valueString = 'John', 'elsa','any name'

			// we will loop over orvlist

			for i, v := range orvList {
				isLast := i == len(orvList)-1
				if isLast {
					if v != nil {
						valueString += fmt.Sprintf("'%v'", v)
					} else {
						valueString += "NULL"
					}
				} else {
					if v != nil {
						valueString += fmt.Sprintf("'%v',", v)
					} else {
						valueString += "NULL,"
					}
				}
			}
			all_or_strings = append(all_or_strings, fmt.Sprintf("%v IN (%v)", ork, valueString))

		} else {
			if len(orvList) == 0 {
				return nil, fmt.Errorf("%v can not have an empty arrary", ork)
			}

			valueString := ""

			// we will loop over orvlist
			for i, v := range orvList {
				isLast := i == len(orvList)-1
				if isLast {
					if v != nil {
						valueString += fmt.Sprintf("%v", v)
					} else {
						valueString += "NULL"
					}
				} else {
					if v != nil {
						valueString += fmt.Sprintf("%v,", v)
					} else {
						valueString += "NULL,"
					}
				}
			}
			all_or_strings = append(all_or_strings, fmt.Sprintf(" %v IN (%v)", ork, valueString))
		}

	}

	// Now we will loop over rows inside andRows
	// nk = key(column name)  nval = value of nk
	var all_and_strings []string = []string{}

	for nk, nval := range andRows {
		// same as or  strings

		dtype := ""

		for _, mycol := range mytable.Columns {
			if nk == mycol.Name {
				dtype = mycol.Type
				break
			}
		}

		nvList, isList := nval.([]any)
		if !isList {
			return nil, fmt.Errorf("values of %v must be an arrary", nk)
		}

		if dtype == "boolean" {
			if len(nvList) == 0 {
				return nil, fmt.Errorf("value for boolean column can not be empty")
			}
			if len(nvList) != 1 {
				return nil, fmt.Errorf("%v is a boolean column and must have one boolean value", nk)
			}
			booleanValue, isBoolean := nvList[0].(bool)
			if !isBoolean {
				return nil, fmt.Errorf("invalid value=(%v).  %v is a boolean column and must have one boolean value", nvList[0], nk)
			}

			all_and_strings = append(all_and_strings, fmt.Sprintf("%v = %v", nk, booleanValue))
		} else if dtype == "text" || dtype == "varchar" {
			if len(nvList) == 0 {
				return nil, fmt.Errorf("%v can not have an empty arrary", nk)
			}

			valueString := ""
			for i, v := range nvList {
				isLast := i == len(nvList)-1
				if isLast {
					if v != nil {
						valueString += fmt.Sprintf("'%v'", v)
					} else {
						valueString += "NULL"
					}
				} else {
					if v != nil {
						valueString += fmt.Sprintf("'%v',", v)
					} else {
						valueString += "NULL,"
					}
				}
			}
			all_and_strings = append(all_and_strings, fmt.Sprintf("%v IN (%v)", nk, valueString))
		} else {
			if len(nvList) == 0 {
				return nil, fmt.Errorf("%v can not have an empty arrary", nk)
			}

			valueString := ""

			// we will loop over orvlist
			for i, v := range nvList {
				isLast := i == len(nvList)-1
				if isLast {
					if v != nil {
						valueString += fmt.Sprintf("%v", v)
					} else {
						valueString += "NULL"
					}
				} else {
					if v != nil {
						valueString += fmt.Sprintf("%v,", v)
					} else {
						valueString += "NULL,"
					}
				}
			}
			all_and_strings = append(all_and_strings, fmt.Sprintf(" %v IN (%v)", nk, valueString))
		}
	}

	mainString := ""
	var all_string []string = []string{}

	for _, str := range all_or_strings {
		if len(all_string) != 0 {
			mainString += fmt.Sprintf("OR %v", str)
			all_string = append(all_string, str)
		} else {
			mainString += fmt.Sprintf("%v", str)
			all_string = append(all_string, str)
		}
	}
	for _, str := range all_and_strings {
		if len(all_string) != 0 {
			mainString += fmt.Sprintf(" AND %v", str)
			all_string = append(all_string, str)
		} else {
			mainString += fmt.Sprintf(" %v", str)
			all_string = append(all_string, str)
		}
	}

	query := fmt.Sprintf(`DELETE FROM "public"."%v" WHERE %v RETURNING *;`, arg.Table, mainString)
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		println(query)
		println(err.Error())
		return nil, err
	}
	cols, _ := rows.Columns()
	results, err := sqlRowsToJson(rows, cols)
	return results, err
}

type UpdateRowsParams struct {
	Table  string                   `json:"table" validate:"required,alphanum,gte=3,lte=60"`
	UserID int64                    `json:"user_id" validate:"required,numeric,min=1"`
	Rows   []map[string]interface{} `json:"rows" validate:"required"`
}

func (store *SQLStore) UpdateRows(ctx context.Context, arg UpdateRowsParams) error {
	validate := validator.New()
	err := validate.Struct(arg)
	if err != nil {
		return err
	}

	// First we will validate so that we dont make interactions in database with invalid data
	// Requirements: Table Schema
	dbtable, err := store.GetTableByUserIdAndTableName(ctx, GetTableByUserIdAndTableNameParams{Name: arg.Table, UserID: arg.UserID})
	// If schema is not found then probably table does not exits or does not belongs to user
	// In both cases we will send table not found.
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(TABLE_NOT_FOUND)
		}
		return err
	}
	schema, err := dbtable.Schema()
	if err != nil {
		return err
	}
	colsThatDontExist := schema.ColumnsThatDontExists(arg.Rows)
	if len(colsThatDontExist) != 0 {
		return fmt.Errorf("columns %v dosen't exists", colsThatDontExist)
	}
	var updateStrings []string

	primaryColumn := schema.PrimaryColumn()

	for i, r := range arg.Rows {
		row := Row{Value: r}
		istr, err := schema.updateRowString(row, primaryColumn, i)
		if err != nil {
			return err
		}
		updateStrings = append(updateStrings, istr)
	}
	// to use and read colums properly we need to format table columns
	updateString := strings.Join(updateStrings, "\n")
	// If insert strings are safely build without erros then execute statements
	_, err = store.db.ExecContext(ctx, updateString)
	return err
}
