package core_repo

import "context"

type QuerierTx interface {
	CreateTableTx(ctx context.Context, arg CreateTableTxParams) (RealTable, error)
	DropTableTx(ctx context.Context, arg DeleteTableWhereUserAndNameParams) error
	GetRows(ctx context.Context, arg GetRowsParams) ([]any, error)
	InsertRows(ctx context.Context, arg InsertRowsParams) error
	AddColumnTx(ctx context.Context, arg AddColumnsTxParams) (RealTable, error)
}
