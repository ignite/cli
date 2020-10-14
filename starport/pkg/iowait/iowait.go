package iowait

import (
	"bufio"
	"io"
	"strings"
)

// Untill waits the apperance of s in the string n times and
// then stops blocking.
func Untill(r io.Reader, s string, n int) (err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if n == 0 {
			return nil
		}
		if strings.Contains(scanner.Text(), s) {
			n--
		}
	}
	return scanner.Err()
}
