package mrentity

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrInternalListOfFieldsIsEmpty = NewFactory(
        "errInternalListOfFieldsIsEmpty", ErrorKindInternalNotice, "the list of fields is empty")
)
