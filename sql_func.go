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
}

//同一个实例多次调用，清除条件
func (db *Db) clear() {
	db.table = ""
	db.sum = ""
	db.count = ""
	db.max = ""
	db.min = ""

	db.join = nil
	db.fields = nil
	db.where = nil
	db.whereRaw = nil
	db.orderBy = nil
	db.groupBy = nil
	db.having = nil
	db.insert = nil
	db.update = nil

	db.limit = 0
	db.offset = 0

	db.buffer = bytes.Buffer{}
	db.err = nil
}
