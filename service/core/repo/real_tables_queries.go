package core_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

type Column struct {
	Name      string `json:"name" validate:"required,alphanum,gte=1,lte=30"`
	Type      string `json:"type" validate:"required,oneof=integer smallint bigint decimal numeric real 'double precision' smallserial serial bigserial varchar char character text timestamp 'timestamp with time zone' 'timestamp without time zone' date 'time with time zone' time 'time without time zone' bool boolean bit 'bit varying' cidr inet macaddr macaddr8 json jsonb money uuid"`
	Length    int64  `json:"length"`
	Primary   bool   `json:"primary"`
	Unique    bool   `json:"unique"`
	Required  bool   `json:"required"`
	Precision int32  `json:"precision"`
	Scale     int32  `json:"scale"`
	Default   string `json:"default"`
}

type RealTable struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	UserID  int64     `json:"user_id"`
	Columns []Column  `json:"columns"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func getColumnString(col Column) (string, error) {
	var err error
	var validate = validator.New()
	err = validate.Struct(col)
	if err != nil {
		return "", err
	}

	if _, err := strconv.Atoi(string(col.Name[0])); err == nil {
		return "", fmt.Errorf("invalid column name %s . Column name can not start with a number", col.Name)
	}

	main_string := col.Name
	column_type := ""
	upper_col_type := strings.ToUpper(col.Type)
	lower_col_type := strings.ToLower(col.Type)
	switch lower_col_type {
	// First three can be simplified to one case but it will become too lenthy
	// so spliting into three for simplicity
	case "integer", "smallint", "bigint", "smallserial", "serial", "bigserial", "real", "money", "uuid":
		column_type = upper_col_type
	case "double precision", "json", "jsonb", "boolean", "cidr", "inet", "macaddr", "macaddr8":
		column_type = upper_col_type
	case "date", "time", "timestamp", "timestamp with time zone", "timestamp without time zone", "time with time zone", "time without time zone":
		column_type = upper_col_type
	case "decimal", "numeric":
		if col.Precision > 0 {
			if col.Scale > 0 {
				column_type = fmt.Sprintf("%s(%d,%d)", upper_col_type, col.Precision, col.Scale)
				break
			} else {
				column_type = fmt.Sprintf("%s(%d)", upper_col_type, col.Precision)
				break
			}
		} else {
			column_type = upper_col_type
			break
		}
	case "varchar", "char", "character", "text", "bit", "bit varying":
		if col.Length > 0 {
			column_type = fmt.Sprintf("%s(%d)", upper_col_type, col.Length)
		} else {
			column_type = upper_col_type
		}
	default:
		err = fmt.Errorf("column=(%s) contains invalid type=(%s)", col.Name, col.Type)
	}
	if len(column_type) < 1 {
		err = fmt.Errorf("column=(%s) contains invalid type=(%s)", col.Name, col.Type)
	}
	if err != nil {
		return "", err
	}
	main_string += " " + column_type
	if col.Primary {
		main_string += " PRIMARY KEY"
	}
	if col.Unique {
		main_string += " UNIQUE"
	}
	if col.Required {
		main_string += " NOT NULL"
	}
	if len(col.Default) > 0 {
		main_string += " DEFAULT (" + col.Default + ")"
	}
	return main_string, err
}

func getColumnsFromSchema(coreTable CoreTable) ([]Column, error) {
	schema_bytes, err := json.Marshal(coreTable)
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	table_schema_string := []byte(string(schema_bytes))
	if err := json.Unmarshal(table_schema_string, &raw); err != nil {
		return nil, err
	}

	var mycolumns []Column
	err = json.Unmarshal([]byte(fmt.Sprintf("%s", raw["columns"])), &mycolumns)
	if err != nil {
		return nil, err
	}
	return mycolumns, nil

}

func GenerateFieldString(fields []string, columns []Column) (string, error) {
	fieldString := ""
	if len(fields) == 0 {
		return "*", nil
	}
	var fieldsThatDontExists []string = []string{}
	for i, f := range fields {
		exits := false
		for _, c := range columns {
			if f == c.Name {
				exits = true
				break
			}
		}
		if !exits {
			fieldsThatDontExists = append(fieldsThatDontExists, f)
		} else {
			// if last field
			if i == len(fields)-1 {
				fieldString += f
			} else {
				fieldString += f + ","
			}
		}
	}
	if len(fieldsThatDontExists) != 0 {
		return "", fmt.Errorf("%v dont exits", fieldsThatDontExists)
	}
	return fieldString, nil
}

func ColumnsToJsonString(columns []Column) (string, error) {
	// Build columns json string
	column_bytes, err := json.Marshal(columns)
	if err != nil {
		return "", err
	}
	return string(column_bytes), err
}

func FormatTableEntryToTable(coreTable CoreTable) (RealTable, error) {
	var err error
	var columns []Column
	var table RealTable
	err = json.Unmarshal([]byte(coreTable.Columns), &columns)
	if err != nil {
		return table, err
	}
	return RealTable{
		ID:      coreTable.ID,
		UserID:  coreTable.UserID,
		Name:    coreTable.Name,
		Columns: columns,
		Created: coreTable.Created,
		Updated: coreTable.Updated,
	}, err
}

type CreateTableTxParams struct {
	Name    string   `json:"table" binding:"required,gte=3,lte=60"`
	UserID  int64    `json:"user_id" binding:"required,numeric,min=1"`
	Columns []Column `json:"columns" binding:"required"`
}

func (store *SQLStore) CreateTableTx(ctx context.Context, arg CreateTableTxParams) (RealTable, error) {
	var result RealTable
	err := store.execTx(ctx, func(q *Queries) error {
		var all_columns_string string = ""
		// Process columns
		//  1. validate column : types, name
		// 2. generate column string:  NAME VARCHAR(50) NOT NULL
		for i, col := range arg.Columns {
			column_string, err := getColumnString(col)
			if err != nil {
				return err
			}
			if i == len(arg.Columns)-1 {
				// If last column then dont add comma (,)
				all_columns_string = all_columns_string + column_string
			} else {
				// If not a last column add comma (,)
				all_columns_string = all_columns_string + column_string + ", "
			}
		}
		// Build columns json string
		columnsString, err := ColumnsToJsonString(arg.Columns)
		if err != nil {
			return err
		}

		// if no error mean this table is new
		// 	Create Real Table in database
		create_string := fmt.Sprintf("CREATE TABLE %s ( %s );", arg.Name, all_columns_string)
		_, err = q.db.ExecContext(ctx, create_string)
		if err != nil {
			// If any error occured while creating real table then we will delete table entry
			q.DeleteTableWhereUserAndName(ctx, DeleteTableWhereUserAndNameParams{UserID: arg.UserID, Name: arg.Name})
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "duplicate_table":
					return fmt.Errorf("table with name=(%s) already exists", arg.Name)
				}
			}
			return err
		}

		// Store created table details: name, user.id , columns json as string
		created_table, err := q.CreateTable(ctx, CreateTableParams{Name: arg.Name, UserID: arg.UserID, Columns: columnsString})
		// if any error occurs return err
		// 1. Error will occur if table with same name aleady exists, uniqye key violation on tablename
		// 2. any other database error
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					return fmt.Errorf("table=(%s) already exists", arg.Name)
				}
			}
			return err
		}
		result, err = FormatTableEntryToTable(created_table)
		return err
	})
	return result, err
}

func (store *SQLStore) DropTableTx(ctx context.Context, arg DeleteTableWhereUserAndNameParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		drop_table_string := fmt.Sprintf("DROP TABLE IF EXISTS %s;", arg.Name)
		// First we will delete entry
		err := q.DeleteTableWhereUserAndName(ctx, DeleteTableWhereUserAndNameParams{UserID: arg.UserID, Name: arg.Name})
		if err != nil {
			return err
		}
		// Then drop table
		_, err = q.db.ExecContext(ctx, drop_table_string)
		if err != nil {
			return err
		}
		return err
	})
	return err
}
