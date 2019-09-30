package fake_orm

import (
	"fmt"
	"strings"
)

/**
组合sql语句
*/
func sqlCompose(sql ...string) string {
	return strings.Join(sql, " ")
}

func (db *db) addSelect() string {
	return SELECT
}

func (db *db) addUpdate() string {
	return UPDATE
}

func (db *db) addDelete() string {
	return DELETE
}

/**
添加字段
*/
func (db *db) addFields() (sqlStr string) {
	if len(db.fields) == 0 {
		sqlStr = "*"
	} else {
		var fieldStr []string
		for _, f := range db.fields {
			fieldStr = append(fieldStr, f)
		}
		sqlStr = strings.Join(fieldStr, ",")
	}
	return
}

/**
添加join
*/
func (db *db) addJoin() (sqlStr string) {
	if len(db.join) > 0 {
		join := make([]string, 0, 2)
		for _, j := range db.join {
			join = append(join, fmt.Sprintf("%s %s %s %s", j.direction, j.table, ON, j.on))
		}
		sqlStr = strings.Join(join, " ")
	}
	return
}

/**
添加where条件
*/
func (db *db) addWhere() string {
	if len(db.whereRaw) > 0 || len(db.where) > 0 {
		sqlTmp := make([]string, 0, 5)
		if len(db.where) > 0 {
			for _, w := range db.where {
				switch w.operator {
				case IN, NOT_IN:
					sqlTmp = append(sqlTmp, fmt.Sprintf("%s %s(%s)", w.field, w.operator, arrayToStrPlace(w.conditionArray)))
				case LIKE, NOT_LIKE:
					sqlTmp = append(sqlTmp, fmt.Sprintf("%s %s '%%s%'", w.field, w.operator, w.condition))
				case BETWEEN:
					sqlTmp = append(sqlTmp, fmt.Sprintf("%s %s %s AND %s", w.field, w.operator, "?", "?"))
				default:
					sqlTmp = append(sqlTmp, fmt.Sprintf("%s %s %s", w.field, w.operator, "?"))
				}
			}
		}
		if len(db.whereRaw) > 0 {
			for _, w := range db.whereRaw {
				sqlTmp = append(sqlTmp, w)
			}
		}

		return fmt.Sprintf("%s %s", WHERE, strings.Join(sqlTmp, fmt.Sprintf(" %s ", AND)))
	}
	return ""
}

/**
添加排序条件
*/
func (db *db) addOrderBy() string {
	if len(db.orderBy) > 0 {
		order := make([]string, 0, 2)
		for _, o := range db.orderBy {
			order = append(order, fmt.Sprintf("%s %s", o.field, o.by))
		}
		return fmt.Sprintf("%s %s", ORDER_BY, strings.Join(order, ","))
	}
	return ""
}

/**
添加limit
*/
func (db *db) addLimit() string {
	if db.limit > 0 {
		return fmt.Sprintf("%s %d", LIMIT, db.limit)
	}
	return ""
}

/**
添加数据表
*/
func (db *db) addTable() string {
	return db.table
}

/**
添加form
*/
func (db *db) addFrom() string {
	return FROM
}

func (db *db) addSum() string {
	if db.sum == "" {
		return ""
	}
	return fmt.Sprintf("SUM(%s) AS sum", db.sum)
}

func (db *db) addMax() string {
	if db.max == "" {
		return ""
	}
	return fmt.Sprintf("MAX(%s) AS max", db.max)
}

func (db *db) addMin() string {
	if db.min == "" {
		return ""
	}
	return fmt.Sprintf("MIN(%s) AS min", db.min)
}

func (db *db) addCount() string {
	return "COUNT(*) AS count"
}

/**
条件转SQL语句
*/
func (db *db) whereToSql() string {
	sqlStr := sqlCompose(
		db.addSelect(),
		db.addFields(),
		db.addFrom(),
		db.addTable(),
		db.addJoin(),
		db.addWhere(),
		db.addOrderBy(),
		db.addLimit(),
	)
	return retSql(sqlStr)
}

func (db *db) countToSql() string {
	sqlStr := sqlCompose(
		db.addSelect(),
		db.addCount(),
		db.addFrom(),
		db.addTable(),
		db.addJoin(),
		db.addWhere(),
		db.addOrderBy(),
		db.addLimit(),
	)
	return retSql(sqlStr)
}

func (db *db) sumToSql() string {
	sqlStr := sqlCompose(
		db.addSelect(),
		db.addSum(),
		db.addFrom(),
		db.addTable(),
		db.addJoin(),
		db.addWhere(),
		db.addOrderBy(),
		db.addLimit(),
	)
	return retSql(sqlStr)
}

func (db *db) maxToSql() string {
	sqlStr := sqlCompose(
		db.addSelect(),
		db.addMax(),
		db.addFrom(),
		db.addTable(),
		db.addJoin(),
		db.addWhere(),
		db.addOrderBy(),
		db.addLimit(),
	)
	return retSql(sqlStr)
}

func (db *db) minToSql() string {
	sqlStr := sqlCompose(
		db.addSelect(),
		db.addMin(),
		db.addFrom(),
		db.addTable(),
		db.addJoin(),
		db.addWhere(),
		db.addOrderBy(),
		db.addLimit(),
	)
	return retSql(sqlStr)
}

func retSql(sqlStr string) string {
	fmt.Println(sqlStr)
	return sqlStr
}

/**
条件转SQL语句, 自定义LIMIT
*/
func (db *db) whereToSqlForLimit(count int64) string {
	db.limit = count
	return db.whereToSql()
}

func arrayToStrPlace(arr []interface{}) string {
	strTmp := make([]string, 0, 5)
	for i := 0; i < len(arr); {
		strTmp = append(strTmp, "?")
	}
	return strings.Join(strTmp, ",")
}
