package mrreq

import (
    "net/http"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrreq"
)

const (
    maxListPagerSize = 1000000
)

func ParseListPager(r *http.Request, keyIndex string, keySize string) (mrentity.ListPager, error) {
    pager := mrentity.ListPager{}
    index, err := mrreq.ParseInt64(r, keyIndex, false)

    if index < 0 {
        index = 0
    }

    if err != nil {
        return pager, err
    }

    size, err := mrreq.ParseInt64(r, keySize, false)

    if err != nil {
        return pager, err
    }

    if size < 0 {
        size = 0
    }

    if size > maxListPagerSize {
        return pager, mrcore.FactoryErrHttpRequestParamMax.New(keySize, maxListPagerSize)
    }

    pager.Index = uint64(index)
    pager.Size = uint64(size)

    return pager, nil
}
