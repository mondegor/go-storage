package mrentity

type (
    ListPager struct {
        Index uint64 // pageIndex
        Size  uint64 // pageSize
    }

    ListSorter struct {
        FieldName string // sortField
        Direction SortDirection // sortDirection
    }
)
