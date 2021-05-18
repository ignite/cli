package starportcmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/placeholder"
)

// FlagNoValidation disables validation.
const FlagNoValidation = "no-validation"

// RegisterValidationFlags adds validation specific flags to cobra.Command.
func RegisterValidationFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagNoValidation, false, "Disable validation.")
}

// WithValidation will configure validation based on the value of the flags.
func WithValidation(ctx context.Context, cmd *cobra.Command) context.Context {
	val, err := cmd.Flags().GetBool(FlagNoValidation)
	if err != nil {
		panic(err.Error())
	}
	if !val {
		ctx = placeholder.EnableTracing(ctx)
	}
	return ctx
}
