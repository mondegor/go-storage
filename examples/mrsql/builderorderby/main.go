package main

import (
	"os"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/slog"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/mrtype/enums"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

func main() {
	logger, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))

	mrlog.Info(logger, "SAMPLE1:")
	orderByBuilder := part.NewSQLOrderByBuilder(
		mrtype.SortParams{
			FieldName: "id",
			Direction: enums.SortDirectionDESC,
		},
	)

	orderBy := orderByBuilder.BuildFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field("caption", enums.SortDirectionASC),
				o.Field("createdAt", enums.SortDirectionDESC),
			)
		},
	)

	mrlog.Info(logger, "generated sql", "value", orderBy.String())

	mrlog.Info(logger, "SAMPLE2:")
	orderByBuilder = part.NewSQLOrderByBuilder(
		mrtype.SortParams{
			FieldName: "id",
			Direction: enums.SortDirectionDESC,
		},
	)

	orderBy = orderByBuilder.Build(nil) // return default value

	mrlog.Info(logger, "generated sql", "value", orderBy.String())

	mrlog.Info(logger, "SAMPLE3:")
	type OrderedStruct struct {
		ID        string    `sort:"id"`
		Caption   string    `sort:"caption"`
		CreatedAt time.Time `sort:"createdAt,default,desc"`
		NotSorted string
		IsRemoved bool `sort:"isRemoved"`
	}

	meta, _ := mrsql.NewEntityMetaOrderBy(logger, OrderedStruct{})
	mrlog.Info(logger, "caption is registered?", "value", meta.HasField("caption"))
	mrlog.Info(logger, "NotSorted is registered?", "value", meta.HasField("NotSorted"))

	orderByBuilder = part.NewSQLOrderByBuilder(meta.DefaultSort())

	orderBy = orderByBuilder.Build(nil) // return default value

	mrlog.Info(logger, "generated sql", "value", orderBy.String())
}
