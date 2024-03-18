package main

import (
	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel)
	bw := mrpostgres.NewSqlBuilderWhere()

	where := bw.JoinOr(
		bw.JoinAnd(
			bw.Equal("equal_field1-1", "1-1"),
			bw.NotEqual("not_equal_field1-2", "1-2"),
			bw.FilterLike("like_field1-3", "1-3"),
			bw.FilterEqualInt64("equalInt_field1-4", 10000, 0),
		),
		bw.JoinAnd(
			bw.Equal("equal_field2-1", "2-1"),
			bw.NotEqual("not_equal_field2-2", "2-2"),
			bw.FilterLike("like_field2-3", "2-3"),
			bw.FilterEqualBool("bool_field2-4", mrtype.BoolPointer(true)),
			bw.Less("equal_field2-5", "2-5"),
			bw.LessOrEqual("equal_field2-6", "2-6"),
		),
		bw.JoinAnd(
			bw.JoinOr(
				bw.Equal("equal_field3-1-1", "3-1-1"),
				bw.NotEqual("not_equal_field3-1-2", "3-1-2"),
				bw.FilterLikeFields([]string{"like_field3-1-3#1", "like_field3-1-3#2"}, "3-1-3"),
				bw.Greater("equal_field3-1-4", "3-1-4"),
				bw.GreaterOrEqual("equal_field3-1-5", "3-1-5"),
			),
			bw.JoinOr(
				bw.Equal("equal_field3-2-1", "3-2-1"),
				bw.NotEqual("not_equal_field3-2-2", "3-2-2"),
				bw.FilterLike("like_field3-2-3", "3-2-3"),
			),
			bw.FilterEqualUUID("like_field3-2-4", uuid.New()),
		),
	)

	cc, vv := where(5)

	logger.Info().Msgf("%v", cc)
	logger.Info().Msgf("%v", vv)
}
