package numbers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	separator = ","
	sepRange  = "-"
)

// ParseListRange parses comma separated numbers and range to []uint64.
func ParseListRange(arg string) ([]uint64, error) {
	result := make([]uint64, 0)
	listNumbers := make(map[uint64]struct{})
	for _, numberRange := range strings.Split(arg, separator) {
		trimmedRange := strings.TrimSpace(numberRange)
		if trimmedRange == "" {
			continue
		}

		numbers := strings.Split(trimmedRange, sepRange)
		switch len(numbers) {
		case 1:
			trimmed := strings.TrimSpace(numbers[0])
			i, err := strconv.ParseUint(trimmed, 10, 32)
			if err != nil {
				return nil, err
			}
			if _, ok := listNumbers[i]; ok {
				continue
			}
			listNumbers[i] = struct{}{}
			result = append(result, i)
		case 2:
			var (
				startN = strings.TrimSpace(numbers[0])
				endN   = strings.TrimSpace(numbers[1])
			)
			if startN == "" {
				startN = endN
			}
			if endN == "" {
				endN = startN
			}
			if startN == "" {
				continue
			}
			start, err := strconv.ParseUint(startN, 10, 32)
			if err != nil {
				return nil, err
			}
			end, err := strconv.ParseUint(endN, 10, 32)
			if err != nil {
				return nil, err
			}
			if start > end {
				start, end = end, start
			}
			for ; start <= end; start++ {
				if _, ok := listNumbers[start]; ok {
					continue
				}
				listNumbers[start] = struct{}{}
				result = append(result, start)
			}
		default:
			return nil, fmt.Errorf("cannot parse the number range: %s", trimmedRange)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result, nil
}

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

// List creates a comma separated int list with optional prefix for each uint64.
func List(numbers []uint64, prefix string) string {
	var s []string
	for _, n := range numbers {
		s = append(s, fmt.Sprintf("%s%d", prefix, n))
	}
	return strings.Join(s, ", ")
}
