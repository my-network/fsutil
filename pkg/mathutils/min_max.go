package mathutils

func Min(s0 int, s ...int) int {
	min := s0
	for _, cmp := range s {
		if cmp < min {
			min = cmp
		}
	}
	return min
}
