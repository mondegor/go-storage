package main

import (
	"os"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrstorage"
)

func main() {
	logger, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))

	condBuilder := part.NewSQLConditionBuilder()

	partFunc1 := condBuilder.HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Equal("part1_item1", "equal"),
				c.Expr("part1_item2 IS NULL"),
			)
		},
	)

	partFunc2 := condBuilder.HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("part2_item1 = 'value2_1'"),
				c.FilterEqualInt64("part2_item2", 2222, 0),
				c.FilterEqualString("part2_item3", "value2_3"),
			)
		},
	)

	joinedParts := condBuilder.BuildAnd(partFunc1, partFunc2).WithStartArg(5)
	cc, vv := joinedParts.ToSQL()

	mrlog.Info(logger, "generated sql", "value", cc)
	mrlog.Info(logger, "generated args", "value", vv)
}
