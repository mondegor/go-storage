package mrpostgres

import (
    "context"
    "strings"

    "github.com/mondegor/go-webcore/mrctx"
)

func (c *ConnAdapter) debugQuery(ctx context.Context, query string) {
    mrctx.Logger(ctx).Debug("SQL Query: %s", strings.Join(strings.Fields(query), " "))
}
