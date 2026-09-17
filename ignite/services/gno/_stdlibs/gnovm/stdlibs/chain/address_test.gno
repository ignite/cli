package chain_test

import (
	"chain"
	"testing"
)

func TestValid(t *testing.T) {
	type test struct {
		inputAddress address
		expected     bool
	}

	testCases := []test{
		{inputAddress: "g1f4v282mwyhu29afke4vq5r2xzcm6z3ftnugcnv", expected: true},
		{inputAddress: "g127jydsh6cms3lrtdenydxsckh23a8d6emqcvfa", expected: true},
		{inputAddress: "g1u7y667z64x2h7vc6fmpcprgey4ck233jaww9zq", expected: true},
		{inputAddress: "g14da4n9hcynyzz83q607uu8keuh9hwlv42ra6fa", expected: true},

		// Bech32 doesn't allow '1', 'b', 'i', 'o' for data part
		//
		// 	https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki#bech32
		{inputAddress: "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5", expected: true},
		{inputAddress: "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf1", expected: false},
		{inputAddress: "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqfb", expected: false},
		{inputAddress: "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqfi", expected: false},
		{inputAddress: "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqfo", expected: false},

		{inputAddress: "g1u0000000000000000000000000000000000000", expected: false},
		{inputAddress: "", expected: false},
		{inputAddress: "000000000000", expected: false},
		{inputAddress: "0000000000000000000000000000000000000000000000000000000000000000000000", expected: false},
	}

	for _, tc := range testCases {
		result := tc.inputAddress.IsValid()

		if result != tc.expected {
			t.Fatalf("Expected: %t, got: %t", tc.expected, result)
		}
	}
}

func TestDerivePkgAddress(t *testing.T) {
	type test struct {
		inputPath string
		expected  string
	}

	testCases := []test{
		{inputPath: "gno.land/r/gnoland/faucet", expected: "g1ttrq7mp4zy6dssnmgyyktnn4hcj3ys8xhju0n7"},
		{inputPath: "gno.land/r/demo/tamagotchi", expected: "g1a3tu874agjlkrpzt9x90xv3uzncapcn959yte4"},
	}

	for _, tc := range testCases {
		result := chain.PackageAddress(tc.inputPath)
		if result.String() != tc.expected {
			t.Fatalf("Expected: %t, got: %t", tc.expected, result)
		}
	}
}
