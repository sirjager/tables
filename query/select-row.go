package query

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirjager/gotables/pkg/utils"
	"github.com/sirjager/gotables/pkg/validator"
)

type selectQueryBuilder struct {
	columns          []string
	from             string
	joins            []string
	where            []string
	limit            int
	offset           int
	orderByColumn    string
	orderByDirection string
	groupBy          string
	unions           []string
	having           []string
	rawSql           string
	violations       []string
}

func Select(comma_seperated_columns_or_aggregate_functions string) SelectQueryBuilder {
	return &selectQueryBuilder{columns: strings.Split(comma_seperated_columns_or_aggregate_functions, ",")}
}

func (qb *selectQueryBuilder) From(t string) SelectQueryBuilder {
	if err := validator.ValidateTableName(t); err != nil {
		qb.violations = append(qb.violations, err.Error())
		return qb
	}
	qb.from = t
	return qb
}

const all_joins = "inner,left,right,full outer,self,cross"

func (qb *selectQueryBuilder) Join(join, table, cond string) SelectQueryBuilder {
	_joinType := strings.ToLower(join)
	valid := false
	for _, v := range strings.Split(all_joins, ",") {
		if v == _joinType {
			valid = true
		}
	}
	if !valid {
		qb.violations = append(qb.violations, fmt.Sprintf("invalid join : %s", join))
		return qb
	}
	_joinStatement := fmt.Sprintf(" %s JOIN %s ON %s", strings.ToUpper(join), table, cond)
	if utils.ValueExist(_joinStatement, qb.joins) {
		qb.violations = append(qb.violations, fmt.Sprintf("duplicate join statement: %s", _joinStatement))
		return qb
	}
	qb.joins = append(qb.joins, _joinStatement)
	return qb
}

func (qb *selectQueryBuilder) WhereRaw(raw string) SelectQueryBuilder {
	if len(qb.where) == 0 {
		qb.where = append(qb.where, " "+raw)
	} else {
		qb.where = append(qb.where, fmt.Sprintf(" AND %s", raw))
	}
	return qb
}

func (qb *selectQueryBuilder) Where(expr string, cond string, value interface{}) SelectQueryBuilder {
	_value, err := toPgString(value)
	if err != nil {
		qb.violations = append(qb.violations, err.Error())
		return qb
	}
	if len(qb.where) == 0 {
		qb.where = append(qb.where, " "+expr+" "+cond+" "+" "+_value)
	} else {
		qb.where = append(qb.where, " AND "+expr+" "+cond+" "+" "+_value)
	}
	return qb
}

func (qb *selectQueryBuilder) OrWhere(expr string, cond string, value interface{}) SelectQueryBuilder {
	_value, err := toPgString(value)
	if err != nil {
		qb.violations = append(qb.violations, err.Error())
		return qb
	}
	if len(qb.where) == 0 {
		qb.where = append(qb.where, " "+expr+" "+cond+" "+" "+_value)
	} else {
		qb.where = append(qb.where, " OR "+expr+" "+cond+" "+" "+_value)
	}
	return qb
}

func (qb *selectQueryBuilder) AndWhere(expr string, cond string, value interface{}) SelectQueryBuilder {
	qb.Where(expr, cond, value)
	return qb
}

func (qb *selectQueryBuilder) WhereNull(column string) SelectQueryBuilder {
	if len(qb.where) == 0 {
		qb.where = append(qb.where, fmt.Sprintf(" %s IS NULL ", column))
	} else {
		qb.where = append(qb.where, fmt.Sprintf(" AND %s IS NULL ", column))
	}
	return qb
}

func (qb *selectQueryBuilder) WhereNotNull(column string) SelectQueryBuilder {
	if len(qb.where) == 0 {
		qb.where = append(qb.where, fmt.Sprintf(" %s IS NOT NULL ", column))
	} else {
		qb.where = append(qb.where, fmt.Sprintf(" AND %s IS NOT NULL ", column))
	}
	return qb
}

func (qb *selectQueryBuilder) WhereIn(column string, values []interface{}) SelectQueryBuilder {
	allValues := []string{}
	for _, v := range values {
		_value, err := toPgString(v)
		if err != nil {
			qb.violations = append(qb.violations, err.Error())
			return qb
		}
		allValues = append(allValues, _value)
	}
	list := strings.Join(allValues, ",")
	if len(qb.where) == 0 {
		qb.where = append(qb.where, fmt.Sprintf(" %s IN (%s) ", column, list))
	} else {
		qb.where = append(qb.where, fmt.Sprintf(" AND %s IN (%s) ", column, list))
	}
	return qb
}

