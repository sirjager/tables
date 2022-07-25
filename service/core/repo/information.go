package core_repo

import (
	"context"
	"fmt"
)

type ColumnValue struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}
type InfoGetCountParams struct {
	Table         string        `json:"table"`
	AndConditions []ColumnValue `json:"andconditions"`
}

func (q *Queries) InfoGetCount(ctx context.Context, arg InfoGetCountParams) (int64, error) {
	var count int64
	var err error
	var query string = `SELECT COUNT(*) FROM "public"."` + arg.Table + `"`
	if len(arg.AndConditions) != 0 {
		query = query + " WHERE "
		for i, conditon := range arg.AndConditions {
			isLast := i == len(arg.AndConditions)-1
			if isLast {
				query = query + conditon.Column + " = " + conditon.Value + ";"
			} else {
				query = query + conditon.Column + " = " + conditon.Value + " AND "
			}
		}
	} else {
		query = query + ";"
	}
	rows := q.db.QueryRowContext(ctx, query)
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

const infoGetCoreTables = "SELECT * FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"

type InfoGetCoreTable struct {
	SchemaName  string      `json:"schemaname"`
	TableName   string      `json:"tablename"`
	TableOwner  string      `json:"tableowner"`
	TableSpace  interface{} `json:"tablespace"`
	HasIndexes  bool        `json:"hasindexes"`
	HasRules    bool        `json:"hasrules"`
	HasTriggers bool        `json:"hastriggers"`
	RowSecurity bool        `json:"rowsecurity"`
}

func (q *Queries) InfoGetCoreTables(ctx context.Context) ([]InfoGetCoreTable, error) {
	rows, err := q.db.QueryContext(ctx, infoGetCoreTables)
	if err != nil {
		// return nil, err
		println(err)
	}
	defer rows.Close()
	items := []InfoGetCoreTable{}
	for rows.Next() {
		var table InfoGetCoreTable
		err := rows.Scan(
			&table.SchemaName,
			&table.TableName,
			&table.TableOwner,
			&table.TableSpace,
			&table.HasIndexes,
			&table.HasRules,
			&table.HasTriggers,
			&table.RowSecurity,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, table)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, item := range items {
		println(fmt.Sprintf("%s", item.SchemaName))
		println(fmt.Sprintf("%s", item.TableName))
		println(fmt.Sprintf("%s", item.TableOwner))
	}
	return items, nil
}
