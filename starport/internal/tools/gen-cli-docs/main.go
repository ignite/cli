// this tool generates Starport CLI docs to be placed in the docs/cli dir and deployed
// on docs.starport.network.
package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	starportcmd "github.com/tendermint/starport/starport/cmd"
)

const head = `---
order: 1
description: Starport CLI docs. 
parent:
  order: 8
  title: CLI 
---

# CLI
Documentation for Starport CLI.
Note: CLI docs reflect the changes in an upcoming release.
`

func main() {
	outPath := flag.String("out", ".", ".md file path to place Starport CLI docs inside")
	flag.Parse()

	if err := generate(starportcmd.New(context.Background()), *outPath); err != nil {
		log.Fatal(err)
	}
}

func generate(cmd *cobra.Command, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, head); err != nil {
		return err
	}

	return generateCmd(cmd, f)
}

func generateCmd(cmd *cobra.Command, w io.Writer) error {
	cmd.DisableAutoGenTag = true

	b := &bytes.Buffer{}
	if err := doc.GenMarkdownCustom(cmd, b, linkHandler); err != nil {
		return err
	}

	// here we change sub titles to bold styling. Otherwise, these titles will get
	// printed in the right menu of docs.starport.network which is unpleasant because
	// we only want to see a list of all available commands without the extra noise.
	sc := bufio.NewScanner(b)
	for sc.Scan() {
		t := sc.Text()
		if strings.HasPrefix(t, "###") {
			t = strings.TrimPrefix(t, "### ")
			t = fmt.Sprintf("**%s**", t)
		}
		if _, err := fmt.Fprintln(w, t); err != nil {
			return err
		}
	}

	for _, cmd := range cmd.Commands() {
		if !cmd.IsAvailableCommand() || cmd.IsAdditionalHelpTopicCommand() {
			continue
		}

		io.WriteString(w, "\n")

		if err := generateCmd(cmd, w); err != nil {
			return err
		}
	}

	return nil
}

func linkHandler(link string) string {
	link = strings.ReplaceAll(link, "_", "-")
	link = strings.TrimSuffix(link, ".md")
	link = "#" + link
	return link
}
