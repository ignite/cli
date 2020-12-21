package numbers

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseList parses comma separated numbers to []int.
func ParseList(list string) ([]int, error) {
	ints := []int{}
	for _, number := range strings.Split(list, ",") {
		trimmed := strings.TrimSpace(number)
		if trimmed == "" {
			continue
		}
		i, err := strconv.ParseInt(trimmed, 10, 32)
		if err != nil {
			return nil, err
		}
		ints = append(ints, int(i))
	}
	return ints, nil
}

// List creates a comma separated int list with optional prefix for each int.
func List(numbers []int, prefix string) string {
	var s []string
	for _, n := range numbers {
		s = append(s, fmt.Sprintf("%s%d", prefix, n))
	}
	return strings.Join(s, ", ")
}
