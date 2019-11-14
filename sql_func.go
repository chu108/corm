package corm

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

/**
query 查询语句
args 查询参数
scan 结果绑定参数
*/
func (db *db) QueryRow(query string, args []interface{}, scan ...interface{}) error {
	if db.err != nil {
		return db.err
	}
	return db.conn.QueryRow(query, args...).Scan(scan...)
}

/**
query 查询语句
args 查询参数
*/
func (db *db) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db.err != nil {
		return nil, db.err
	}
	return db.conn.Query(query, args...)
}

/**
query 查询语句
args 查询参数
*/
func (db *db) Exec(sql string, args ...interface{}) (sql.Result, error) {
	if db.err != nil {
		return nil, db.err
	}
	stmt, err := db.conn.Prepare(sql)
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

func (db *db) pushErr(err error) {
	db.err = fmt.Errorf("%w", err)
}

func (db *db) getErr() error {
	return db.err
}

func (db *db) check() {
	if db.table == "" {
		db.pushErr(errors.New("未定义数据表"))
	}
}
