package mrreq

import (
    "net/http"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrreq"
)

func ParseRangeInt64(r *http.Request, key string) (mrentity.RangeInt64, error) {
    min, err := mrreq.ParseInt64(r, key + "-min", false)

    if err != nil {
        return mrentity.RangeInt64{}, err
    }

    max, err := mrreq.ParseInt64(r, key + "-max", false)

    if err != nil {
        return mrentity.RangeInt64{}, err
    }

    if min > max { // change
        return mrentity.RangeInt64{
            Min: max,
            Max: min,
        }, nil
    }

    return mrentity.RangeInt64{
        Min: min,
        Max: max,
    }, nil
}
