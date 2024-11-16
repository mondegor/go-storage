package mrsql

// MergeArgs - принимает список массивов,
// формирует из них один линейный массив и возвращает его.
func MergeArgs(args ...[]any) []any {
	var total int

	for i := range args {
		total += len(args[i])
	}

	if total == 0 {
		return nil
	}

	// оптимизация, когда не требуется объединения
	for i := range args {
		if len(args[i]) == total {
			return args[i]
		}
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
