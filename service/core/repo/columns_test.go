package core_repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddColumns(t *testing.T) {
	store := NewStore(testDb)
	user := createRandomUser(t)
	table := createRandomTable(t, user)
	updated, err := store.AddColumnTx(context.Background(), AddColumnsTxParams{
		Table:  table.Name,
		UserID: table.UserID,
		Columns: []Column{
			{Name: "col1", Type: "varchar", Length: 50, Unique: true, Required: true},
			{Name: "col2", Type: "integer", Precision: 10, Scale: 4, Unique: true, Required: true},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, updated)
	require.Equal(t, table.Name, updated.Name)
	require.Equal(t, table.UserID, updated.UserID)
	require.Equal(t, len(table.Columns), len(updated.Columns)-2)
}
