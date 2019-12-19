package corm

import (
	"bytes"
	"database/sql"
)

type where struct {
	field          string
	operator       string
	condition      interface{}
	conditionArray []interface{}
}

type having struct {
	field     string
	operator  string
	condition interface{}
}

type join struct {
	table     string
	direction string
	on        string
}

type orderBy struct {
	field string
	by    string
}

type Db struct {
	conn     *sql.DB
	tx       *sql.Tx
	err      []error
	table    string
	force    string
	join     []join
	fields   []string
	where    []where
	whereRaw []string
	orderBy  []orderBy
	groupBy  []string
	limit    int
	offset   int
	having   []having
	sum      string
	count    string
	max      string
	min      string
	insert   map[string]interface{}
	update   map[string]interface{}
	compose  []string
	buffer   bytes.Buffer
}
