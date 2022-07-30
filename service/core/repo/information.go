package core_repo

import (
	"context"
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

const infoGetCoreTables = `-- name: Display ALl Tables
SELECT tablename,schemaname,tableowner 
FROM pg_catalog.pg_tables 
WHERE schemaname != 'pg_catalog' 
AND schemaname != 'information_schema';`

type InfoGetCoreTable struct {
	SchemaName string `json:"schemaname"`
	TableName  string `json:"tablename"`
	TableOwner string `json:"tableowner"`
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
	return items, nil
}

const infoGetTableSchema = `-- name: Get Table Schema
SELECT
    ordinal_position as "position",
    column_name as "column",
    data_type as "type",
    is_nullable as "nullable",
	character_maximum_length as "length",
    numeric_precision as "precision",
    numeric_scale as "scale",
    column_default as "default" 
FROM information_schema.columns
WHERE table_name = $1;`

type InfoColumnSchema struct {
	Position  int32  `json:"position"`
	Column    string `json:"column"`
	Type      string `json:"type"`
	Primary   bool   `json:"primary"`
	Unique    bool   `json:"unique"`
	Nullable  bool   `json:"nullable"`
	Length    int64  `json:"length"`
	Precision int32  `json:"precision"`
	Scale     int32  `json:"scale"`
	Default   string `json:"default"`
}

func (q *Queries) InfoGetTableSchema(ctx context.Context, tablename string) ([]InfoColumnSchema, error) {
	rows, err := q.db.QueryContext(ctx, infoGetTableSchema, tablename)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InfoColumnSchema = []InfoColumnSchema{}

	for rows.Next() {
		var table InfoColumnSchema

		var _nullable interface{}
		var _length interface{}
		var _precision interface{}
		var _scale interface{}
		var _default interface{}

		err := rows.Scan(
			&table.Position,
			&table.Column,
			&table.Type,
			&_nullable,
			&_length,
			&_precision,
			&_scale,
			&_default,
		)
		if err != nil {
			return nil, err
		}

		if _nullable != nil {
			vnullable, isString := _nullable.(string)
			if isString {
				if vnullable == "NO" {
					table.Nullable = false
				} else {
					table.Nullable = true
				}
			} else {
				table.Nullable = true
			}
		}
		if _length != nil {
			val, isInt64 := _length.(int64)
			if isInt64 {
				table.Length = val
			}
		}
		if _precision != nil {
			val, isInt32 := _precision.(int32)
			if isInt32 {
				table.Precision = val
			}
		}
		if _scale != nil {
			val, isInt32 := _scale.(int32)
			if isInt32 {
				table.Precision = val
			}
		}
		if _default != nil {
			val, isString := _default.(string)
			if isString {
				table.Default = val
			} else {
				table.Default = ""
			}
		}
		items = append(items, table)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
