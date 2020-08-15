package starportcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

var (
	Output    = "yaml"
	Date      = "unknown"
	Version   = "unset"
	Shortened = false
)

func NewVersion() *cobra.Command {

	c := &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			resp := goVersion.FuncWithOutput(Shortened, Version, Date, Output)
			fmt.Print(resp)
		},
	}
	return c
}
