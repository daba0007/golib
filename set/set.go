package set

// CompareStringSlice compare two string slice
func CompareStringSlice(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	if (s1 == nil) != (s2 == nil) {
		return false
	}
	hash := make(map[string]struct{})
	for _, v := range s1 {
		hash[v] = struct{}{}
	}
	for _, v := range s2 {
		// 如果找不到，说明两个切片不相等
		if _, ok := hash[v]; !ok {
			return false
		}
	}
	return true
}

// CompareIntSlice compare two int slice
func CompareIntSlice(i1 []int, i2 []int) bool {
	if len(i1) != len(i2) {
		return false
	}
	if (i1 == nil) != (i2 == nil) {
		return false
	}
	hash := make(map[int]struct{})
	for _, v := range i1 {
		hash[v] = struct{}{}
	}
	for _, v := range i2 {
		// 如果找不到，说明两个切片不相等
		if _, ok := hash[v]; !ok {
			return false
		}
	}
	return true
}
