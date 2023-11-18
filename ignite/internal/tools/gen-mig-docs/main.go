package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

const (
	igniteCliRepository = "http://github.com/ignite/cli.git"
	igniteBinaryPath    = "dist/ignite"
)

var scaffoldCommands = map[string][]string{
	"chain":  {"chain example --no-module --skip-git"},
	"module": {"chain example --skip-git"},
	"list": {
		"chain example --skip-git",
		"list list1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example",
		"map map1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example",
	},
	"map": {
		"chain example --skip-git",
		"map map1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example",
	},
	"single": {
		"chain example --skip-git",
		"single single1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example",
	},
	"type": {
		"chain example --skip-git",
		"type type1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example",
	},
	"message": {
		"chain example --skip-git",
		"message message1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example",
	},
	// "query": {
	// 	"chain example --skip-git",
	// 	"query query1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example",
	// },
	"packet": {
		"chain example --no-module --skip-git",
		"module example --ibc",
		"chain example --skip-git",
		"packet packet1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example",
	},
	"band": {
		"chain example --no-module --skip-git",
		"module example --ibc",
		"band band1 --module example --yes",
	},
	"vue": {
		"chain example --skip-git",
		"vue",
	},
	"react": {
		"chain example --skip-git",
		"react",
	},
}

func main() {
	var logger = log.New(os.Stdout, "", log.LstdFlags)

	fromFlag := flag.String("from", "", "Semver tag to generate migration document from")
	toFlag := flag.String("to", "", "Semver tag to generate migration document to")
	flag.Parse()

	var (
		fromVer, toVer *semver.Version
		err            error
	)
	if fromFlag != nil && *fromFlag != "" {
		fromVer, err = semver.NewVersion(*fromFlag)
		if err != nil {
			logger.Fatalf("Invalid semver tag: %s", *fromFlag)
		}
	}
	if toFlag != nil && *toFlag != "" {
		toVer, err = semver.NewVersion(*toFlag)
		if err != nil {
			logger.Fatalf("Invalid semver tag: %s", *toFlag)
		}
	}

	tmpdir, err := os.MkdirTemp("", "migdoc")
	defer os.RemoveAll(tmpdir)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println("Created temporary directory:", tmpdir)

	logger.Println("Cloning", igniteCliRepository)
	repoDir := filepath.Join(tmpdir, "src/github.com/ignite/cli")
	repo, err := git.PlainClone(repoDir, false, &git.CloneOptions{
		URL:      igniteCliRepository,
		Progress: os.Stdout,
	})
	if err != nil {
		logger.Fatalln(err)
	}

	tags, err := repo.Tags()
	if err != nil {
		logger.Fatalln(err)
	}

	var versions semver.Collection
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name()
		if name.IsTag() {
			ver, err := semver.NewVersion(name.Short())
			if err != nil {
				// Do nothing as it's not a semver tag
				return nil
			}
			versions = append(versions, ver)
		}
		return nil
	})
	if err != nil {
		logger.Fatalln(err)
	}

	if versions.Len() < 2 {
		logger.Fatalln("At least two semver tags are required")
	}

	sort.Sort(versions)

	if fromVer == nil {
		if toVer != nil {
			sort.Search(versions.Len(), func(i int) bool {
				if versions[i].LessThan(toVer) {
					fromVer = versions[i]
					return false
				}
				return true
			})
		} else {
			fromVer = versions[versions.Len()-2]
		}
	}
	if toVer == nil {
		toVer = versions[versions.Len()-1]
	}

	logger.Printf("Generating migration document for %s->%s\n\n", fromVer, toVer)

	// Checkout to previous tag and build ignite cli with make build
	logger.Printf("Checking out to %s\n", fromVer)
	wt, err := repo.Worktree()
	if err != nil {
		logger.Fatalln(err)
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(fromVer.Original()),
	})
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("Building ignite cli...")
	err = runCommand(repoDir, "make", "build")
	if err != nil {
		logger.Fatalln(err)
	}

	err = executeScaffoldCommands(logger, filepath.Join(repoDir, igniteBinaryPath), filepath.Join(tmpdir, fromVer.Original()))
	if err != nil {
		logger.Fatalln(err)
	}

	// Checkout to latest tag and build ignite cli with make build
	logger.Printf("Checking out to %s\n", toVer)
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(toVer.Original()),
	})
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("Building ignite cli...")
	err = runCommand(repoDir, "make", "build")
	if err != nil {
		logger.Fatalln(err)
	}

	err = executeScaffoldCommands(logger, filepath.Join(repoDir, igniteBinaryPath), filepath.Join(tmpdir, toVer.Original()))
	if err != nil {
		logger.Fatalln(err)
	}

	// Run diff between two directories
	logger.Println("Generating diff...")
	diffs, err := Diff(filepath.Join(tmpdir, fromVer.Original()), filepath.Join(tmpdir, toVer.Original()))
	if err != nil {
		logger.Fatalln(err)
	}
	outf, err := os.Create(fmt.Sprintf("migdoc-%s-%s.diff", fromVer, toVer))
	if err != nil {
		logger.Fatalln(err)
	}
	defer outf.Close()
	for _, diff := range diffs {
		outf.WriteString(fmt.Sprint(diff))
		outf.WriteString("\n")
	}
}

func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeScaffoldCommands(logger *log.Logger, ignitePath string, outputDir string) error {
	for name, cmds := range scaffoldCommands {
		logger.Println("Scaffolding", name)
		for _, cmd := range cmds {
			args := []string{"scaffold"}
			args = append(args, strings.Fields(cmd)...)
			args = append(args, "--path", filepath.Join(outputDir, name))
			err := runCommand("", ignitePath, args...)
			if err != nil {
				return errors.Wrapf(err, "failed to execute ignite scaffold command: %s", cmd)
			}
		}
	}
	return nil
}
