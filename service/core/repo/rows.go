package core_repo

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type KeyValueParams struct {
	K string `json:"k"`
	V string `json:"v"`
}

type InsertRowsParams struct {
	Table  string             `json:"table" validate:"required,alphanum,min=1"`
	UserID int64              `json:"user_id" validate:"required,numeric,min=1"`
	Rows   [][]KeyValueParams `json:"rows" validate:"required"`
}

func (store *SQLStore) InsertRows(ctx context.Context, arg InsertRowsParams) error {
	validate := validator.New()
	var err error
	err = validate.Struct(arg)
	if err != nil {
		return err
	}
	err = store.execTx(ctx, func(q *Queries) error {
		// First we will retrieve table schema for column data type validation
		table_schema, err := q.GetTableWhereNameAndUser(ctx, GetTableWhereNameAndUserParams{Name: arg.Table, UserID: arg.UserID})
		if err != nil {
			return err
		}
		// From table schema generate struct for all columns
		mycolumns, err := getColumnsFromSchema(table_schema)
		if err != nil {
			return err
		}

		insert_string := ""

		// Process each rows
		// 1. Validate data types
		// 2. generate insert string : INSERT INTO tablename (c1,c2,c3) VALUES(v1,v2,v3);

		for _, entry := range arg.Rows {
			entry_string := "INSERT INTO " + arg.Table + " "
			columns_string := ""
			values_string := ""
			for c, column := range entry {
				isLegitColumn := false
				var columnType string
				for _, col := range mycolumns {
					if column.K == col.Name {
						// If column name is same as stored in our table schema then column is legit
						isLegitColumn = true
						// store column data type for coming steps
						columnType = col.Type
					}
				}

				// If column is not legin send back error with details
				if !isLegitColumn {
					return fmt.Errorf("column=(%s) not found in table=(%s)", column.K, arg.Table)
				}

				// If column is legit then continue

				isLastColumn := c == len(entry)-1
				if isLastColumn {
					// if last column then after column name put nothing
					// example: INSERT INTO STUDENTS (email,username,age)
					columns_string = columns_string + column.K

					if columnType == "varchar" || columnType == "text" {
						values_string = values_string + "'" + column.V + "'"
					} else {
						values_string = values_string + column.V
					}
				} else {
					// if not a last column then after column name put (,) comma
					// example: INSERT INTO STUDENTS (email,username,age)
					columns_string = columns_string + column.K + ","
					if columnType == "varchar" || columnType == "text" {
						values_string = values_string + "'" + column.V + "'" + ","
					} else {
						values_string = values_string + column.V + ","
					}
				}
			}

			// Build full string
			entry_string = entry_string + "(" + columns_string + ")"
			entry_string = entry_string + " VALUES (" + values_string + "); "
			insert_string += entry_string

		}

		// If insert strings are safely build without erros then execute statements
		// we dont wanna send back same data which user just inserted
		_, err = q.db.ExecContext(ctx, insert_string)

		// In the end we will send err object
		// this is will be nil if process is successful
		// this is will contain error if not successful
		return err
	})
	return err
}

type DeleteRowsParams struct {
	Table   string                   `json:"table" validate:"required,alphanum,min=1"`
	UserID  int64                    `json:"user_id" validate:"required,numeric,min=1"`
	Filters map[string][]interface{} `json:"filters" validate:"required,gte=1"`
}

func (q *Queries) DeleteRows(ctx context.Context, arg DeleteRowsParams) error {
	var err error
	validate := validator.New()
	err = validate.Struct(arg)
	if err != nil {
		return err
	}

	// Get Table schema
	table, err := q.GetTableWhereNameAndUser(ctx, GetTableWhereNameAndUserParams{Name: arg.Table, UserID: arg.UserID})
	if err != nil {
		return err
	}

	mytable, err := FormatTableEntryToTable(table)
	if err != nil {
		return err
	}

	// This will extract all the keys (all column names)
	columns := make([]string, len(arg.Filters))
	//DELETE FROM links WHERE id IN (6,5) RETURNING *;
	i := 0
	mainExecuteString := ""
	for col, v := range arg.Filters {
		// We got the column name now we need to know what data type the column is
		colType := ""
		columnExists := false
		for _, mycol := range mytable.Columns {
			if col == mycol.Name {
				colType = mycol.Type
				columnExists = true
				break
			}
		}
		// We will also check if column given by user exists or not.
		// this can save us unawanted invalid interaction with database
		if !columnExists {
			return fmt.Errorf("column [%s] does not exits", col)
		}

		columns[i] = col
		i++
		mainString := "DELETE FROM " + arg.Table + " WHERE " + col + " IN ("
		for valueIndex, value := range v {
			isLast := valueIndex == len(v)-1
			if isLast {
				// We will check if the value is a text then we need single quote
				if colType == "varchar" || colType == "text" {
					mainString += fmt.Sprintf("'%v');", value)
				} else {
					mainString += fmt.Sprintf("%v);", value)
				}
			} else {
				// We will check if the value is a text then we need single quote
				if colType == "varchar" || colType == "text" {
					mainString += fmt.Sprintf("'%v',", value)
				} else {
					mainString += fmt.Sprintf("%v,", value)
				}
			}
		}
		mainExecuteString += " " + mainString
	}
	_, err = q.db.ExecContext(ctx, mainExecuteString)
	return err
}

