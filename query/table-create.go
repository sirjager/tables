package query

import (
	"fmt"
	"strings"

	"github.com/sirjager/gotables/pkg/utils"
	"github.com/sirjager/gotables/pkg/validator"
)

func SafeName(db, str string) string {
	switch strings.ToLower(db) {
	case "pg", "postgres", "postgresql":
		return `"` + str + `"`
	case "mysql", "mysqldb", "maria", "mariadb":
		return "`" + str + "`"
	case "sqlite", "sqlite3":
		return "`" + str + "`"
	default:
		return str
	}
}

type createTable struct {
	// optional schema of all tables
	// if provided then will be used to validate queries by querybuilder
	schemas []Schema
	// postgres, mysql, mariadb, sqlite, cockroachdb

	tableName     string
	ifNotExists   bool
	rawQuery      string
	columnNames   []string
	columnStrings []string
	foreginkeys   []string
	violations    []string
	indexes       []string
}

func CreateTable(name string) CreateTableQueryBuilder {
	return &createTable{tableName: name}
}

func (t *createTable) AddColumn(c *Column) CreateTableQueryBuilder {
	if c.Raw != "" {
		cname := strings.Split(c.Raw, " ")[0]
		if err := validator.ValidateColumnName(cname); err != nil {
			t.violations = append(t.violations, err.Error())
			return t
		}
		// find duplicate columns
		for _, rr := range t.columnNames {
			if rr == cname {
				t.violations = append(t.violations, c.Name+" multiple times found. Keep single instance")
				return t
			}
		}

		t.columnNames = append(t.columnNames, cname)
		t.columnStrings = append(t.columnStrings, c.Raw)
		return t
	}

	// find duplicate columns
	for _, rr := range t.columnNames {
		if rr == c.Name {
			t.violations = append(t.violations, c.Name+" multiple times found. Keep single instance")
			return t
		}
	}

	if err := validator.ValidateColumnName(c.Name); err != nil {
		t.violations = append(t.violations, err.Error())
		return t
	}

	t.columnNames = append(t.columnNames, c.Name)

	str := name(c.Name)

	column_type := ""
	upper_col_type := strings.ToUpper(c.Type)
	lower_col_type := strings.ToLower(c.Type)

	switch lower_col_type {
	// First three can be simplified to one case but it will become too lenthy
	// so spliting into three for simplicity
	case "integer", "smallint", "bigint", "smallserial", "serial", "bigserial", "real", "money", "uuid":
		column_type = upper_col_type
	case "double precision", "json", "jsonb", "boolean", "cidr", "inet", "macaddr", "macaddr8":
		column_type = upper_col_type
	case "date", "time", "timestamp", "timestamp with time zone", "timestamp without time zone", "time with time zone", "time without time zone":
		column_type = upper_col_type
	case "decimal", "numeric":
		if c.Precision > 0 {
			if c.Scale > 0 {
				column_type = fmt.Sprintf("%s(%d,%d)", upper_col_type, c.Precision, c.Scale)
				break
			} else {
				column_type = fmt.Sprintf("%s(%d)", upper_col_type, c.Precision)
				break
			}
		} else {
			column_type = upper_col_type
			break
		}
	case "text":
		column_type = upper_col_type
	case "varchar", "char", "character", "bit", "bit varying":
		if c.Length > 0 {
			column_type = fmt.Sprintf("%s(%d)", upper_col_type, c.Length)
		} else {
			column_type = upper_col_type
		}
	default:
		t.violations = append(t.violations, fmt.Sprintf("column=(%s) contains invalid type=(%s)", c.Name, c.Type))
		return t
	}

	if len(column_type) < 1 {
		t.violations = append(t.violations, fmt.Sprintf("column=(%s) contains invalid type=(%s)", c.Name, c.Type))
		return t
	}

	str += " " + column_type
	if c.Primary {
		str += " PRIMARY KEY"
	}
	if c.Unique {
		str += " UNIQUE"
	}
	if c.NotNull {
		str += " NOT NULL"
	}

	if c.Default != "" {
		str += " DEFAULT" + "(" + c.Default + ")"
	}

	t.columnStrings = append(t.columnStrings, str)
	return t
}

