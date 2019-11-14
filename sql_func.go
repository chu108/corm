package corm

import (
	"clone_packge/tools/go/ssa/interp/testdata/src/strings"
	"database/sql"
	"fmt"
)

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
