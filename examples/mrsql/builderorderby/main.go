package main

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel)
	ctx := mrlog.WithContext(context.Background(), logger)

	logger.Info().Msg("SAMPLE1:")
	o := mrpostgres.NewSQLBuilderOrderBy(
		ctx,
		mrtype.SortParams{
			FieldName: "id",
			Direction: mrenum.SortDirectionDESC,
		},
	)

	orderBy := o.Join(
		o.Field("caption", mrenum.SortDirectionASC),
		o.Field("createdAt", mrenum.SortDirectionDESC),
	)

	cc, _ := orderBy(0)
	logger.Info().Msgf("%v", cc)

	logger.Info().Msg("SAMPLE2:")
	o = mrpostgres.NewSQLBuilderOrderBy(
		ctx,
		mrtype.SortParams{
			FieldName: "id",
			Direction: mrenum.SortDirectionDESC,
		},
	)

	orderBy = o.DefaultField()

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

	o = mrpostgres.NewSQLBuilderOrderBy(ctx, meta.DefaultSort())

	orderBy = o.DefaultField()

	cc, _ = orderBy(0)
	logger.Info().Msgf("%v", cc)
}
