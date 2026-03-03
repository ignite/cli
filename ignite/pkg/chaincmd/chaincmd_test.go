package chaincmd

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
)

func TestInitCommandBuildsExpectedCommand(t *testing.T) {
	cmd := New("simd", WithChainID("my-chain"), WithHome("/tmp/simd"))

	s := step.New(cmd.InitCommand("my-moniker"))

	require.Equal(t, "simd", s.Exec.Command)
	require.Equal(t, []string{
		"init",
		"my-moniker",
		"--chain-id",
		"my-chain",
		"--home",
		"/tmp/simd",
	}, s.Exec.Args)
}

func TestAddKeyCommandAddsOptionalFieldsAndKeyringBackend(t *testing.T) {
	cmd := New(
		"simd",
		WithKeyringBackend(KeyringBackendTest),
		WithHome("/tmp/simd"),
	)

	s := step.New(cmd.AddKeyCommand("alice", "118", "3", "1"))

	require.Equal(t, "simd", s.Exec.Command)
	require.Equal(t, []string{
		"keys",
		"add",
		"alice",
		"--output",
		"json",
		"--coin-type",
		"118",
		"--account",
		"3",
		"--index",
		"1",
		"--keyring-backend",
		"test",
		"--home",
		"/tmp/simd",
	}, s.Exec.Args)
}

func TestAddKeyCommandSkipsEmptyOptionalFields(t *testing.T) {
	cmd := New("simd")
	s := step.New(cmd.AddKeyCommand("alice", "", "", ""))

	require.Equal(t, "simd", s.Exec.Command)
	require.Equal(t, []string{
		"keys",
		"add",
		"alice",
		"--output",
		"json",
	}, s.Exec.Args)
}

func TestStatusCommandAddsNodeFlag(t *testing.T) {
	cmd := New(
		"simd",
		WithNodeAddress("http://127.0.0.1:26657"),
		WithHome("/tmp/simd"),
	)

	s := step.New(cmd.StatusCommand())

	require.Equal(t, "simd", s.Exec.Command)
	require.Equal(t, []string{
		"status",
		"--node",
		"http://127.0.0.1:26657",
		"--home",
		"/tmp/simd",
	}, s.Exec.Args)
}

func TestCopyOverridesOptionsWithoutMutatingOriginal(t *testing.T) {
	original := New("simd", WithChainID("chain-A"))
	copied := original.Copy(WithChainID("chain-B"))

	originalStep := step.New(original.InitCommand("alice"))
	copiedStep := step.New(copied.InitCommand("alice"))

	require.Equal(t, []string{
		"init",
		"alice",
		"--chain-id",
		"chain-A",
	}, originalStep.Exec.Args)
	require.Equal(t, []string{
		"init",
		"alice",
		"--chain-id",
		"chain-B",
	}, copiedStep.Exec.Args)
}

func TestKeyringBackendFromString(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  KeyringBackend
		shouldErr bool
	}{
		{
			name:     "unspecified",
			input:    "",
			expected: KeyringBackendUnspecified,
		},
		{
			name:     "os",
			input:    "os",
			expected: KeyringBackendOS,
		},
		{
			name:     "file",
			input:    "file",
			expected: KeyringBackendFile,
		},
		{
			name:     "pass",
			input:    "pass",
			expected: KeyringBackendPass,
		},
		{
			name:     "test",
			input:    "test",
			expected: KeyringBackendTest,
		},
		{
			name:     "kwallet",
			input:    "kwallet",
			expected: KeyringBackendKwallet,
		},
		{
			name:      "invalid",
			input:     "invalid",
			expected:  KeyringBackendUnspecified,
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := KeyringBackendFromString(tc.input)
			require.Equal(t, tc.expected, got)
			if tc.shouldErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
