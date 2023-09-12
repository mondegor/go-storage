package mrentity

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    factoryErrInternalListOfFieldsIsEmpty = NewFactory(
        "errInternalListOfFieldsIsEmpty", ErrorKindInternalNotice, "the list of fields is empty")
)
