package main

import (
	"fmt"
	"time"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrenum"
)

func main() {
	fmt.Println("SAMPLE1:")

	o := mrpostgres.NewSqlBuilderOrderBy("id", mrenum.SortDirectionASC)

	orderBy := o.Join(
		o.Field("caption", mrenum.SortDirectionDESC),
		o.Field("createdAt", mrenum.SortDirectionDESC),
	)

	cc, _ := orderBy(0)
	fmt.Printf("%v\n\n", cc)

	fmt.Println("SAMPLE2:")
	o = mrpostgres.NewSqlBuilderOrderBy("id", mrenum.SortDirectionDESC)

	orderBy = o.WrapWithDefault(
		o.Field("", mrenum.SortDirectionASC),
	)

	cc, _ = orderBy(0)
	fmt.Printf("%v\n\n", cc)

	fmt.Println("SAMPLE3:")

	type OrderedStruct struct {
		ID        string    `sort:"id"`
		Caption   string    `sort:"caption"`
		CreatedAt time.Time `sort:"createdAt,default,desc"`
		NotSorted string
		IsRemoved bool `sort:"isRemoved"`
	}

	meta, _ := mrsql.NewEntityMetaOrderBy(OrderedStruct{})
	fmt.Printf("caption is registered? %v\n", meta.CheckField("caption"))
	fmt.Printf("NotSorted is registered? %v\n", meta.CheckField("NotSorted"))

	o = mrpostgres.NewSqlBuilderOrderByWithDefaultSort(meta.DefaultSort())

	orderBy = o.WrapWithDefault(
		o.Field("", mrenum.SortDirectionASC),
	)

	cc, _ = orderBy(0)
	fmt.Printf("%v\n\n", cc)
}
