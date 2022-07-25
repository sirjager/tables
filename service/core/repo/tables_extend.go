package core_repo

import (
	"context"
)

const countTables = `SELECT COUNT(*) FROM "public"."core_tables";`

func (q *Queries) CountTables(ctx context.Context) (int64, error) {
	var count int64
	rows := q.db.QueryRowContext(ctx, countTables)
	err := rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

const countTablesWithUserID = `SELECT COUNT(*) FROM "public"."core_tables" WHERE user_id = $1;`

/*This counts any boolean column in users table

example:
	totalTablesByUser,err := CountTablesWithUserID(context.Background(), 546)

println(count) // int64 value */
func (q *Queries) CountTablesWithUserID(ctx context.Context, UserID int64) (int64, error) {
	var count int64
	rows := q.db.QueryRowContext(ctx, countTablesWithUserID, UserID)
	err := rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}
