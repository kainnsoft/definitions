package toolkit

// убирает из первого массива все значения, которые содержатся во втором
func ArrayInt64Diff(a1, a2 []int64) (res []int64) {
	var inc bool
	res = make([]int64, 0, len(a1))
	for _, i1 := range a1 {
		inc = true
		for _, i2 := range a2 {
			if i1 == i2 {
				inc = false
				break
			}
		}
		if inc {
			res = append(res, i1)
		}
	}
	return
}

func ArrayInt64Unique(input []int64) []int64 {
	u := make([]int64, 0, len(input))
	m := make(map[int64]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		} else {
			m[val] = false
		}
	}

	return u
}
