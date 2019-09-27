package fake_orm

import "database/sql"

type where struct {
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

type db struct {
	conn     *sql.DB
	table    string
	join     []join
	fields   []string
	where    []where
	whereRaw []string
	orderBy  []orderBy
	limit    int64
}
