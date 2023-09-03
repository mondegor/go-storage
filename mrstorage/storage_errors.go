package mrstorage

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryConnectionIsAlreadyCreated = NewFactory(
        "errStorageConnectionIsAlreadyCreated", ErrorKindInternal, "connection '{{ .name }}' is already created")

    FactoryConnectionIsNotOpened = NewFactory(
        "errStorageConnectionIsNotOpened", ErrorKindInternal, "connection '{{ .name }}' is not opened")

    FactoryConnectionFailed = NewFactory(
        "errStorageConnectionFailed", ErrorKindSystem, "connection '{{ .name }}' is failed")

    FactoryQueryFailed = NewFactory(
        "errStorageQueryFailed", ErrorKindInternal, "query is failed")

    FactoryFetchDataFailed = NewFactory(
        "errStorageFetchDataFailed", ErrorKindInternal, "fetching data is failed")

    FactoryFetchedInvalidData = NewFactory(
        "errStorageFetchedInvalidData", ErrorKindInternal, "fetched data '{{ .value }}' is invalid")

    FactoryNoRowFound = NewFactory(
        "errStorageNoRowFound", ErrorKindInternalNotice, "no row found")

    FactoryRowsNotAffected = NewFactory(
        "errStorageRowsNotAffected", ErrorKindInternalNotice, "rows not affected")
)
