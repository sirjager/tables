package core_repo

import (
	"context"
	"database/sql"
	"testing"

	"github.com/SirJager/tables/service/core/utils"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func createRandomTable(t *testing.T, user User) TableSchema {
	store := NewStore(testDb)
	validate := validator.New()
	var columns []Column
	for c := 0; c < 3; c++ {
		col := Column{Name: utils.RandomString(5)}
		if c == 0 {
			col.Primary = true
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

	err := validate.Struct(arg)
	require.NoError(t, err)

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
	return createdTable
}

func TestCreateTableTx(t *testing.T) {
	user := createRandomUser(t)
	createRandomTable(t, user)
}

func TestDropTableTx(t *testing.T) {
	store := NewStore(testDb)
	user := createRandomUser(t)
	count := 5
	var createdTableNames []string
	for i := 0; i < count; i++ {
		table := createRandomTable(t, user)
		createdTableNames = append(createdTableNames, table.Name)
	}
	require.Equal(t, count, len(createdTableNames))
	for _, name := range createdTableNames {
		err := store.DropTableTx(context.Background(), Name{Value: name})
		require.NoError(t, err)
		_, err = store.GetTableWhereName(context.Background(), name)
		require.EqualError(t, err, sql.ErrNoRows.Error())
	}
}
