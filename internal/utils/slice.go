package utils

type Number interface {
	int | int64 | float64
}

func SliceContains[T Number](s []T, t T) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}

	return false
}

func Max[T Number](s []T) T {
	return maxMin(s, true, false)
}

func Min[T Number](s []T) T {
	return maxMin(s, false, false)
}

func AbsMax[T Number](s []T) T {
	return maxMin(s, true, true)
}

func AbsMin[T Number](s []T) T {
	return maxMin(s, false, true)
}

func maxMin[T Number](s []T, isMax bool, shouldUseAbs bool) T {
	if len(s) == 0 {
		var zero T
		return zero
	}

	if len(s) == 1 {
		return s[0]
	}

	m := s[0]
	for i := 1; i < len(s); i += 1 {
		v := s[i]
		if shouldUseAbs {
			v = Abs(v)
		}

		if isMax {
			if v > m {
				m = v
			}
		} else {
			if v < m {
				m = v
			}
		}
	}

	return m
}

func MidPoint[T Number](s []T) float64 {
	if len(s) == 0 {
		return 0
	}

	var tot T
	for _, v := range s {
		tot += v
	}

	return float64(tot) / float64(len(s))
}

func Abs[T Number](n T) T {
	if n >= 0 {
		return n
	}

	return -1 * n
}
