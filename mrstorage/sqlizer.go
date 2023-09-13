package mrstorage

type Sqlizer interface {
    ToSql() (string, []interface{}, error)
}
