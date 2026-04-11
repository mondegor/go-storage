package mrsql

// MergeArgs - объединяет несколько массивов аргументов в один линейный массив.
// Оптимизировано: если передан один непустой массив, возвращает его без копирования.
// Если все массивы пустые или nil, возвращает nil.
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

	mergedArgs := make([]any, 0, total)

	for i := range args {
		if len(args[i]) == 0 {
			continue
		}

		mergedArgs = append(mergedArgs, args[i]...)
	}

	return mergedArgs
}
