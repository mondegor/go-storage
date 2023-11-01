package mrentity

import "io"

const (
    ModelNameFile = "File"
)

type (
    File struct {
        ContentType string
        Name string
        Size int64
        Body io.ReadCloser
    }
)
