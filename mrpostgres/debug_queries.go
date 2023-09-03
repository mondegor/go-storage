package mrpostgres

import (
    "context"
    "strings"

    "github.com/mondegor/go-core/mrcore"
)

func (c *Connection) debugQuery(ctx context.Context, query string) {
    mrcore.ExtractLogger(ctx).Debug("SQL Query: %s", strings.Join(strings.Fields(query), " "))
}
