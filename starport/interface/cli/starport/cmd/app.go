package starportcmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gobuffalo/genny"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/app"
)

var (
	commitMessage = "Initialized with Starport"
	devXAuthor    = &object.Signature{
		Name:  "Developer Experience team at Tendermint",
		Email: "hello@tendermint.com",
		When:  time.Now(),
	}
)

// NewApp creates new command named `app` to create Cosmos scaffolds customized
// by the user given options.
func NewApp() *cobra.Command {
	c := &cobra.Command{
		Use:   "app [github.com/org/repo]",
		Short: "Generates an empty application",
		Args:  cobra.ExactArgs(1),
		RunE:  appHandler,
	}
	c.Flags().StringP("denom", "d", "token", "Token denomination")
	c.Flags().String("address-prefix", "cosmos", "Address prefix")
	return c
}

func appHandler(cmd *cobra.Command, args []string) error {
	path, err := gomodulepath.Parse(args[0])
	if err != nil {
		return err
	}
	denom, _ := cmd.Flags().GetString("denom")
	addressPrefix, _ := cmd.Flags().GetString("address-prefix")
	g, _ := app.New(&app.Options{
		ModulePath:       path.RawPath,
		AppName:          path.Package,
		BinaryNamePrefix: path.Root,
		Denom:            denom,
		AddressPrefix:    addressPrefix,
	})
	run := genny.WetRunner(context.Background())
	run.With(g)
	pwd, _ := os.Getwd()
	run.Root = pwd + "/" + path.Root
	run.Run()
	if err := initGit(path.Root); err != nil {
		return err
	}
	message := `
⭐️ Successfully created a Cosmos app '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% starport serve

NOTE: add -v flag for advanced use.
`
	fmt.Printf(message, path.Root)
	return nil
}

func initGit(path string) error {
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	if _, err := wt.Add("."); err != nil {
		return err
	}
	_, err = wt.Commit(commitMessage, &git.CommitOptions{
		All:    true,
		Author: devXAuthor,
	})
	return err
}
