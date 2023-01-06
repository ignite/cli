package protoutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	protoPath = "../testdata/liquidity"
)

// Sanity check: Ensure that parsing works fine.
func TestParseSuccess(t *testing.T) {
	files := []string{"genesis", "liquidity", "msg", "query", "tx"}
	for _, file := range files {
		file = fmt.Sprintf(`../testdata/liquidity/%[1]v.proto`, file)
		_, err := ParseProtoPath(file)
		require.NoError(t, err)
	}

	// Cover the error case 1) -- non existent file:
	_, err := ParseProtoPath("p.proto")
	require.Error(t, err)
	// Cover the error case 2) -- invalid file type
	_, err = ParseProtoPath("parser.go")
	require.Error(t, err)
}

// Sanity check: Ensure that parsing works fine with a string input.
func TestParseString(t *testing.T) {
	_, err := parseStringProto(`syntax = "proto3";
	
	package test;
	import "github.com/cosmos/cosmos-sdk/codec";
	import "github.com/cosmos/cosmos-sdk/codec/types";

	message Msg {
		string name = 1;
		string description = 2;
	}`)
	require.NoError(t, err)

	// Cover the error case.
	_, err = parseStringProto(`var b = "go"`)
	require.Error(t, err)
}

func TestParseProtoFiles(t *testing.T) {
	files := []string{"genesis", "liquidity", "msg", "query", "tx"}
	for _, f := range files {
		f = fmt.Sprintf(`%[1]v/%[2]v.proto`, protoPath, f)
		fp, err := os.Open(f)
		require.NoError(t, err)

		nodes, err := ParseProtoFile(fp)
		require.NoError(t, err)

		// Pass through printer and check that it still parses
		// afterwards:
		out := Print(nodes)
		_, err = parseStringProto(out)
		require.NoError(t, err)
	}
}
