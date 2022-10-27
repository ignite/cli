package iowait

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Until waits for the appearance of s in the string n times and
// then stops blocking.
func Until(r io.Reader, s string, n int) (capturedLines []string, err error) {
	total := n
	scanner := bufio.NewScanner(r)
	for {
		if n == 0 {
			return capturedLines, nil
		}
		if !scanner.Scan() {
			if n != 0 {
				return capturedLines, fmt.Errorf("could not find %d out of %d", n, total)
			}
			return capturedLines, scanner.Err()
		}
		if strings.Contains(scanner.Text(), s) {
			capturedLines = append(capturedLines, scanner.Text())
			n--
		}
	}
}
