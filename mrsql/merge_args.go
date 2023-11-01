package mrsql

func MergeArgs(args ...[]any) []any {
    var total int

    for i := range args {
        total += len(args[i])
    }

    mergedArgs := make([]any, total)
    n := 0

    for i := range args {
        for j := range args[i] {
            mergedArgs[n] = args[i][j]
            n++
        }
    }

    return mergedArgs
}
