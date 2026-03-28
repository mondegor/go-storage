package main

import (
	"os"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-storage/mrpostgres/builder/helper"
	"github.com/mondegor/go-storage/mrpostgres/builder/part"
)

func main() {
	logger, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))

	condBuilder := part.NewSQLConditionBuilder()
	condHelper := helper.NewSQLCondition()

	partFunc1 := condHelper.JoinAnd(
		condHelper.Equal("part1_item1", "equal"),
		condHelper.Expr("part1_item2 IS NULL"),
	)

	partFunc2 := condHelper.JoinAnd(
		condHelper.Expr("part2_item1 = 'value2_1'"),
		condHelper.FilterEqualInt64("part2_item2", 2222, 0),
		condHelper.FilterEqualString("part2_item3", "value2_3"),
	)

	joinedParts := condBuilder.BuildAnd(partFunc1, partFunc2).WithStartArg(5)
	cc, vv := joinedParts.ToSQL()

	mrlog.Info(logger, "generated sql", "value", cc)
	mrlog.Info(logger, "generated args", "value", vv)
}
