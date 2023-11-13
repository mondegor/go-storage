package mrsql

import "regexp"

const (
    fieldTagDBFieldName = "db"
    fieldTagFieldUpdate  = "upd"
    fieldTagSortByField = "sort"
)

var (
    regexpDbName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

