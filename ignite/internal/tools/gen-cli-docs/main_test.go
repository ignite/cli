package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestEscapeMDX(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "tag at line start",
			line:     "<to> is a bech32 address",
			expected: "&lt;to&gt; is a bech32 address",
		},
		{
			name:     "tag mid sentence",
			line:     "on-chain. <name> is either a bare name",
			expected: "on-chain. &lt;name&gt; is either a bare name",
		},
		{
			name:     "uppercase tag with underscore",
			line:     "address use '<FIELD_NAME>:address' to scaffold types",
			expected: "address use '&lt;FIELD_NAME&gt;:address' to scaffold types",
		},
		{
			name:     "closing tag",
			line:     "as shown </name> above",
			expected: "as shown &lt;/name&gt; above",
		},
		{
			name:     "inline code span untouched",
			line:     "inside `x/<module>/module.go` file",
			expected: "inside `x/<module>/module.go` file",
		},
		{
			name:     "tag outside span escaped, inside span untouched",
			line:     "<name> is set in `chain/<name>/config`",
			expected: "&lt;name&gt; is set in `chain/<name>/config`",
		},
		{
			name:     "unbalanced backticks left unchanged",
			line:     "use `coins or <name>",
			expected: "use `coins or <name>",
		},
		{
			name:     "no tags",
			line:     "nothing to escape here",
			expected: "nothing to escape here",
		},
		{
			name:     "arrows and comparisons untouched",
			line:     "a < b and x -> y",
			expected: "a < b and x -> y",
		},
		{
			name:     "hyphenated tag",
			line:     "<chain-id> is the chain id",
			expected: "&lt;chain-id&gt; is the chain id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeMDX(tt.line)
			assert.Equal(t, got, tt.expected)
		})
	}
}
