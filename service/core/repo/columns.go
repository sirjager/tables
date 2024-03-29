package core_repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type AddColumnsTxParams struct {
	Table   string   `json:"table" validate:"required,alphanum,gte=3,lte=60"`
	Columns []Column `json:"columns" validate:"required"`
	UserID  int64    `json:"user_id" validate:"required,numeric,min=1"`
}

func (store *SQLStore) AddColumnTx(ctx context.Context, arg AddColumnsTxParams) (TableSchema, error) {
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
	for i, col := range arg.Columns {
		column_string, err := getColumnString(col)
		if err != nil {
			return result, err
		}
		if i == len(arg.Columns)-1 {
			// If last column then dont add comma (,)
			// ADD COLUMN fax VARCHAR,
			all_columns_string = all_columns_string + " ADD COLUMN " + column_string
		} else {
			// If not a last column add comma (,)
			all_columns_string = all_columns_string + " ADD COLUMN " + column_string + ", "
		}
	}

	coretable, err := store.GetTableByUserIdAndTableName(ctx, GetTableByUserIdAndTableNameParams{Name: arg.Table, UserID: arg.UserID})
	// This will make sure table belongs to user
	if err != nil {
		if err == sql.ErrNoRows {
			return result, fmt.Errorf("table %s not found", arg.Table)
		}
		return TableSchema{}, err
	}

	// Now before creating real table we will check if columns exists
	mytable, err := FormatTableEntryToTable(coretable)
	if err != nil {
		return result, err
	}
	var alreadyExistingColumns []string = []string{}
	for _, newCol := range arg.Columns {
		for _, existCol := range mytable.Columns {
			if newCol.Name == existCol.Name {
				alreadyExistingColumns = append(alreadyExistingColumns, newCol.Name)
			}
			if newCol.Primary {
				// Check if table already have primary key
				for _, cc := range mytable.Columns {
					if cc.Primary {
						return result, fmt.Errorf("'%s' already have primary key, a table can not have more multiple primary keys", cc.Name)
					}
				}
			}
		}
	}

	if len(alreadyExistingColumns) > 0 {
		return result, fmt.Errorf("%s already exists in the table", alreadyExistingColumns)
	}
	err = store.execTx(ctx, func(q *Queries) error {

		/* At this point we have new valid columns to be added in table which belongs to the user
		now we will first add columns in real table. then we will update the record
		*/
		alterTableString := fmt.Sprintf(`ALTER TABLE "public"."%s" %s;`, arg.Table, all_columns_string)
		_, err = q.db.ExecContext(ctx, alterTableString)
		if err != nil {
			return err
		}
		// if no errors the we must update the records
		mytable.Columns = append(mytable.Columns, arg.Columns...)
		// Build columns json string
		columnsString, err := ColumnsToJsonString(mytable.Columns)
		if err != nil {
			return err
		}
		updatedTable, err := q.UpdateTableColumns(ctx, UpdateTableColumnsParams{ID: mytable.ID, Columns: columnsString})
		if err != nil {
			return err
		}
		result, err = FormatTableEntryToTable(updatedTable)
		return err
	})
	return result, err
}

type DropColumnsTxParams struct {
	Table   string   `json:"table" validate:"required,alphanum,gte=3,lte=60"`
	UserID  int64    `json:"user_id" validate:"required,numeric,min=1"`
	Columns []string `json:"columns" validate:"required"`
}

func (store *SQLStore) DropColumnTx(ctx context.Context, arg DropColumnsTxParams) (TableSchema, error) {
	var result TableSchema
	// First we will get table and check if it belongs to the user
	coretable, err := store.GetTableByUserIdAndTableName(ctx, GetTableByUserIdAndTableNameParams{Name: arg.Table, UserID: arg.UserID})
	if err != nil {
		if err == sql.ErrNoRows {
			return result, fmt.Errorf("table %s not found", arg.Table)
		}
		return result, err
	}

	mytable, err := FormatTableEntryToTable(coretable)
	if err != nil {
		return result, err
	}

	// No we will check if the columns exists or not
	//
	numberOfColumnsToDelete := 0
	for _, colToDel := range arg.Columns {
		var columExists bool
		for _, existingCol := range mytable.Columns {
			if colToDel == existingCol.Name {
				// This means columns does exists
				numberOfColumnsToDelete += 1
				columExists = true
			}
		}
		if !columExists {
			return result, fmt.Errorf("column %s does not exists", colToDel)
		}
	}
	// We will check if user wants to delete all the columns or not
	// if he/she does wants to delete all the columns then we send error: delete table instead
	if numberOfColumnsToDelete == len(mytable.Columns) {
		return result, fmt.Errorf("table has %d columns and and you are deleting %d columns. which literally means you want all columns gone, delete table instead",
			len(mytable.Columns), numberOfColumnsToDelete)
	}
	// at this point we have table which belongs to the user
	// valid names of columns which user wants to delete from his table
	// Now we build drop column strings
	var all_columns_string string = ""
	existingColumns := mytable.Columns
	for i, delcol := range arg.Columns {
		isLastColumn := i == len(arg.Columns)-1
		for j, existingCol := range existingColumns {
			if delcol == existingCol.Name {
				// name matches then this is the column we want to delete
				if isLastColumn {
					all_columns_string += " DROP COLUMN " + delcol
				} else {
					all_columns_string += " DROP COLUMN " + delcol + ","
				}
				existingColumns = append(existingColumns[:j], existingColumns[j+1:]...)
			}
		}
	}
	if all_columns_string == "" {
		return result, fmt.Errorf("no columns deleted")
	}

	updatedColumnsString, err := ColumnsToJsonString(existingColumns)
	if err != nil {
		return result, err
	}
	err = store.execTx(ctx, func(q *Queries) error {
		// No make Drop Column string
		alterTableString := fmt.Sprintf(`ALTER TABLE "public"."%s" %s;`, arg.Table, all_columns_string)
		// At this point all validations is done we will start making changes
		_, err = q.db.ExecContext(ctx, alterTableString)
		if err != nil {
			println(err.Error())
			return err
		}
		updatedTable, err := q.UpdateTableColumns(ctx, UpdateTableColumnsParams{ID: mytable.ID, Columns: updatedColumnsString})
		if err != nil {
			return err
		}
		result, err = FormatTableEntryToTable(updatedTable)
		return err
	})
	return result, err
}
