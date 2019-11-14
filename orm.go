package corm

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

/**
获取一个新的DB
conn 数据库连接
*/
func GetDb(conn *sql.DB) *db {
	db := new(db)
	db.conn = conn
	return db
}

/**
设置数据表
table 表名
*/
func (db *db) Tab(table string) *db {
	db.table = table
	return db
}

/**
设置查询字段，格式：Select("id", "name", "age")
field 查询字段
*/
func (db *db) Select(field ...string) *db {
	db.fields = append(db.fields, field...)
	return db
}

/**
设置查询字段原生格式，格式：Select("id, name, age, IFNULL(sex=1,1,2) AS sex")
field 查询字段
*/
func (db *db) SelectRaw(field string) *db {
	fieldTmp := strings.Split(field, ",")
	db.fields = append(db.fields, fieldTmp...)
	return db
}

/**
查询条件，格式：Where("id", ">", 100).where("name", "=", "张三")
field 查询字段
operator 条件符号 >、<、=、<>、like、in 等
condition 条件值
*/
func (db *db) Where(field, operator string, condition interface{}) *db {
	db.where = append(db.where, where{
		field:     field,
		operator:  operator,
		condition: condition,
	})
	return db
}

/**
将数字字符串转换成INT
*/
func (db *db) WhereStrToInt(field, operator string, condition string) *db {
	//strInt, err := strconv.ParseInt(condition, 10, 64)
	if condition == "" {
		db.pushErr(errors.New("func:WhereStrToInt condition is empty"))
	}
	strInt, err := strconv.Atoi(condition)
	if err != nil {
		db.pushErr(err)
	}
	return db.Where(field, operator, strInt)
}

/**
将int64转换成字符串
*/
func (db *db) WhereInt64ToStr(field, operator string, condition int64) *db {
	intStr := strconv.FormatInt(condition, 10)
	return db.Where(field, operator, intStr)
}

/**
将int转换成字符串
*/
func (db *db) WhereIntToStr(field, operator string, condition int) *db {
	intStr := strconv.Itoa(condition)
	return db.Where(field, operator, intStr)
}

/**
查询条件原生格式，格式：Where("id > 100 and name = '张三'")
where 条件字符串
*/
func (db *db) WhereRaw(where string) *db {
	db.whereRaw = append(db.whereRaw, where)
	return db
}

/**
查询 In 条件，格式：WhereIn("name", "张")
where 条件字符串
*/
func (db *db) WhereIn(field string, condition ...interface{}) *db {
	db.where = append(db.where, where{
		field:          field,
		operator:       IN,
		conditionArray: condition,
	})
	return db
}

/**
查询 Not In 条件，格式：WhereNotIn("name", "张")
where 条件字符串
*/
func (db *db) WhereNotIn(field string, condition ...interface{}) *db {
	db.where = append(db.where, where{
		field:          field,
		operator:       NOT_IN,
		conditionArray: condition,
	})
	return db
}

/**
查询 like 条件，格式：WhereLike("name", "张")
where 条件字符串
*/
func (db *db) WhereLike(field string, condition string) *db {
	condition = "%" + condition + "%"
	db.where = append(db.where, where{
		field:     field,
		operator:  LIKE,
		condition: condition,
	})
	return db
}

/**
查询 not like 条件，格式：WhereNotLike("name", "张")
where 条件字符串
*/
func (db *db) WhereNotLike(field string, condition string) *db {
	condition = "%" + condition + "%"
	db.where = append(db.where, where{
		field:     field,
		operator:  NOT_LIKE,
		condition: condition,
	})
	return db
}

/**
查询 Between 条件，格式：WhereBetween("id", 100, 1000)
where 条件字符串
*/
func (db *db) WhereBetween(field string, startCondition interface{}, endCondition interface{}) *db {
	db.where = append(db.where, where{
		field:          field,
		operator:       BETWEEN,
		conditionArray: []interface{}{startCondition, endCondition},
	})
	return db
}

/**
查询结果过滤 Having ，格式：Having("name", "=", "张三").Having("age", ">", 18)
Having 条件字符串
*/
func (db *db) Having(field, operator string, condition interface{}) *db {
	db.having = append(db.having, having{
		field:     field,
		operator:  operator,
		condition: condition,
	})
	return db
}

/**
排序，格式：OrderBy("id", "desc").OrderBy("name", "asc")
field 字段
by asc或desc
*/
func (db *db) OrderBy(field, by string) *db {
	db.orderBy = append(db.orderBy, orderBy{
		field: field,
		by:    by,
	})
	return db
}

