package main

import (
	"os"
	"time"

	"github.com/mondegor/go-sysmess/mrlog/litelog"
	"github.com/mondegor/go-sysmess/mrlog/slog"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/mrtype/enums"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

func main() {
	l, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))
	logger := litelog.NewLogger(l)

	logger.Info("SAMPLE1:")
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

	logger.Info("generated sql", "value", orderBy.String())

	logger.Info("SAMPLE2:")
	orderByBuilder = part.NewSQLOrderByBuilder(
		mrtype.SortParams{
			FieldName: "id",
			Direction: enums.SortDirectionDESC,
		},
	)

	orderBy = orderByBuilder.Build(nil) // return default value

	logger.Info("generated sql", "value", orderBy.String())

	logger.Info("SAMPLE3:")
	type OrderedStruct struct {
		ID        string    `sort:"id"`
		Caption   string    `sort:"caption"`
		CreatedAt time.Time `sort:"createdAt,default,desc"`
		NotSorted string
		IsRemoved bool `sort:"isRemoved"`
	}

	meta, _ := mrsql.NewEntityMetaOrderBy(l, OrderedStruct{})
	logger.Info("caption is registered?", "value", meta.CheckField("caption"))
	logger.Info("NotSorted is registered?", "value", meta.CheckField("NotSorted"))

	orderByBuilder = part.NewSQLOrderByBuilder(meta.DefaultSort())

	orderBy = orderByBuilder.Build(nil) // return default value

	logger.Info("generated sql", "value", orderBy.String())
}
