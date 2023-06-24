package query

import (
	"sort"
)

const DATA_TYPES = "integer,smallint,bigint,smallserial,serial,bigserial,decimal,numeric,real,money,uuid" +
	"double precision,json,jsonb,boolean,cidr,inet,macaddr,macaddr8" + "varchar,char,character,text,bit,bit varying" +
	"date,time,timestamp,timestamp with time zone,timestamp without time zone,time with time zone,time without time zone"

func (s *Schema) findRequiredColumns() []string {
	var c []string
	for _, f := range s.Columns {
		if f.NotNull {
			c = append(c, f.Name)
		}
	}
	return c
}

func findSchema(tableName string, schemas []*Schema) *Schema {
	for _, r := range schemas {
		if r.Name == tableName {
			return r
		}
	}
	return &Schema{}
}

func (s *Schema) findNotNullColumnNames() (c []string) {
	for _, t := range s.Columns {
		if t.NotNull {
			c = append(c, t.Name)
		}
	}
	return
}

func (s *Schema) findNotNullColumns() (cols []*Column) {
	for _, t := range s.Columns {
		if t.NotNull {
			cols = append(cols, t)
		}
	}
	return
}

func (s *Schema) findUniqueColumnNames() (cols []string) {
	for _, t := range s.Columns {
		if t.Unique {
			cols = append(cols, t.Name)
		}
	}
	return
}

func (s *Schema) findUniqueColumns() (cols []*Column) {
	for _, t := range s.Columns {
		if t.Unique {
			cols = append(cols, t)
		}
	}
	return
}

func (s *Schema) findPrimaryColumn() *Column {
	for _, t := range s.Columns {
		if t.Primary {
			return t
		}

	}
	return nil
}

func name(c string) string {
	return `"` + c + `"`
}

func (s *Schema) findColumnNames() (c []string) {
	for _, t := range s.Columns {
		c = append(c, t.Name)
	}
	return
}

func (s *Schema) findColumn(name string) (c *Column) {
	for _, t := range s.Columns {
		if t.Name == name {
			return c
		}
	}
	return
}

func getMapKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func sortMapByKeyValue(data map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make(map[string]interface{})
	for _, k := range keys {
		sorted[k] = data[k]
	}
	return sorted
}
