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
			TableName: utils.RandomString(8),
			Uid:       user.Uid,
			Columns:   columns,
		}
		createdTable, err := store.CreateTableTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, createdTable)
		require.Equal(t, arg.TableName, createdTable.Tablename)
		require.Equal(t, arg.Uid, createdTable.Uid)
		require.NotZero(t, createdTable.Tid)

		foundCreatedTable, err := store.GetCoreTableWithTid(context.Background(), createdTable.Tid)
		require.NoError(t, err)
		require.NotEmpty(t, foundCreatedTable)

		require.Equal(t, createdTable.Tid, foundCreatedTable.Tid)
		require.Equal(t, createdTable.Uid, foundCreatedTable.Uid)
		require.Equal(t, createdTable.Uid, arg.Uid)
		require.Equal(t, createdTable.Tablename, foundCreatedTable.Tablename)
		require.Equal(t, createdTable.Tablename, arg.TableName)
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
			TableName: utils.RandomString(8),
			Uid:       user.Uid,
			Columns:   columns,
		}
		createdTable, err := store.CreateTableTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, createdTable)
		require.Equal(t, arg.TableName, createdTable.Tablename)
		require.Equal(t, arg.Uid, createdTable.Uid)
		require.NotZero(t, createdTable.Tid)
		createdTableNames = append(createdTableNames, createdTable.Tablename)
	}

	require.Equal(t, count, len(createdTableNames))
	for _, name := range createdTableNames {
		err := store.DropTableTx(context.Background(), RemoveCoreTableWithUidAndNameParams{Uid: user.Uid, Tablename: name})
		require.NoError(t, err)
	}
}
