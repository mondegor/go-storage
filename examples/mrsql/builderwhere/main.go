package main

import (
	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrstorage"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel)
	condBuilder := part.NewSQLConditionBuilder()

	partSql := condBuilder.BuildFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinOr(
				c.JoinAnd(
					c.Equal("equal_field1-1", "1-1"),
					c.NotEqual("not_equal_field1-2", "1-2"),
					c.FilterLike("like_field1-3", "1-3"),
					c.FilterEqualInt64("equalInt_field1-4", 10000, 0),
					c.FilterRangeFloat64("equalInt_field1-5", mrtype.RangeFloat64{Min: 1.34, Max: 2.81}, 0, 0.0001),
				),
				c.JoinAnd(
					c.Equal("equal_field2-1", "2-1"),
					c.NotEqual("not_equal_field2-2", "2-2"),
					c.FilterLike("like_field2-3", "2-3"),
					c.FilterEqualBool("bool_field2-4", mrtype.CastBoolToPointer(true)),
					c.Less("equal_field2-5", "2-5"),
					c.LessOrEqual("equal_field2-6", "2-6"),
				),
				c.JoinAnd(
					c.JoinOr(
						c.Equal("equal_field3-1-1", "3-1-1"),
						c.NotEqual("not_equal_field3-1-2", "3-1-2"),
						c.FilterLikeFields([]string{"like_field3-1-3#1", "like_field3-1-3#2"}, "3-1-3"),
						c.Greater("equal_field3-1-4", "3-1-4"),
						c.GreaterOrEqual("equal_field3-1-5", "3-1-5"),
					),
					c.JoinOr(
						c.Equal("equal_field3-2-1", "3-2-1"),
						c.NotEqual("not_equal_field3-2-2", "3-2-2"),
						c.FilterLike("like_field3-2-3", "3-2-3"),
					),
					c.FilterEqual("like_field3-2-4", uuid.New()),
				),
			)
		},
	)

	cc, vv := partSql.WithStartArg(4).ToSQL()

	logger.Info().Msgf("generated sql: %v", cc)
	logger.Info().Msgf("generated args: %v", vv)
}
