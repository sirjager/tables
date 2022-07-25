package core_repo

import (
	"context"
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
			Name:    utils.RandomString(8),
			UserID:  user.ID,
			Columns: columns,
		}
		createdTable, err := store.CreateTableTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, createdTable)
		require.Equal(t, arg.Name, createdTable.Name)
		require.Equal(t, arg.UserID, createdTable.UserID)
		require.NotZero(t, createdTable.ID)

		foundCreatedTable, err := store.GetCoreTableWithTid(context.Background(), createdTable.ID)
		require.NoError(t, err)
		require.NotEmpty(t, foundCreatedTable)

		require.Equal(t, createdTable.ID, foundCreatedTable.ID)
		require.Equal(t, createdTable.UserID, foundCreatedTable.UserID)
		require.Equal(t, createdTable.UserID, arg.UserID)
		require.Equal(t, createdTable.Name, foundCreatedTable.Name)
		require.Equal(t, createdTable.Name, arg.Name)
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
			Name:    utils.RandomString(8),
			UserID:  user.ID,
			Columns: columns,
		}
		createdTable, err := store.CreateTableTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, createdTable)
		require.Equal(t, arg.Name, createdTable.Name)
		require.Equal(t, arg.UserID, createdTable.UserID)
		require.NotZero(t, createdTable.ID)
		createdTableNames = append(createdTableNames, createdTable.Name)
	}

	require.Equal(t, count, len(createdTableNames))
	for _, name := range createdTableNames {
		err := store.DropTableTx(context.Background(), RemoveCoreTableWithUidAndNameParams{UserID: user.ID, Name: name})
		require.NoError(t, err)
	}
}
