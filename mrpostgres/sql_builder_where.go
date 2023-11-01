package mrpostgres

import (
    "fmt"
    "reflect"
    "strings"

    "github.com/lib/pq"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrsql"
    "github.com/mondegor/go-storage/mrstorage"
)

// go get -u github.com/lib/pq

type (
    SqlBuilderWhere struct {
    }
)

func NewSqlBuilderWhere() *SqlBuilderWhere {
    return &SqlBuilderWhere{
    }
}

func (b *SqlBuilderWhere) JoinAnd(conds ...mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
    return b.join(" AND ", conds)
}

func (b *SqlBuilderWhere) JoinOr(conds ...mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
    return b.join(" OR ", conds)
}

func (b *SqlBuilderWhere) Expr(expr string) mrstorage.SqlBuilderPartFunc {
    if expr == "" {
        return nil
    }

    return func(paramNumber int) (string, []any) {
        return expr, []any{}
    }
}

func (b *SqlBuilderWhere) ExprWithValue(expr string, value any) mrstorage.SqlBuilderPartFunc {
    if expr == "" {
        return nil
    }

    return func(paramNumber int) (string, []any) {
        return fmt.Sprintf(expr, paramNumber), []any{value}
    }
}

func (b *SqlBuilderWhere) Equal(dbName string, value any) mrstorage.SqlBuilderPartFunc {
    return func (paramNumber int) (string, []any) {
        return fmt.Sprintf("%s = $%d", dbName, paramNumber), []any{value}
    }
}

func (b *SqlBuilderWhere) NotEqual(dbName string, value any) mrstorage.SqlBuilderPartFunc {
    return func (paramNumber int) (string, []any) {
        return fmt.Sprintf("%s <> $%d", dbName, paramNumber), []any{value}
    }
}

func (b *SqlBuilderWhere) FilterLike(dbName string, value string) mrstorage.SqlBuilderPartFunc {
    return b.FilterLikeFields([]string{dbName}, value)
}

func (b *SqlBuilderWhere) FilterLikeFields(dbNames []string, value string) mrstorage.SqlBuilderPartFunc {
   if value == "" {
       return nil
   }

    return func (paramNumber int) (string, []any) {
       var conds []string

       for i := range dbNames {
            conds = append(conds, fmt.Sprintf("%s LIKE '%%' || $%d || '%%'", dbNames[i], paramNumber))
       }

       return fmt.Sprintf("(%s)", strings.Join(conds, " OR ")), []any{value}
   }
}

func (b *SqlBuilderWhere) FilterEqualInt64(dbName string, value int64, empty int64) mrstorage.SqlBuilderPartFunc {
    if value == empty {
        return nil
    }

    return func (paramNumber int) (string, []any) {
        return fmt.Sprintf("%s = $%d", dbName, paramNumber), []any{value}
    }
}

func (b *SqlBuilderWhere) FilterRangeInt64(dbName string, value mrentity.RangeInt64, empty int64) mrstorage.SqlBuilderPartFunc {
   if value.Min != empty {
       if value.Max != empty {
           if value.Min > value.Max {
               return nil
           }

           return func (paramNumber int) (string, []any) {
               return fmt.Sprintf("(%s BETWEEN $%d AND $%d)", dbName, paramNumber, paramNumber + 1), []any{value.Min, value.Max}
           }
       } else {
           return func (paramNumber int) (string, []any) {
               return fmt.Sprintf("%s >= $%d", dbName, paramNumber), []any{value.Min}
           }
       }
   } else if value.Max != empty {
       return func (paramNumber int) (string, []any) {
           return fmt.Sprintf("%s <= $%d", dbName, paramNumber), []any{value.Max}
       }
   }

   return nil
}

func (b *SqlBuilderWhere) FilterAnyOf(dbName string, values any) mrstorage.SqlBuilderPartFunc {
    val := reflect.ValueOf(values)

    if val.Kind() != reflect.Slice || val.Len() == 0 {
        return nil
    }

    return func (paramNumber int) (string, []any) {
        return fmt.Sprintf("%s = ANY($%d)", dbName, paramNumber), []any{pq.Array(values)}
    }
}

func (b *SqlBuilderWhere) join(separator string, conds []mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
    conds = mrstorage.SqlBuilderPartFuncRemoveNil(conds)

    if len(conds) == 0 {
        return nil
    }

    return func(paramNumber int) (string, []any) {
        var prepared []string
        var args []any

        for i := range conds {
            item, itemArgs := conds[i](paramNumber + len(args))
            prepared = append(prepared, item)
            args = mrsql.MergeArgs(args, itemArgs)
        }

        return fmt.Sprintf("(%s)", strings.Join(prepared, separator)), args
    }
}


