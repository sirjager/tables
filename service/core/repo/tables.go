package core_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func (t *Table) Schema() (TableSchema, error) {
	var err error
	var s TableSchema
	var c []Column
	err = json.Unmarshal([]byte(t.Columns), &c)
	if err != nil {
		return s, err
	}
	return TableSchema{ID: t.ID, UserID: t.UserID, Name: t.Name, Columns: c, Created: t.Created, Updated: t.Updated}, err
}

func (t *TableSchema) ColumnsThatDontExists(m []map[string]interface{}) []string {
	invalidColumns := []string{}
	for ri, r := range m {
		for k := range r {
			exists := false
			for _, c := range t.Columns {
				if c.Name == k {
					exists = true
					break
				}
			}
			if !exists {
				invalidColumns = append(invalidColumns, fmt.Sprintf("%v(row#%d)", k, ri+1))
			}
		}
	}
	return invalidColumns
}

func (s *TableSchema) RequiredColumns() []Column {
	var c []Column
	for _, f := range s.Columns {
		if f.Required {
			c = append(c, f)
		}
	}
	return c
}

// Though there will only one
func (s *TableSchema) PrimaryColumn() Column {
	var c Column
	for _, f := range s.Columns {
		if f.Primary {
			c = f
		}
	}
	return c
}

func (s *TableSchema) UniqueColumns() []Column {
	var c []Column
	for _, f := range s.Columns {
		if f.Unique {
			c = append(c, f)
		}
	}
	return c
}

/*
This will check if a table has any not nullable column.
If Table has any not nullable column then We must check our rows has that column and value should not be null.
If any required column is not present in row then we must return the missing column or column that
have null values.
*/
func (s *TableSchema) ValidateRequiredColumns(rows []map[string]interface{}) error {
	reqCols := s.RequiredColumns()
	if len(reqCols) == 0 {
		return nil
	}
	for rowIndex, row := range rows { // row is a map
		for _, rqCol := range reqCols {
			exits := false
			valid := false
			for k, v := range row {
				if k == rqCol.Name {
					exits = true
					if v != nil {
						valid = true
					}
				}
			}
			if !exits {
				return fmt.Errorf("missing required column [%s] in row [%d]", rqCol.Name, rowIndex+1)
			}
			if !valid {
				return fmt.Errorf("null value given for not nullable column [%s] in row [%d]", rqCol.Name, rowIndex+1)
			}
		}
	}
	return nil
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

func FormatTableEntryToTable(coreTable Table) (TableSchema, error) {
	var err error
	var columns []Column
	var table TableSchema
	err = json.Unmarshal([]byte(coreTable.Columns), &columns)
	if err != nil {
		return table, err
	}
	return TableSchema{
		ID:      coreTable.ID,
		UserID:  coreTable.UserID,
		Name:    coreTable.Name,
		Columns: columns,
		Created: coreTable.Created,
		Updated: coreTable.Updated,
	}, err
}

type CreateTableTxParams struct {
	Table   string   `json:"table" validate:"required,alphanum,gte=3,lte=60"`
	UserID  int64    `json:"user_id" validate:"required,numeric,min=1"`
	Columns []Column `json:"columns" validate:"required"`
}

func (store *SQLStore) CreateTableTx(ctx context.Context, arg CreateTableTxParams) (TableSchema, error) {
	var result TableSchema

	validate := validator.New()
	err := validate.Struct(arg)
	if err != nil {
		return result, err
	}

	var all_columns_string string = ""
	// Process columns
	//  1. validate column : types, name
	// 2. generate column string:  NAME VARCHAR(50) NOT NULL
	// 3. We will make sure that user gives one primary column
	primaryColumns := []string{}
	for i, col := range arg.Columns {
		column_string, err := getColumnString(col)
		if err != nil {
			return result, err
		}
		if i == len(arg.Columns)-1 {
			// If last column then dont add comma (,)
			all_columns_string = all_columns_string + column_string
		} else {
			// If not a last column add comma (,)
			all_columns_string = all_columns_string + column_string + ", "
		}
		if col.Primary {
			primaryColumns = append(primaryColumns, col.Name)
		}
	}
	if len(primaryColumns) == 0 {
		return result, fmt.Errorf("table must contain one primary column")
	}
	if len(primaryColumns) > 1 {
		return result, fmt.Errorf("table can not have multiple primary columns, you have %s as primary columns", primaryColumns)
	}
	// Build columns json string
	columnsString, err := ColumnsToJsonString(arg.Columns)
	if err != nil {
		return result, err
	}

	// if no error mean this table is new

	err = store.execTx(ctx, func(q *Queries) error {
		// 	Create Real Table in database
		create_string := fmt.Sprintf("CREATE TABLE %s ( %s );", arg.Table, all_columns_string)
		_, err = q.db.ExecContext(ctx, create_string)
		if err != nil {
			// If any error occured while creating real table then we will delete table entry
			q.DeleteTableWhereName(ctx, arg.Table)
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "duplicate_table":
					return fmt.Errorf("table with name=(%s) already exists", arg.Table)
				}
			}
			return err
		}

		// Store created table details: name, user.id , columns json as string
		created_table, err := q.CreateTable(ctx, CreateTableParams{Name: arg.Table, UserID: arg.UserID, Columns: columnsString})
		// if any error occurs return err
		// 1. Error will occur if table with same name aleady exists, uniqye key violation on tablename
		// 2. any other database error
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					return fmt.Errorf("table=(%s) already exists", arg.Table)
				}
			}
			return err
		}

		result, err = FormatTableEntryToTable(created_table)

		return err
	})

	return result, err
}

type Name struct {
	Value string `json:"name" validate:"required,alphanum,gte=3,lte=60"`
}

/*
This function does not check for users table this delete the table whose name is provided
*/
func (store *SQLStore) DropTableTx(ctx context.Context, arg Name) error {
	validate := validator.New()
	// Validate Name
	err := validate.Struct(arg)
	if err != nil {
		return err
	}

	err = store.execTx(ctx, func(q *Queries) error {
		drop_table_string := fmt.Sprintf(`DROP TABLE "public".%v;`, arg.Value)
		// First we will Drop the table. This will give error if table does not exists
		// which is good. We dont have to fetch/look for table if it exists or not
		_, err = q.db.ExecContext(ctx, drop_table_string)
		if err != nil {
			return err
		}
		// No Error then we remove entry also
		err := q.DeleteTableWhereName(ctx, arg.Value)
		return err
	})
	return err
}
