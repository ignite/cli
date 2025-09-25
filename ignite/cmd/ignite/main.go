package main

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"sync"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"google.golang.org/grpc/status"

	ignitecmd "github.com/ignite/cli/v29/ignite/cmd"
	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/internal/analytics"
	"github.com/ignite/cli/v29/ignite/pkg/clictx"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xstrings"
	"github.com/ignite/cli/v29/ignite/version"
)

const exitCodeOK, exitCodeError = 0, 1

func main() {
	os.Exit(run())
}

func run() int {
	ctx := clictx.From(context.Background())
	cmd, cleanUp, err := ignitecmd.New(ctx)
	if err != nil {
		fmt.Printf("%v\n", err)
		return exitCodeError
	}
	defer cleanUp()

	// find command and send to analytics
	subCmd, _, err := cmd.Find(os.Args[1:])
	if err != nil {
		fmt.Printf("%v\n", err)
		return exitCodeError
	}
	var wg sync.WaitGroup
	analytics.SendMetric(&wg, subCmd)
	analytics.EnableSentry(ctx, &wg)

	// use charm's fang to improve CLI output
	err = fang.Execute(ctx, cmd,
		fang.WithColorSchemeFunc(cliColorScheme),
		fang.WithVersion(version.Version),
	)
	if err != nil {
		err = ensureError(err)
	}

	if errors.Is(ctx.Err(), context.Canceled) || errors.Is(err, context.Canceled) {
		fmt.Println("aborted")
		return exitCodeOK
	}

	if err != nil {
		var (
			validationErr errors.ValidationError
			versionErr    chainconfig.VersionError
			msg           string
		)

		if errors.As(err, &validationErr) {
			msg = validationErr.ValidationInfo()
		} else {
			msg = err.Error()
		}

		// Make sure the error message starts with an upper case character
		msg = xstrings.ToUpperFirst(msg)

		fmt.Printf("%s %s\n", icons.NotOK, colors.Error(msg))

		if errors.As(err, &versionErr) {
			fmt.Println("Use a more recent CLI version or upgrade blockchain app's config")
		}

		return exitCodeError
	}

	// waits for analytics to finish
	wg.Wait()

	return exitCodeOK
}

func ensureError(err error) error {
	// Extract gRPC error status.
	// These errors are returned by the plugins.
	s, ok := status.FromError(err)
	if !ok {
		// The error is not a gRPC error
		return err
	}

	// Get the error message
	cause := s.Proto().GetMessage()
	if cause == "" {
		return err
	}

	// Restore context canceled errors
	if cause == context.Canceled.Error() {
		return context.Canceled
	}

	// Use the gRPC description as error to avoid printing
	// extra gRPC error information like code or prefix.
	return errors.New(cause)
}

// cliColorScheme returns a ColorScheme for the CLI.
var cliColorScheme = func(c lipgloss.LightDarkFunc) fang.ColorScheme {
	return fang.ColorScheme{
		Base:           c(lipgloss.Color("#2F2E36"), lipgloss.Color(colors.White)),
		Title:          lipgloss.Color(colors.HiBlue),
		Codeblock:      c(lipgloss.Color("#F5F5F5"), lipgloss.Color("#2F2E36")),
		Program:        c(lipgloss.Color(colors.Blue), lipgloss.Color(colors.Cyan)),
		Command:        c(lipgloss.Color(colors.Magenta), lipgloss.Color(colors.HiBlue)),
		DimmedArgument: c(lipgloss.Color(colors.Magenta), lipgloss.Color("#AAAAAA")),
		Comment:        c(lipgloss.Color("#666666"), lipgloss.Color("#CCCCCC")),
		Flag:           c(lipgloss.Color(colors.Green), lipgloss.Color(colors.Green)),
		Argument:       c(lipgloss.Color("#2F2E36"), lipgloss.Color(colors.White)),
		Description:    c(lipgloss.Color("#2F2E36"), lipgloss.Color(colors.White)),    // flag and command descriptions
		FlagDefault:    c(lipgloss.Color(colors.Blue), lipgloss.Color(colors.HiBlue)), // flag default values in descriptions
		QuotedString:   c(lipgloss.Color(colors.Yellow), lipgloss.Color(colors.Yellow)),
		ErrorHeader: [2]color.Color{
			lipgloss.Color(colors.Yellow),
			lipgloss.Color(colors.Red),
		},
	}
}
