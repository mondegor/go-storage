package main

import (
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel)
	cond := mrpostgres.NewSQLBuilderCondition(
		mrpostgres.NewSQLBuilderWhere(),
	)

	part1 := cond.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
		return w.JoinAnd(
			w.Equal("part1_item1", "equal"),
			w.Expr("part1_item2 IS NULL"),
		)
	})

	part2 := cond.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
		return w.JoinAnd(
			w.Expr("part2_item1 = 'value2_1'"),
			w.FilterEqualInt64("part2_item2", 2222, 0),
			w.FilterEqualString("part2_item3", "value2_3"),
		)
	})

	joinedParts := part1.WithPart(" AND ", part2).WithParam(5)
	cc, vv := joinedParts.ToSQL()

	logger.Info().Msgf("%v", cc)
	logger.Info().Msgf("%v", vv)
}
