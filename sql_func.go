package corm

import (
	"bytes"
	"database/sql"
	"github.com/pkg/errors"
	"strings"
)

/**
query 查询语句
args 查询参数
scan 结果绑定参数
*/
func (db *Db) queryRow(query string, args []interface{}, scan ...interface{}) error {
	defer db.putPool()
	if db.getErr() != nil {
		return db.getErr()
	}

	if db.tx != nil {
		defer db.clear()
		return db.tx.QueryRow(query, args...).Scan(scan...)
	}
	return db.conn.QueryRow(query, args...).Scan(scan...)
}

/**
query 查询语句
args 查询参数
*/
func (db *Db) query(query string, args ...interface{}) (*sql.Rows, error) {
	defer db.putPool()
	if db.getErr() != nil {
		return nil, db.getErr()
	}

	if db.tx != nil {
		defer db.clear()
		return db.tx.Query(query, args...)
	}
	return db.conn.Query(query, args...)
}

/**
query 执行语句
args 查询参数
*/
func (db *Db) exec(sqlStr string, args ...interface{}) (sql.Result, error) {
	defer db.putPool()
	if db.getErr() != nil {
		return nil, db.getErr()
	}

	var stmt *sql.Stmt
	var err error

	if db.tx != nil {
		defer db.clear()
		stmt, err = db.tx.Prepare(sqlStr)
	} else {
		stmt, err = db.conn.Prepare(sqlStr)
	}

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rest, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return rest, nil
}

func errs(err error) error {
	if err != nil {
		//数据库记录不存在，报此错，可忽略
		if err == sql.ErrNoRows {
			return nil
		}
		//字段值为NULL时，报此错，可忽略，默认为类型的零值
		if strings.Index(err.Error(), "sql: Scan error on column index") != -1 {
			return nil
		}
	}
	return err
}

func (db *Db) pushErr(err error) {
	if err != nil {
		db.err = append(db.err, err)
		//db.err = fmt.Errorf("%w", err)
	}
}

func (db *Db) getErr() error {
	if len(db.err) > 0 {
		return db.err[0]
	}
	return nil
}

func (db *Db) check() {
	if db.table == "" {
		db.pushErr(errors.New("未定义数据表"))
	}
	db.buffer = bytes.Buffer{}
}

//同一个实例多次调用，清除条件
func (db *Db) clear() {
	//*db = Db{conn: db.conn, tx: db.tx}
	db.table, db.sum, db.count, db.max, db.min = "", "", "", "", ""
	db.join, db.fields, db.where, db.whereRaw, db.orderBy, db.groupBy, db.having, db.insert, db.update, db.err, db.tx = nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	db.limit, db.offset = 0, 0
	db.buffer = bytes.Buffer{}
}

//检测whereIn条件的参数类型是否正确
func checkWhereIn(condition []interface{}) bool {
	if len(condition) > 0 {
		switch condition[0].(type) {
		case string, int, int32, int64:
		default:
			return true
		}
	}
	return false
}

//将对象放回池中
func (db *Db) putPool() {
	db.clear()
	dbPool.Put(db)
}

/**
克隆DB对象
注意：此方法只是结构的浅拷贝，只能拷贝slice、map的指针，并不能拷贝值，只能用在相同条件下的不同操作
	 克隆后不能修改对象的slice及map类型成员
*/
func (db *Db) clone() *Db {
	newDB := dbPool.Get().(*Db)
	*newDB = *db
	return newDB
}

/**
判断是否为0值
*/
func IsDefaultValue(obj interface{}) bool {
	switch value := obj.(type) {
	case int, int16, int32, int64:
		if value == 0 {
			return true
		}
	case string:
		if value == "" {
			return true
		}
	case float32, float64:
		if value == 0. {
			return true
		}
	}
	return false
}
