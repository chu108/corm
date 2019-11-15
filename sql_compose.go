package corm

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

func (db *Db) addSelect() string {
	return SELECT
}

func (db *Db) addUpdate() string {
	return UPDATE
}

func (db *Db) addDelete() string {
	return DELETE
}

func (db *Db) addInsert() string {
	return INSERT
}

func (db *Db) addTable() string {
	return db.table
}

func (db *Db) addFrom() string {
	return FROM
}

func (db *Db) addSet() string {
	return SET
}

func (db *Db) addSum() string {
	if db.sum == "" {
		return ""
	}
	return fmt.Sprintf("SUM(%s) AS sum", db.sum)
}

func (db *Db) addMax() string {
	if db.max == "" {
		return ""
	}
	return fmt.Sprintf("MAX(%s) AS max", db.max)
}

func (db *Db) addMin() string {
	if db.min == "" {
		return ""
	}
	return fmt.Sprintf("MIN(%s) AS min", db.min)
}

func (db *Db) addCount() string {
	return "COUNT(*) AS count"
}

/**
添加字段
*/
func (db *Db) addFields() (sqlStr string) {
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
func (db *Db) addJoin() (sqlStr string) {
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
func (db *Db) addWhere() string {
	if len(db.whereRaw) > 0 || len(db.where) > 0 {
		sqlTmp := make([]string, 0, 5)
		if len(db.where) > 0 {
			for _, w := range db.where {
				switch w.operator {
				case IN, NOT_IN:
					sqlTmp = append(sqlTmp, fmt.Sprintf("%s %s(%s)", w.field, w.operator, arrayToStrPlace(w.conditionArray)))
				case LIKE, NOT_LIKE:
					sqlTmp = append(sqlTmp, fmt.Sprintf("%s %s %s", w.field, w.operator, "?"))
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
func (db *Db) addOrderBy() string {
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
func (db *Db) addLimit() string {
	if db.limit > 0 {
		return fmt.Sprintf("%s %d", LIMIT, db.limit)
	}
	return ""
}

/**
条件转SQL语句
*/
func (db *Db) whereToSql() string {
	db.check()
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

func (db *Db) countToSql() string {
	db.check()
	sqlStr := sqlCompose(
		db.addSelect(),
		db.addCount(),
		db.addFrom(),
		db.addTable(),
		db.addJoin(),
		db.addWhere(),
	)
	return retSql(sqlStr)
}

func (db *Db) sumToSql() string {
	db.check()
	sqlStr := sqlCompose(
		db.addSelect(),
		db.addSum(),
		db.addFrom(),
		db.addTable(),
		db.addJoin(),
		db.addWhere(),
	)
	return retSql(sqlStr)
}

func (db *Db) maxToSql() string {
	db.check()
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

func (db *Db) minToSql() string {
	db.check()
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

func (db *Db) insertToStrAndArr() (string, []interface{}) {
	var keys []string
	var keyVals []string
	var vals []interface{}

	for k, v := range db.insert {
		keys = append(keys, k)
		keyVals = append(keyVals, "?")
		vals = append(vals, v)
	}

	keysToStr := db.addTable() + "(" + strings.Join(keys, ", ") + ")"
	keyValsToStr := " VALUES(" + strings.Join(keyVals, ", ") + ")"
	insertStr := keysToStr + keyValsToStr

	return insertStr, vals
}

func (db *Db) updateToStrAndArr() (string, []interface{}) {
	var keys []string
	var vals []interface{}

	for k, v := range db.update {
		keys = append(keys, fmt.Sprintf("%s = ?", k))
		vals = append(vals, v)
	}

	return strings.Join(keys, ", "), vals
}

func (db *Db) insertToSql() (sql string, arr []interface{}) {
	db.check()
	insertStr, vals := db.insertToStrAndArr()

	sqlStr := sqlCompose(
		db.addInsert(),
		insertStr,
	)
	return retSql(sqlStr), vals
}

func (db *Db) updateToSql() (sql string, arr []interface{}) {
	db.check()
	updateStr, vals := db.updateToStrAndArr()

	sqlStr := sqlCompose(
		db.addUpdate(),
		db.addTable(),
		db.addSet(),
		updateStr,
		db.addWhere(),
	)
	return retSql(sqlStr), vals
}

func retSql(sqlStr string) string {
	fmt.Println(sqlStr)
	return sqlStr
}

/**
条件转SQL语句, 自定义LIMIT
*/
func (db *Db) whereToSqlForLimit(count int) string {
	db.limit = count
	return db.whereToSql()
}

func arrayToStrPlace(arr []interface{}) string {
	strTmp := make([]string, 0, 5)
	for i := 0; i < len(arr); i++ {
		strTmp = append(strTmp, "?")
	}
	return strings.Join(strTmp, ",")
}

func (db *Db) getWhereValue() []interface{} {
	where := make([]interface{}, 0, 5)
	if len(db.where) > 0 {
		for _, w := range db.where {
			switch w.operator {
			case IN, NOT_IN:
				where = append(where, w.conditionArray...)
			case LIKE, NOT_LIKE:
				where = append(where, w.condition)
			case BETWEEN:
				where = append(where, w.conditionArray...)
			default:
				where = append(where, w.condition)
			}
		}
	}
	return where
}
