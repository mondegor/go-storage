package main

import (
	"time"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel)

	logger.Info().Msg("SAMPLE1:")
	orderByBuilder := part.NewSQLOrderByBuilder(
		mrtype.SortParams{
			FieldName: "id",
			Direction: mrenum.SortDirectionDESC,
		},
	)

	orderBy := orderByBuilder.BuildFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field("caption", mrenum.SortDirectionASC),
				o.Field("createdAt", mrenum.SortDirectionDESC),
			)
		},
	)

	logger.Info().Msgf("generated sql: %v", orderBy.String())

	logger.Info().Msg("SAMPLE2:")
	orderByBuilder = part.NewSQLOrderByBuilder(
		mrtype.SortParams{
			FieldName: "id",
			Direction: mrenum.SortDirectionDESC,
		},
	)

	orderBy = orderByBuilder.Build(nil) // return default value

	logger.Info().Msgf("generated sql: %v", orderBy.String())

	logger.Info().Msg("SAMPLE3:")
	type OrderedStruct struct {
		ID        string    `sort:"id"`
		Caption   string    `sort:"caption"`
		CreatedAt time.Time `sort:"createdAt,default,desc"`
		NotSorted string
		IsRemoved bool `sort:"isRemoved"`
	}

	meta, _ := mrsql.NewEntityMetaOrderBy(logger, OrderedStruct{})
	logger.Info().Msgf("caption is registered? %t", meta.CheckField("caption"))
	logger.Info().Msgf("NotSorted is registered? %t", meta.CheckField("NotSorted"))

	orderByBuilder = part.NewSQLOrderByBuilder(meta.DefaultSort())

	orderBy = orderByBuilder.Build(nil) // return default value

	logger.Info().Msgf("generated sql: %v", orderBy.String())
}
