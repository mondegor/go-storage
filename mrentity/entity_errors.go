package mrentity

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryInternalListOfFieldsIsEmpty = NewFactory(
        "errInternalListOfFieldsIsEmpty", ErrorKindInternalNotice, "the list of fields is empty")
)
