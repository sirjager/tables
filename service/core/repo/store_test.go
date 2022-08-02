package core_repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/SirJager/tables/service/core/utils"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestCreateTableTx(t *testing.T) {
	store := NewStore(testDb)
	user := createRandomUser(t)
	count := 5
	validate := validator.New()
	for i := 0; i < count; i++ {
		var columns []Column
		for c := 0; c < 3; c++ {
			col := Column{Name: utils.RandomString(5)}
			if c == 0 {
				col.Type = "integer"
			} else if c == 1 {
				col.Type = "boolean"
			} else {
				col.Type = "varchar"
			}
			err := validate.Struct(col)
			if err != nil {
				require.NoError(t, err)
			}
			columns = append(columns, col)
		}
		arg := CreateTableTxParams{
			Table:   utils.RandomString(8),
			UserID:  user.ID,
			Columns: columns,
		}
		createdTable, err := store.CreateTableTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, createdTable)
		require.Equal(t, arg.Table, createdTable.Name)
		require.Equal(t, arg.UserID, createdTable.UserID)
		require.NotZero(t, createdTable.ID)

		foundCreatedTable, err := store.GetTableByUserIdAndTableName(context.Background(),
			GetTableByUserIdAndTableNameParams{Name: createdTable.Name, UserID: createdTable.UserID})
		require.NoError(t, err)
		require.NotEmpty(t, foundCreatedTable)

		require.Equal(t, createdTable.ID, foundCreatedTable.ID)
		require.Equal(t, createdTable.UserID, foundCreatedTable.UserID)
		require.Equal(t, createdTable.UserID, arg.UserID)
		require.Equal(t, createdTable.Name, foundCreatedTable.Name)
		require.Equal(t, createdTable.Name, arg.Table)
	}

}

func TestDropTableTx(t *testing.T) {
	store := NewStore(testDb)
	user := createRandomUser(t)
	count := 5
	validate := validator.New()
	var createdTableNames []string
	for i := 0; i < count; i++ {
		var columns []Column
		for c := 0; c < 3; c++ {
			col := Column{Name: utils.RandomString(5)}
			if c == 0 {
				col.Type = "integer"
			} else if c == 1 {
				col.Type = "boolean"
			} else {
				col.Type = "varchar"
			}
			err := validate.Struct(col)
			if err != nil {
				require.NoError(t, err)
			}
			columns = append(columns, col)
		}
		arg := CreateTableTxParams{
			Table:   utils.RandomString(8),
			UserID:  user.ID,
			Columns: columns,
		}
		createdTable, err := store.CreateTableTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, createdTable)
		require.Equal(t, arg.Table, createdTable.Name)
		require.Equal(t, arg.UserID, createdTable.UserID)
		require.NotZero(t, createdTable.ID)
		createdTableNames = append(createdTableNames, createdTable.Name)
	}

	require.Equal(t, count, len(createdTableNames))
	for _, name := range createdTableNames {
		err := store.DropTableTx(context.Background(), Name{Value: name})
		require.NoError(t, err)
	}
}

func TestAddColumns(t *testing.T) {
	store := NewStore(testDb)
	_, err := store.AddColumnTx(context.Background(), AddColumnsTxParams{
		Table: "mytable",
		Columns: []Column{
			{Name: "col1", Type: "varchar", Length: 50, Unique: true, Required: true},
			{Name: "col2", Type: "integer", Precision: 10, Scale: 4, Unique: true, Required: true},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestDropColumns(t *testing.T) {
	store := NewStore(testDb)
	_, err := store.DropColumnTx(context.Background(), DropColumnsTxParams{
		Table:   "mytable",
		Columns: []string{"col1", "col2"},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
