package core_repo

import (
	"context"
)

const countUsers = `SELECT COUNT(*) FROM "public"."core_users";`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	rows := q.db.QueryRowContext(ctx, countUsers)
	err := rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

const countUsersWhereBooleanColumnIs = `SELECT COUNT(*) FROM "public"."core_users" WHERE $1 = $2;`

/*This counts any boolean column in users table

boolean columns: "verified" | "blocked" | "public"

example:
	totalVerifiedUsers,err := CountUsersWhereBooleanColumnIs(context.Background(), "verified", true)
	totalBlockedUsers,err := CountUsersWhereBooleanColumnIs(context.Background(), "blocked", true)
	totalPrivateUsers,err := CountUsersWhereBooleanColumnIs(context.Background(), "public", false)

println(count) // int64 value */
func (q *Queries) CountUsersWhereBooleanColumnIs(ctx context.Context, column string, value bool) (int64, error) {
	var count int64
	rows := q.db.QueryRowContext(ctx, countUsersWhereBooleanColumnIs, column, value)
	err := rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}
