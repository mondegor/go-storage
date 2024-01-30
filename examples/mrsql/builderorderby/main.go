package main

import (
	"context"
	"time"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel)
	ctx := mrlog.WithContext(context.Background(), logger)

	logger.Info().Msg("SAMPLE1:")

	o := mrpostgres.NewSqlBuilderOrderBy(ctx, "id", mrenum.SortDirectionASC)

	orderBy := o.Join(
		o.Field("caption", mrenum.SortDirectionDESC),
		o.Field("createdAt", mrenum.SortDirectionDESC),
	)

	cc, _ := orderBy(0)
	logger.Info().Msgf("%v", cc)

	logger.Info().Msg("SAMPLE2:")
	o = mrpostgres.NewSqlBuilderOrderBy(ctx, "id", mrenum.SortDirectionDESC)

	orderBy = o.WrapWithDefault(
		o.Field("", mrenum.SortDirectionASC),
	)

	cc, _ = orderBy(0)
	logger.Info().Msgf("%v", cc)

	logger.Info().Msg("SAMPLE3:")

	type OrderedStruct struct {
		ID        string    `sort:"id"`
		Caption   string    `sort:"caption"`
		CreatedAt time.Time `sort:"createdAt,default,desc"`
		NotSorted string
		IsRemoved bool `sort:"isRemoved"`
	}

	meta, _ := mrsql.NewEntityMetaOrderBy(ctx, OrderedStruct{})
	logger.Info().Msgf("caption is registered? %v", meta.CheckField("caption"))
	logger.Info().Msgf("NotSorted is registered? %v", meta.CheckField("NotSorted"))

	o = mrpostgres.NewSqlBuilderOrderByWithDefaultSort(ctx, meta.DefaultSort())

	orderBy = o.WrapWithDefault(
		o.Field("", mrenum.SortDirectionASC),
	)

	cc, _ = orderBy(0)
	logger.Info().Msgf("%v", cc)
}
