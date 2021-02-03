package xstrings

// AllOrSomeFilter filters elems out from the list as they  present in filterList and
// returns the remaning ones.
// if filterList is empty, all elems from list returned.
func AllOrSomeFilter(list, filterList []string) []string {
	if len(filterList) == 0 {
		return list
	}

	var elems []string

	for _, elem := range list {
		if !SliceContains(filterList, elem) {
			elems = append(elems, elem)
		}
	}

	return elems
}

// SliceContains returns with true if s is a member of ss.
func SliceContains(ss []string, s string) bool {
	for _, e := range ss {
		if e == s {
			return true
		}
	}

	return false
}

// List returns a slice of strings captured after the value returned by do which is
// called n times.
func List(n int, do func(i int) string) []string {
	var list []string

	for i := 0; i < n; i++ {
		list = append(list, do(i))
	}

	return list
}
