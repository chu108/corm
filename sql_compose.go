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
	return "SELECT"
}

func (db *db) addUpdate() string {
	return "UPDATE"
}

func (db *db) addDelete() string {
	return "DELETE"
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
			join = append(join, fmt.Sprintf("%s %s ON %s", j.direction, j.table, j.on))
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
				sqlTmp = append(sqlTmp, fmt.Sprintf("%s %s %s", w.field, w.operator, "?"))
			}
		}
		if len(db.whereRaw) > 0 {
			for _, w := range db.whereRaw {
				sqlTmp = append(sqlTmp, w)
			}
		}

		return "WHERE " + strings.Join(sqlTmp, " AND ")
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
		return "ORDER BY " + strings.Join(order, ",")
	}
	return ""
}

/**
添加limit
*/
func (db *db) addLimit() string {
	if db.limit > 0 {
		return fmt.Sprintf("LIMIT %d", db.limit)
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
	return "FROM"
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
