package starportcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

var (
	Output    = "yaml"
	Date      = "unknown"
	Commit    = "none"
	Version   = "dev"
	Shortened = false
)

func NewVersion() *cobra.Command {

	c := &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			resp := goVersion.FuncWithOutput(Shortened, Version, Commit, Date, Output)
			fmt.Print(resp)
		},
	}
	return c
}
