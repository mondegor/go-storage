package main

import (
	"fmt"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrtype"
)

func main() {
	w := mrpostgres.NewSqlBuilderWhere()

	where := w.JoinOr(
		w.JoinAnd(
			w.Equal("equal_field1-1", "1-1"),
			w.NotEqual("not_equal_field1-2", "1-2"),
			w.FilterLike("like_field1-3", "1-3"),
			w.FilterEqualInt64("equalInt_field1-4", 10000, 0),
		),
		w.JoinAnd(
			w.Equal("equal_field2-1", "2-1"),
			w.NotEqual("not_equal_field2-2", "2-2"),
			w.FilterLike("like_field2-3", "2-3"),
			w.FilterEqualBool("bool_field2-4", mrtype.NullableBoolTrue),
		),
		w.JoinAnd(
			w.JoinOr(
				w.Equal("equal_field3-1-1", "3-1-1"),
				w.NotEqual("not_equal_field3-1-2", "3-1-2"),
				w.FilterLikeFields([]string{"like_field3-1-3#1", "like_field3-1-3#2"}, "3-1-3"),
			),
			w.JoinOr(
				w.Equal("equal_field3-2-1", "3-2-1"),
				w.NotEqual("not_equal_field3-2-2", "3-2-2"),
				w.FilterLike("like_field3-2-3", "3-2-3"),
			),
		),
	)

	cc, vv := where(5)

	fmt.Printf("%v\n", cc)
	fmt.Printf("%v\n", vv)
}
