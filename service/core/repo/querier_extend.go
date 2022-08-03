package core_repo

import "context"

type QuerierTx interface {
	CreateTableTx(ctx context.Context, arg CreateTableTxParams) (TableSchema, error)
	DropTableTx(ctx context.Context, arg Name) error
	GetRows(ctx context.Context, arg GetRowsParams) ([]any, error)
	DeleteRows(ctx context.Context, arg DeleteRowsParams) ([]any, error)
	InsertRows(ctx context.Context, arg InsertRowsParams) error
	UpdateRows(ctx context.Context, arg UpdateRowsParams) error
	AddColumnTx(ctx context.Context, arg AddColumnsTxParams) (TableSchema, error)
	DropColumnTx(ctx context.Context, arg DropColumnsTxParams) (TableSchema, error)
}
