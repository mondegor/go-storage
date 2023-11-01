package mrreq

import (
    "net/http"
    "regexp"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
)

const (
    maxListSorterField = 32
)

var (
    regexpSorterField = regexp.MustCompile(`^[a-z]([a-zA-Z0-9]+)?[a-zA-Z0-9]$`)
)

func ParseListSorter(r *http.Request, keyField string, keyDirection string) (mrentity.ListSorter, error) {
    sorter := mrentity.ListSorter{}
    query := r.URL.Query()

    value := query.Get(keyField)

    if value == "" {
        return sorter, nil
    }

    if len(value) > maxListSorterField {
        return sorter, mrcore.FactoryErrHttpRequestParamLenMax.New(keyField, maxListSorterField)
    }

    if !regexpSorterField.MatchString(value) {
        return sorter, mrcore.FactoryErrHttpRequestParseParam.New("ListSorter", keyField, value)
    }

    direction := query.Get(keyDirection)

    if direction != "" {
        err := sorter.Direction.ParseAndSet(direction)
        return sorter, err
    }

    sorter.FieldName = value

    return sorter, nil
}
