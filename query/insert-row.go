package query

import (
	"fmt"
	"sort"

	"strings"

	"github.com/sirjager/gotables/pkg/utils"
)

type insertQueryBuilder struct {
	tableName string
	schema    *Schema
	schemas   []*Schema
	where     []string
	err       error
	exec      []*ExecPack
}

func InsertInto(tableName string, schemas ...*Schema) InsertQueryBuilder {
	t := &insertQueryBuilder{tableName: tableName, schemas: schemas, schema: &Schema{}}
	if len(schemas) > 0 {
		for _, s := range schemas {
			if tableName == s.Name {
				t.schema = s
				break
			}
		}
	}
	return t
}

func (qb *insertQueryBuilder) Where(expr string, cond string, value interface{}) InsertQueryBuilder {
	_value, err := toPgString(value)
	if err != nil {
		qb.err = err
		return qb
	}
	if len(qb.where) == 0 {
		qb.where = append(qb.where, " "+expr+" "+cond+" "+" "+_value)
	} else {
		qb.where = append(qb.where, " AND "+expr+" "+cond+" "+" "+_value)
	}
	return qb
}

func (i *insertQueryBuilder) NamedExecutables() ([]*ExecPack, error) {
	if i.err != nil {
		return nil, i.err
	}

	return i.exec, nil
}

// Insert Rows
//
// example data:
// unfiltered json data, can have missing/extra columns, can have null values,
// totally unreliable data
//
//	[
//	 {
//	  "userId": yfOhCUHXblDbTWhX,
//	  "id": 1,
//	  "published": false,
//	  "title": "title 1",
//	  "body": "body 1"
//	 },
//	 {
//	  "userId": ewd2dzoGAlAmfhml,
//	  "id": 2,
//	  "published": true,
//	  "title": "title 2",
//	  "body": "body 2"
//	 },
//	 {
//	  "extraCol": 546.22,
//	  "title": "title 3",
//	  "body": null
//	 },
//	]
//
// generated query:
//
// INSERT INTO "table" ("userId","id","published","title","body") VALUES
// ('yfOhCUHXblDbTWhX',1,false,'title 1', 'body 1'),
// ('ewd2dzoGAlAmfhml',2,true,'title 2', 'body 2');
// INSERT INTO "table" ("title","body") VALUES ('title 3', 'body 1');
//
// Steps:
//  1. rows which have same columns will be packed together in one insert statement
//  2. rows which are not packed or left over will have seperate insert statement
//  3. rows will be packed in seperate insert statements because
//     we are assuming some columns are missing, if all rows are legit,
//     then they will be packed in one insert statement
//  4. extra columns which are not in table schema will be omitted from generated query, though error will be sent back
func (i *insertQueryBuilder) Insert(jsonStruct []map[string]interface{}) InsertQueryBuilder {
	if i.schema.Name != i.tableName {
		return i
	}

	//
	// [id title body userId published]
	columnNames := sort.StringSlice(i.schema.findColumnNames())

	// [extraCol]
	rowsPack := map[string][]map[string]interface{}{}

	// columnPacks:= [[title body userId id published] [userId id published title body] [title body]]
	// this loop will do step 1,2
	for _, row := range jsonStruct {
		rowColumns := []string{}
		for k := range row {
			// checking if column exists or not
			if !utils.ValueExist(k, columnNames) {
				i.err = fmt.Errorf("column '%s' does not exists", k)
				return i
			}
			rowColumns = append(rowColumns, k)
		}
		sort.Strings(rowColumns)
		key := strings.Join(rowColumns, "-")
		rowsPack[key] = append(rowsPack[key], row)
	}

	// now we will do validations and create insert statements

	for _cnames, rows := range rowsPack {
		query, err := i.schema.processRows(_cnames, rows)
		if err != nil {
			i.err = err
			return i
		}
		i.exec = append(i.exec, &ExecPack{Query: query, Data: rows})
	}

	return i
}

type ExecPack struct {
	Query string
	Data  []map[string]interface{}
}

func (s *Schema) processRows(cnames string, rows []map[string]interface{}) (str string, err error) {
	if len(rows) == 0 {
		return "", fmt.Errorf("no rows to process")
	}

	// This is how statement will look
	// INSERT INTO "table" ("userId","id","published","title","body") VALUES
	// ('yfOhCUHXblDbTWhX',1,false,'title 1', 'body 1'),
	// ('ewd2dzoGAlAmfhml',2,true,'title 2', 'body 2');

	keys := getMapKeys(rows[0])
	insertInColumns := []string{}
	for _, r := range keys {
		insertInColumns = append(insertInColumns, ":"+r)
	}

	str = "INSERT INTO " + name(s.Name)
	str += fmt.Sprintf(" (%s) \n ", strings.Join(insertInColumns, ","))
	str += fmt.Sprintf(" VALUES (%s) ", strings.Join(insertInColumns, ", "))
	return
}

func getValue(s *Schema, i int, k string, v interface{}) (str string, err error) {
	c := s.findColumn(k)
	if c != nil {
		return "", fmt.Errorf("column '%s' does not exists", k)
	}
	// values := ""

	switch c.Type {
	case "text", "varchar":
		if c.Type == "varchar" {
			if c.Length > 0 {
				if v != nil && len(fmt.Sprintf("%v", v)) > int(c.Length) {
					if v != nil && len(fmt.Sprintf("%v", v)) > int(c.Length) {
						return "", fmt.Errorf("row [%d] column [%v] can not have character more than length [%d]", i+1, k, c.Length)
					}
				}
			}
		}
		return fmt.Sprintf("'%s'", v), nil
	}

	return
}
