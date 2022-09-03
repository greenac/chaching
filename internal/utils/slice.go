package utils

func SliceContains[T int | int64 | string | float64](s []T, t T) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}

	return false
}
