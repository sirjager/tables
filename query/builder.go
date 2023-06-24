package query

type CreateTableQueryBuilder interface {
	AddColumn(c *Column) CreateTableQueryBuilder
	ForeignKey(f *ForeignKey) CreateTableQueryBuilder
	IfNotExists() CreateTableQueryBuilder
	BTreeIndex(idxname, cname string) CreateTableQueryBuilder
	TextSerachIndex(idxname, cname string) CreateTableQueryBuilder
	Query() (string, error)
}

type SelectQueryBuilder interface {
	From(t string) SelectQueryBuilder
	Join(jointype, jointable, matching_condition string) SelectQueryBuilder
	WhereRaw(raw string) SelectQueryBuilder
	Where(expr string, cond string, value interface{}) SelectQueryBuilder
	OrWhere(expr string, cond string, value interface{}) SelectQueryBuilder
	AndWhere(expr string, cond string, value interface{}) SelectQueryBuilder
	WhereNull(expr string) SelectQueryBuilder
	WhereNotNull(expr string) SelectQueryBuilder
	WhereIn(expr string, values []interface{}) SelectQueryBuilder
	WhereNotIn(expr string, values []interface{}) SelectQueryBuilder
	Having(expr string, cond string, value interface{}) SelectQueryBuilder
	OrHaving(expr string, cond string, value interface{}) SelectQueryBuilder
	AndHaving(expr string, cond string, value interface{}) SelectQueryBuilder
	Limit(limit int) SelectQueryBuilder
	Offset(offset int) SelectQueryBuilder
	OrderBy(column string, direction string) SelectQueryBuilder
	GroupBy(columns string) SelectQueryBuilder
	Raw(sql string) SelectQueryBuilder
	Query() (string, error)
}

type InsertQueryBuilder interface {
	Insert([]map[string]interface{}) InsertQueryBuilder
	Where(column string, cond string, value interface{}) InsertQueryBuilder
	NamedExecutables() error
}
