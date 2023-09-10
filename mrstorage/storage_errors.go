package mrstorage

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    ErrFactoryConnectionIsAlreadyCreated = NewFactory(
        "errStorageConnectionIsAlreadyCreated", ErrorKindInternal, "connection '{{ .name }}' is already created")

    ErrFactoryConnectionIsNotOpened = NewFactory(
        "errStorageConnectionIsNotOpened", ErrorKindInternal, "connection '{{ .name }}' is not opened")

    ErrFactoryConnectionFailed = NewFactory(
        "errStorageConnectionFailed", ErrorKindSystem, "connection '{{ .name }}' is failed")

    ErrFactoryQueryFailed = NewFactory(
        "errStorageQueryFailed", ErrorKindInternal, "query is failed")

    ErrFactoryFetchDataFailed = NewFactory(
        "errStorageFetchDataFailed", ErrorKindInternal, "fetching data is failed")

    ErrFactoryFetchedInvalidData = NewFactory(
        "errStorageFetchedInvalidData", ErrorKindInternal, "fetched data '{{ .value }}' is invalid")

    ErrFactoryNoRowFound = NewFactory(
        "errStorageNoRowFound", ErrorKindInternalNotice, "no row found")

    ErrFactoryRowsNotAffected = NewFactory(
        "errStorageRowsNotAffected", ErrorKindInternalNotice, "rows not affected")
)
