package corm

import (
	"fmt"
	"strings"
)

/**
组合sql语句
*/
func (db *Db) writeBuf(strs ...string) {
	for _, s := range strs {
		db.buffer.WriteString(s)
	}
}

func (db *Db) addSelect() {
	db.writeBuf(SELECT, SPACE)
}

func (db *Db) addUpdate() {
	db.writeBuf(UPDATE, SPACE)
}

func (db *Db) addDelete() {
	db.writeBuf(DELETE, SPACE)
}

func (db *Db) addInsert() {
	db.writeBuf(INSERT, SPACE)
}

func (db *Db) addTable() {
	db.writeBuf(db.table, SPACE)
}

func (db *Db) addFrom() {
	db.writeBuf(FROM, SPACE)
}

func (db *Db) addSet() {
	db.writeBuf(SET, SPACE)
}

func (db *Db) addSum() {
	if db.sum == "" {
		return
	}
	db.writeBuf("SUM(", db.sum, ") AS sum", SPACE)
}

func (db *Db) addMax() {
	if db.max == "" {
		return
	}
	db.writeBuf("MAX(", db.max, ") AS max", SPACE)
}

func (db *Db) addMin() {
	if db.min == "" {
		return
	}
	db.writeBuf("MIN(", db.min, ") AS min", SPACE)
}

func (db *Db) addCount() {
	db.writeBuf("COUNT(*) AS count", SPACE)
}

/**
添加字段
*/
func (db *Db) addFields() {
	var sqlStr string
	if len(db.fields) == 0 {
		db.writeBuf("*", SPACE)
	} else {
		var fieldStr []string
		for _, f := range db.fields {
			fieldStr = append(fieldStr, f)
		}
		sqlStr = strings.Join(fieldStr, COMMA)
	}
	db.writeBuf(sqlStr, SPACE)
}

/**
添加join
*/
func (db *Db) addJoin() {
	if len(db.join) > 0 {
		join := make([]string, 0, 2)
		for _, j := range db.join {
			join = append(join, strings.Join([]string{j.direction, j.table, ON, j.on}, SPACE))
		}
		db.writeBuf(strings.Join(join, SPACE), SPACE)
	}
}

/**
添加where条件
*/
func (db *Db) addWhere() {
	if len(db.whereRaw) > 0 || len(db.where) > 0 {
		sqlTmp := make([]string, 0, 5)
		if len(db.where) > 0 {
			for _, w := range db.where {
				switch w.operator {
				case IN, NOT_IN:
					sqlTmp = append(sqlTmp, w.field+SPACE+w.operator+"("+arrayToStrPlace(w.conditionArray)+")")
				case LIKE, NOT_LIKE:
					sqlTmp = append(sqlTmp, w.field+SPACE+w.operator+SPACE+QUES)
				case BETWEEN:
					sqlTmp = append(sqlTmp, w.field+SPACE+w.operator+SPACE+QUES+SPACE+AND+SPACE+QUES)
				default:
					sqlTmp = append(sqlTmp, w.field+SPACE+w.operator+SPACE+QUES)
				}
			}
		}
		if len(db.whereRaw) > 0 {
			for _, w := range db.whereRaw {
				sqlTmp = append(sqlTmp, w)
			}
		}
		db.writeBuf(WHERE, SPACE, strings.Join(sqlTmp, " AND "), SPACE)
	}
}

/**
添加排序条件
*/
func (db *Db) addOrderBy() {
	if len(db.orderBy) > 0 {
		order := make([]string, 0, 2)
		for _, o := range db.orderBy {
			order = append(order, o.field+SPACE+o.by)
		}
		db.writeBuf(ORDER_BY, SPACE, strings.Join(order, COMMA), SPACE)
	}
}

/**
添加limit
*/
func (db *Db) addLimit() {
	if db.limit > 0 {
		db.writeBuf(LIMIT, SPACE, string(db.limit))
	}
}

/**
条件转SQL语句
*/
func (db *Db) whereToSql() string {
	db.check()
	db.addSelect()
	db.addFields()
	db.addFrom()
	db.addTable()
	db.addJoin()
	db.addWhere()
	db.addOrderBy()
	db.addLimit()
	return db.buffer.String()
}

func (db *Db) countToSql() string {
	db.check()
	db.addSelect()
	db.addCount()
	db.addFrom()
	db.addTable()
	db.addJoin()
	db.addWhere()
	return db.buffer.String()
}

func (db *Db) sumToSql() string {
	db.check()
	db.addSelect()
	db.addSum()
	db.addFrom()
	db.addTable()
	db.addJoin()
	db.addWhere()
	return db.buffer.String()
}

func (db *Db) maxToSql() string {
	db.check()
	db.addSelect()
	db.addMax()
	db.addFrom()
	db.addTable()
	db.addJoin()
	db.addWhere()
	db.addOrderBy()
	db.addLimit()
	return db.buffer.String()
}

func (db *Db) minToSql() string {
	db.check()
	db.addSelect()
	db.addMin()
	db.addFrom()
	db.addTable()
	db.addJoin()
	db.addWhere()
	db.addOrderBy()
	db.addLimit()
	return db.buffer.String()
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

	keysToStr := db.table + "(" + strings.Join(keys, ", ") + ")"
	keyValsToStr := " VALUES(" + strings.Join(keyVals, ", ") + ")"
	insertStr := keysToStr + keyValsToStr

	return insertStr, vals
}

func (db *Db) updateToStrAndArr() (string, []interface{}) {
	var keys []string
	var vals []interface{}

	for k, v := range db.update {
		keys = append(keys, k+"=?")
		vals = append(vals, v)
	}

	return strings.Join(keys, ", "), vals
}

func (db *Db) insertToSql() (sql string, arr []interface{}) {
	db.check()
	insertStr, vals := db.insertToStrAndArr()

	db.addInsert()
	db.writeBuf(insertStr, SPACE)

	return db.buffer.String(), vals
}

func (db *Db) updateToSql() (sql string, arr []interface{}) {
	db.check()
	updateStr, vals := db.updateToStrAndArr()

	db.addUpdate()
	db.addTable()
	db.addSet()
	db.writeBuf(updateStr, SPACE)
	db.addWhere()
	return retSql(db.buffer.String()), vals
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