/**
分组，格式：GroupBy("class", "type")
field 字段
by asc或desc
*/
func (db *db) GroupBy(field ...string) *db {
	db.groupBy = append(db.groupBy, field...)
	return db
}

/**
查询结果数量，格式：Limit(100)
limit 数量
*/
func (db *db) Limit(limit int) *db {
	db.limit = limit
	return db
}

/**
查询结果数量，格式：Limit(100)
limit 数量
*/
func (db *db) Offset(offset int) *db {
	db.offset = offset
	return db
}

/**
左连接，格式：LeftJoin("group", "user.group_id = group.id")
table 表名
on 关联条件
*/
func (db *db) LeftJoin(table, on string) *db {
	db.join = append(db.join, join{
		table:     table,
		direction: LEFT_JOIN,
		on:        on,
	})
	return db
}

/**
左连接，格式：RightJoin("group", "user.group_id = group.id")
table 表名
on 关联条件
*/
func (db *db) RightJoin(table, on string) *db {
	db.join = append(db.join, join{
		table:     table,
		direction: RIGHT_JOIN,
		on:        on,
	})
	return db
}

/**
左连接，格式：Join("group", "user.group_id = group.id")
table 表名
on 关联条件
*/
func (db *db) Join(table, on string) *db {
	db.join = append(db.join, join{
		table:     table,
		direction: INNER_JOIN,
		on:        on,
	})
	return db
}

/**
查询一条数据
callable 回调函数
*/
func (db *db) First(result ...interface{}) error {
	err := db.QueryRow(db.whereToSql(), db.getWhereValue(), result...)
	if errs(err) != nil {
		return err
	}
	return nil
}

/**
查询多条数据
callable 回调函数
*/
func (db *db) Get(callable func(rows *sql.Rows)) error {
	rows, err := db.Query(db.whereToSql(), db.getWhereValue()...)
	if errs(err) != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		callable(rows)
	}
	return nil
}

/**
分页数据查询
page 页数
pageCount 每页记录数
callable 回调函数
*/
func (db *db) GetPage(page, pageCount int, callable func(rows *sql.Rows)) (int64, error) {
	db.offset = (page - 1) * pageCount
	db.limit = pageCount
	//总记录数
	totalCount, err := db.Count()
	if err != nil {
		return 0, err
	}
	err = db.Get(callable)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

/**
Sum
*/
func (db *db) Sum(sumField string) (float64, error) {
	db.sum = sumField
	var sum sql.NullFloat64
	err := db.QueryRow(db.sumToSql(), db.getWhereValue(), &sum)
	if errs(err) != nil {
		return 0, err
	}
	return sum.Float64, nil
}

/**
Sum
*/
func (db *db) Max(maxField string) (int64, error) {
	db.max = maxField
	var max sql.NullInt64
	err := db.QueryRow(db.maxToSql(), db.getWhereValue(), &max)
	if errs(err) != nil {
		return 0, err
	}
	return max.Int64, nil
}

/**
Sum
*/
func (db *db) Min(minField string) (int64, error) {
	db.min = minField
	var min sql.NullInt64
	err := db.QueryRow(db.minToSql(), db.getWhereValue(), &min)
	if errs(err) != nil {
		return 0, err
	}
	return min.Int64, nil
}

/**
Count
*/
func (db *db) Count() (int64, error) {
	var count sql.NullInt64
	err := db.QueryRow(db.countToSql(), db.getWhereValue(), &count)
	if errs(err) != nil {
		return 0, err
	}
	return count.Int64, nil
}

/**
Exists 查询数据是否存在
*/
func (db *db) Exists() (bool, error) {
	count, err := db.Count()
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

/**
打印SQL
*/
func (db *db) PrintSql() string {
	return db.whereToSql()
}

/**
插入数据
*/
func (db *db) Insert(insertMap map[string]interface{}) (LastInsertId int64, err error) {
	db.insert = insertMap
	insertStr, vals := db.insertToSql()

	rest, err := db.Exec(insertStr, vals...)
	if err != nil {
		return 0, err
	}
	insertId, err := rest.LastInsertId()
	if err != nil {
		return 0, err
	}
	return insertId, nil
}

/**
修改数据
*/
func (db *db) Update(updateMap map[string]interface{}) (updateNum int64, err error) {
	db.update = updateMap
	updateStr, vals := db.updateToSql()

	vals = append(vals, db.getWhereValue()...)
	rest, err := db.Exec(updateStr, vals...)
	if err != nil {
		return 0, err
	}
	rows, err := rest.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}