func (qb *selectQueryBuilder) WhereNotIn(column string, values []interface{}) SelectQueryBuilder {
	allValues := []string{}
	for _, v := range values {
		_value, err := toPgString(v)
		if err != nil {
			qb.violations = append(qb.violations, err.Error())
			return qb
		}
		allValues = append(allValues, _value)
	}
	list := strings.Join(allValues, ",")
	if len(qb.where) == 0 {
		qb.where = append(qb.where, fmt.Sprintf(" %s NOT IN (%s) ", column, list))
	} else {
		qb.where = append(qb.where, fmt.Sprintf(" AND %s NOT IN (%s) ", column, list))
	}
	return qb
}

func (qb *selectQueryBuilder) Having(expr string, cond string, value interface{}) SelectQueryBuilder {
	_value, err := toPgString(value)
	if err != nil {
		qb.violations = append(qb.violations, err.Error())
		return qb
	}

	if len(qb.having) == 0 {
		qb.having = append(qb.having, " "+expr+" "+cond+" "+_value+" ")
	} else {
		qb.having = append(qb.having, " AND "+expr+" "+cond+" "+_value+" ")
	}
	return qb
}

func (qb *selectQueryBuilder) OrHaving(expr string, cond string, value interface{}) SelectQueryBuilder {
	_value, err := toPgString(value)
	if err != nil {
		qb.violations = append(qb.violations, err.Error())
	}
	if len(qb.having) == 0 {
		qb.having = append(qb.having, " "+expr+" "+cond+" "+_value+" ")
	} else {
		qb.having = append(qb.having, " OR "+expr+" "+cond+" "+_value+" ")
	}
	return qb
}

func (qb *selectQueryBuilder) AndHaving(expr string, cond string, value interface{}) SelectQueryBuilder {
	qb.Having(expr, cond, value)
	return qb
}

func (qb *selectQueryBuilder) Limit(limit int) SelectQueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *selectQueryBuilder) Offset(offset int) SelectQueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *selectQueryBuilder) OrderBy(column string, direction string) SelectQueryBuilder {
	qb.orderByColumn = column
	qb.orderByDirection = direction
	return qb
}

func (qb *selectQueryBuilder) GroupBy(column string) SelectQueryBuilder {
	qb.groupBy = column
	return qb
}

func (qb *selectQueryBuilder) Raw(sql string) SelectQueryBuilder {
	qb.rawSql = sql
	return qb
}

func (qb *selectQueryBuilder) Query() (string, error) {
	if qb.rawSql != "" {
		return qb.rawSql, nil
	}

	// see if any violations
	if len(qb.violations) > 0 {
		all := strings.Join(qb.violations, ", ")
		return "", fmt.Errorf("%d errors: %s", len(qb.violations), all)
	}

	// select col1, col2, col3
	query := "SELECT "
	query += strings.Join(qb.columns, ",")

	// from table
	query += " FROM " + qb.from

	// joins
	if len(qb.joins) > 0 {
		joinstr := strings.Join(qb.joins, " ")
		query += " " + joinstr
	}

	// where
	if len(qb.where) > 0 {
		query += " WHERE" + strings.Join(qb.where, "")
	}

	// group
	if qb.groupBy != "" {
		query += " GROUP BY " + qb.groupBy
	}

	// having
	if len(qb.having) > 0 {
		query += " HAVING" + strings.Join(qb.having, "")
	}

	// unions
	if len(qb.unions) > 0 {
		query = " " + strings.Join(qb.unions, " ") + " " + query
	}

	// order by
	if qb.orderByColumn != "" {
		query += " ORDER BY " + qb.orderByColumn + " " + qb.orderByDirection
	}

	// limit
	if qb.limit > 0 {
		query += " LIMIT " + strconv.Itoa(qb.limit)
	}

	//offset
	if qb.offset > 0 {
		query += " OFFSET " + strconv.Itoa(qb.offset)
	}

	// end with semicolor;
	query += ";"

	return query, nil
}

func toPgString(val interface{}) (string, error) {
	_type := reflect.TypeOf(val)
	switch _type.Name() {
	case "string":
		return fmt.Sprintf("'%s'", val), nil
	case "bool", "boolean":
		return fmt.Sprintf("%v", val), nil
	case "int":
		return strconv.Itoa(val.(int)), nil
	case "int64":
		return strconv.FormatInt(val.(int64), 10), nil
	case "uint64":
		return strconv.FormatUint(val.(uint64), 10), nil
	case "float64":
		return fmt.Sprintf("%g", val), nil
	}
	return "", fmt.Errorf("invalid value: %s", _type.Elem())
}
