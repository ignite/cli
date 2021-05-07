// Package integration_test integration tests Starport by using the prebuilt starport binary.
package integration_test

// spn defines the SPN version to run to test commands that interacts with SPN.
var spn = func() (repo, hash string) {
	return "https://github.com/tendermint/spn", "0cc93df25c782c67edd9abd67cced5cb689ace80"
}