func (t *createTable) ForeignKey(f *ForeignKey) CreateTableQueryBuilder {

	if !utils.ValueExist(f.Column, t.columnNames) {
		t.violations = append(t.violations, fmt.Sprintf("column %s in foregin key does not exists", f.Column))
		return t
	}
	if err := validator.ValidateTableName(f.RefTable); err != nil {
		t.violations = append(t.violations, fmt.Sprintf("invalid foreign key reference table name: %s", err.Error()))
		return t
	}

	if err := validator.ValidateColumnName(f.RefColumn); err != nil {
		t.violations = append(t.violations, fmt.Sprintf("invalid foreign key reference column name: %s", err.Error()))
		return t
	}

	// if schema is provided then check if ref table has this column or not
	if len(t.schemas) > 0 {
		for _, s := range t.schemas {
			if s.Name == f.RefTable {
				refColExists := false
				// loop over columns and check if column is present or not
				for _, sc := range s.Columns {
					if sc.Name == f.RefColumn {
						refColExists = true
						break
					}
				}
				if !refColExists {
					t.violations = append(t.violations, fmt.Sprintf("reference column '%s' does not exists in reference table '%s'", f.RefColumn, f.RefTable))
					return t
				}

				// if table found then no break the loop
				break
			}
		}
	}

	fk := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)", name(f.Column), name(f.RefTable), name(f.RefColumn))

	ondel := ""
	if f.OnDelCascade {
		ondel = " ON DELETE CASCADE"
	} else if f.OnDelSetNull {
		ondel = " ON DELETE SET NULL"
	} else if f.OnDelSetValue != "" {
		ondel = " ON DELETE SET " + f.OnDelSetValue
	}
	fk += ondel

	// check for duplicates
	for _, r := range t.foreginkeys {
		if r == fk {
			t.violations = append(t.violations, fmt.Sprintf("duplicate foregin key: %s", fk))
			return t
		}
	}

	t.foreginkeys = append(t.foreginkeys, fk)
	return t
}

func (t *createTable) IfNotExists() CreateTableQueryBuilder {
	t.ifNotExists = true
	return t
}

func (t *createTable) BTreeIndex(idx, cname string) CreateTableQueryBuilder {
	// check if column exists
	if !utils.ValueExist(cname, t.columnNames) {
		t.violations = append(t.violations, fmt.Sprintf("failed to create index, '%s' column does not exists", cname))
		return t
	}

	// validate index name
	if err := validator.ValidateIsAlphaNumUnderscore(idx); err != nil {
		t.violations = append(t.violations, "invalid index name: "+idx+" "+err.Error())
		return t
	}

	index := fmt.Sprintf("CREATE INDEX %s ON %s (%s)", idx, name(t.tableName), name(cname))

	// check if exact same index exists
	if utils.ValueExist(index, t.indexes) {
		t.violations = append(t.violations, "duplicate index: "+index)
		return t
	}
	t.indexes = append(t.indexes, index)
	return t
}

func (t *createTable) TextSerachIndex(idx, cname string) CreateTableQueryBuilder {
	// check if column exists
	if !utils.ValueExist(cname, t.columnNames) {
		t.violations = append(t.violations, fmt.Sprintf("failed to create index, '%s' column does not exists", cname))
		return t
	}
	// validate index name
	if err := validator.ValidateIsAlphaNumUnderscore(idx); err != nil {
		t.violations = append(t.violations, "invalid index name: '"+idx+"' "+err.Error())
		return t
	}

	// CREATE INDEX idx_text_search_column_name ON table_name USING gin(to_tsvector('english', column_name));
	index := fmt.Sprintf("CREATE INDEX %s ON %s USING gin(to_tsvector('english', %s))", idx, name(t.tableName), name(cname))

	// check if exact same index exists
	if utils.ValueExist(index, t.indexes) {
		t.violations = append(t.violations, "duplicate index: "+index)
		return t
	}
	t.indexes = append(t.indexes, index)
	return t
}

func (t *createTable) Query() (string, error) {
	if len(t.violations) > 0 {
		all := strings.Join(t.violations, ", ")
		return "", fmt.Errorf("%d errors: %s", len(t.violations), all)
	}
	if t.rawQuery != "" {
		return t.rawQuery, nil
	}
	str := "CREATE TABLE " + name(t.tableName)
	str += " (\n"
	str += strings.Join(t.columnStrings, ",\n")
	if len(t.foreginkeys) > 0 {
		str += ",\n"
		str += strings.Join(t.foreginkeys, ",\n")
	}
	str += "\n);\n"

	if len(t.indexes) > 0 {
		str += "\n"
		allindex := strings.Join(t.indexes, ";\n")
		str += allindex
		str += ";\n"
	}

	return str, nil
}
