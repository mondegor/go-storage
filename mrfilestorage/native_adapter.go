package mrfilestorage

import "strings"

type (
    nativeAdapter struct {
        rootDir string
    }
)

func New(rootDir string) *nativeAdapter {
    return &nativeAdapter{
        rootDir: strings.TrimRight(rootDir, "/"),
    }
}
