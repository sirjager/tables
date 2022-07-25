package core_repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInfoGetCoreTables(t *testing.T) {
	items, err := testQueries.InfoGetCoreTables(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, items)
}