type GetRowsParams struct {
	UserID int64  `json:"user_id" validate:"required,numeric,min=1"`
	Table  string `json:"table" validate:"required,alphanum,min=1"`
}

func (q *Queries) GetRows(ctx context.Context, arg GetRowsParams) ([]any, error) {
	var err error
	validate := validator.New()
	err = validate.Struct(arg)
	if err != nil {
		return nil, err
	}

	query := " SELECT * FROM " + arg.Table + " ;"
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()

	// This will be main result and will have list of col:value
	var results []any

	// Now we will do some hard core magic i barely understand
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		fields := make(map[string]interface{})
		for k, v := range m {
			fields[k] = v
		}
		results = append(results, fields)
	}
	return results, err
}

type GetRowParams struct {
	UserID  int64                  `json:"user_id" validate:"required,numeric,min=1"`
	Table   string                 `json:"table" validate:"required,alphanum,min=1"`
	Fields  []string               `json:"fields" validate:""`
	Filters map[string]interface{} `json:"filters" validate:""`
}

func (q *Queries) GetRow(ctx context.Context, arg GetRowParams) ([]any, error) {
	var err error
	validate := validator.New()
	err = validate.Struct(arg)
	if err != nil {
		return nil, err
	}

	if len(arg.Filters) == 0 {
		return nil, fmt.Errorf("no filters provided")
	}

	/*
		So if the request is customized then we expect the request in
		below format.
		{
			"rows": {
		        "column_name" : [ "value to match","value to match" ] /
		    }
		}
		Kepoint :
			1. If the value is null then simply pass null

		# Suppose we have table with columns: id(integer), name(varchar/text), verified(boolean)
		Scenarios :
		1. Request is to get all the rows where id = 1. so the request would look like
			 {
				"rows": {
					"id" : [ 1 ]
				}
			}
			Result will be a list of all the rows where id = 1
		2. Request is to get all the rows where id = 3, 4, 5 . so the request would look like
			 {
				"rows": {
					"id" : [ 3, 4, 5 ]
				}
			}
			Result will be a list of  all the rows where id = 3 , 4 , 5
		3. Request is to get all the rows where
			name = "John Doe", "Elsa"  or id = 256 . so the request would look like
			{
				"rows": {
					"id" : [ 256 ],
					"name" : [ "John Doe","Elsa" ]
				}
			}
			Result will be a list of  all the rows where either name is "John Doe", "Elsa" or id = 256
			No Rows will be repeated in the result. Even if John Doe has id 256.
			This is Equivalted to OR in SQL
		4. Request is to get all the rows where
			name = null (NULL/ NIL) and has verifed = false. So the request will look like
			{
				"rows" : {
					"name" : [ null ],
					"&" : {
						"verified" : [ "false" ]
					}
				}
			}
			So this "&" will make sure that our result has both
			name = null and verified = false
			This is equivaltent to
			SELECT * FROM tablename WHERE name is NULL and verified = false ;

		5. Table users : id, firstname, lastname, email, username, hashedpass, phone, isBlocked, isVerified
			Request is to fetch all the users where
			firstname = null, lastname = null, isBlocked = true  and isVerified = false
			// for some reason we just want to fetch all the users with above conditions then
			request will look like:
				{
					"rows" : {
						"&" : {
							"firstname" : [null],
							"lastname" : [null],
							"isBlocked": [true],
							"isVerified": [false]
						}
					}
				}
			So the result will be a list of users where all the conditions match
			This is equivalent to
			SELECT * FROM tablename
				WHERE firstname is null
				AND lastname is null
				AND isBlocked = true
				AND isVerified = false ;

			response : [
				{
					"id" : 645,
					"firstname" null,
					"lastname" null,
					"lastname" null,
					"email" "something.@email.com",
					"username" "someusername",
					"hashedpass" "somehasedpass",
					"phone" "some phone number",
					"isBlocked" : true,
					"isVerified" : false
				}
				...
			]
	*/

	/* Sample request :
	{ "rows":
		{
			"verified": [ null ],
			"&": {
					"name": [ "user one" ]
				}
		}
	}
	*/

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

	table, err := q.GetTableWhereNameAndUser(ctx, GetTableWhereNameAndUserParams{Name: arg.Table, UserID: int64(arg.UserID)})

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

	// if len(arg.Fields) == 0 {
	// 	fieldsString = "*"
	// } else {
	// 	for i, f := range arg.Fields {
	// 		exists := false
	// 		isLast := i == len(arg.Fields)-1
	// 		for _, c := range mytable.Columns {
	// 			if f == c.Name {
	// 				exists = true
	// 				break
	// 			}
	// 		}
	// 		if !exists {
	// 			return nil, fmt.Errorf("column %v does not exits", f)
	// 		}
	// 		if isLast {
	// 			fieldsString += f
	// 		} else {
	// 			fieldsString += f + ","
	// 		}
	// 	}
	// }

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
	query := fmt.Sprintf(`SELECT %v FROM  "public"."%v" WHERE %v ;`, fieldsString, arg.Table, mainString)
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		println(query)
		println(err.Error())
		return nil, err
	}

	cols, _ := rows.Columns()

	// This will be main result and will have list of col:value
	var results []any = []any{}

	// Now we will do some hard core magic i barely understand
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		fields := make(map[string]interface{})
		for k, v := range m {
			fields[k] = v
		}
		results = append(results, fields)
	}
	return results, err
}
