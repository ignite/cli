package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/pkg/errors"
)

const (
	igniteCliRepository = "http://github.com/ignite/cli.git"
	igniteBinaryPath    = "dist/ignite"
	igniteRepoPath      = "src/github.com/ignite/cli"
)

var scaffoldCommands = map[string][]string{
	"chain":  {"chain example --no-module"},
	"module": {"chain example"},
	"list": {
		"chain example",
		"list list1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
	},
	"map": {
		"chain example",
		"map map1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --index f10:string --module example --yes",
	},
	"single": {
		"chain example",
		"single single1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
	},
	"type": {
		"chain example",
		"type type1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
	},
	"message": {
		"chain example",
		"message message1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
	},
	"query": {
		"chain example",
		"query query1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints --module example --yes",
	},
	"packet": {
		"chain example --no-module",
		"module example --ibc",
		"packet packet1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --ack f1:string,f2:strings,f3:bool,f4:int,f5:ints,f6:uint,f7:uints,f8:coin,f9:coins --module example --yes",
	},
}

func main() {
	var logger = log.New(os.Stdout, "", log.LstdFlags)

	fromFlag := flag.String("from", "", "Semver tag to generate migration document from")
	toFlag := flag.String("to", "", "Semver tag to generate migration document to")
	sourceFlag := flag.String("source", "", "Source code directory of ignite cli repository (will be cloned if not provided)")
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

	err = run(fromVer, toVer, sourceFlag, logger)
	if err != nil {
		logger.Fatal(err)
	}
}

func run(fromVer, toVer *semver.Version, source *string, logger *log.Logger) error {
	// A temporary directory is created to clone ignite cli repository and build it
	tmpdir, err := os.MkdirTemp("", ".migdoc")
	defer os.RemoveAll(tmpdir)
	if err != nil {
		return err
	}
	logger.Println("Created temporary directory:", tmpdir)

	var (
		repoDir string
		repo    *git.Repository
	)
	if source != nil && *source != "" {
		logger.Println("Using source code directory:", *source)
		repoDir = *source
		repo, err = git.PlainOpen(*source)
		if err != nil {
			return err
		}
	} else {
		logger.Println("Cloning", igniteCliRepository)
		repoDir := filepath.Join(tmpdir, igniteRepoPath)
		repo, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			URL: igniteCliRepository,
		})
		if err != nil {
			return err
		}
	}

	versions, err := getRepositoryVersionTags(repo)
	if err != nil {
		return err
	}

	fromVer, toVer, err = validateVersionRange(fromVer, toVer, versions)

	logger.Printf("Generating migration document for %s->%s\n", fromVer, toVer)

	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}

	// Run scaffolds for fromVer and toVer
	fromVerDir := filepath.Join(tmpdir, fromVer.Original())
	err = runScaffoldsForVersion(wt, repoDir, fromVerDir, fromVer)
	if err != nil {
		return errors.Wrapf(err, "failed to run scaffolds for tag %s", fromVer)
	}
	toVerDir := filepath.Join(tmpdir, toVer.Original())
	err = runScaffoldsForVersion(wt, repoDir, toVerDir, toVer)
	if err != nil {
		return errors.Wrapf(err, "failed to run scaffolds for tag %s", toVer)
	}

	// Run diff between two directories
	logger.Println("Generating diff...")
	diffMap, err := calculateDiff(fromVerDir, toVerDir)
	if err != nil {
		return err
	}
	err = subtractBaseDiffs(diffMap)
	if err != nil {
		return err
	}

	outputDir := fmt.Sprintf("migdoc-%s-%s", fromVer, toVer)
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = saveDiffMap(diffMap, outputDir)
	if err != nil {
		return err
	}
	logger.Println("Migration document generated successfully at", outputDir)

	return nil
}

func getRepositoryVersionTags(repo *git.Repository) (semver.Collection, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tags")
	}

	// Iterate over all tags in the repository and pick valid semver tags
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
		return nil, errors.Wrap(err, "failed to iterate over tags")
	}

	sort.Sort(versions)

	return versions, nil
}

func validateVersionRange(fromVer, toVer *semver.Version, versions semver.Collection) (*semver.Version, *semver.Version, error) {
	// Unable to generate migration document if there are less than two releases!
	if versions.Len() < 2 {
		return nil, nil, errors.New("At least two semver tags are required")
	}

	// Replace fromVer and toVer with equivalent semver tags from versions
	if fromVer != nil {
		found := false
		for _, ver := range versions {
			if ver.Equal(fromVer) {
				fromVer = ver
				found = true
				break
			}
		}
		if !found {
			return nil, nil, errors.Errorf("tag %s not found", fromVer)
		}
	}
	if toVer != nil {
		found := false
		for _, ver := range versions {
			if ver.Equal(toVer) {
				toVer = ver
				found = true
				break
			}
		}
		if !found {
			return nil, nil, errors.Errorf("tag %s not found", toVer)
		}
	}

	// Picking default values for fromVer and toVer such that:
	// If both fromVer and toVer are not provided, then generate migration document for second last and last semver tags
	// If only fromVer is not provided, then use the tag before toVer as fromVer
	// If only toVer is not provided, then use the last tag as toVer
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

	// Unable to generate migration document if fromVer is greater or equal to toVer
	if fromVer.GreaterThan(toVer) || fromVer.Equal(toVer) {
		return nil, nil, errors.Errorf("from version %s should be less than to version %s", fromVer, toVer)
	}

	return fromVer, toVer, nil
}

// Run scaffolds commands one by one with the given version of ignite cli and save the output in the output directory
func runScaffoldsForVersion(wt *git.Worktree, repoDir, outputDir string, ver *semver.Version) error {
	err := checkoutAndBuildIgniteCli(wt, ver.Original(), repoDir)
	if err != nil {
		return err
	}

	binPath := filepath.Join(repoDir, igniteBinaryPath)
	err = executeScaffoldCommands(binPath, outputDir, ver)
	if err != nil {
		return err
	}

	err = applyVersionExceptions(outputDir, ver)
	if err != nil {
		return err
	}

	return nil
}

func checkoutAndBuildIgniteCli(wt *git.Worktree, tag, repoDir string) error {
	err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(tag),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to checkout tag %s", tag)
	}

	err = exec.Exec(context.Background(), []string{"make", "build"}, exec.StepOption(step.Workdir(repoDir)))
	if err != nil {
		return errors.Wrap(err, "failed to build ignite cli using make build")
	}

	return nil
}

func executeScaffoldCommands(ignitePath, outputDir string, ver *semver.Version) error {
	for name, cmds := range scaffoldCommands {
		for _, cmd := range cmds {
			args := []string{ignitePath, "scaffold"}
			args = append(args, strings.Fields(cmd)...)
			pathFlag := filepath.Join(outputDir, name)
			if !strings.HasPrefix(cmd, "chain") && ver.LessThan(semver.MustParse("v0.27.0")) {
				pathFlag = filepath.Join(outputDir, name, "example")
			}
			args = append(args, "--path", pathFlag)
			err := exec.Exec(context.Background(), args, exec.StepOption(step.Stdout(os.Stdout)), exec.StepOption(step.Stderr(os.Stderr)))
			if err != nil {
				return errors.Wrapf(err, "failed to execute ignite scaffold command: %s", cmd)
			}
		}

	}
	return nil
}

func applyVersionExceptions(outputDir string, ver *semver.Version) error {
	if ver.LessThan(semver.MustParse("v0.27.0")) {
		// Move files from the "example" directory to the parent directory for each scaffold directory
		for name := range scaffoldCommands {
			err := os.Rename(filepath.Join(outputDir, name, "example"), filepath.Join(outputDir, "example_tmp"))
			if err != nil {
				return errors.Wrapf(err, "failed to move %s directory to tmp directory", name)
			}

			err = os.RemoveAll(filepath.Join(outputDir, name))
			if err != nil {
				return errors.Wrapf(err, "failed to remove %s directory", name)
			}

			err = os.Rename(filepath.Join(outputDir, "example_tmp"), filepath.Join(outputDir, name))
			if err != nil {
				return errors.Wrapf(err, "failed to move tmp directory to %s directory", name)
			}
		}
	}

	return nil
}

func calculateDiff(fromVerDir, toVerDir string) (map[string][]gotextdiff.Unified, error) {
	diffMap := make(map[string][]gotextdiff.Unified)
	for name := range scaffoldCommands {
		diffs, err := diff(filepath.Join(fromVerDir, name), filepath.Join(toVerDir, name))
		if err != nil {
			return nil, err
		}
		diffMap[name] = diffs
	}

	return diffMap, nil
}

func diff(dir1, dir2 string) ([]gotextdiff.Unified, error) {
	var (
		diffs  []gotextdiff.Unified
		marked map[string]struct{} = make(map[string]struct{})
	)

	dirsLCP := longestCommonPrefix(dir1, dir2)

	// Consider dir1 as reference and walk through all of the files comparing them with files in dir2.
	err := filepath.Walk(dir1, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if isException(path) {
			return nil
		}

		relPath, err := filepath.Rel(dir1, path)
		if err != nil {
			return err
		}
		marked[relPath] = struct{}{}

		file1, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		str1 := string(file1)

		file2, err := os.ReadFile(filepath.Join(dir2, relPath))
		if !os.IsNotExist(err) && err != nil {
			return err
		}
		str2 := string(file2)

		edits := myers.ComputeEdits(span.URIFromPath(relPath), str1, str2)
		if len(edits) > 0 {
			fromPath, err := filepath.Rel(dirsLCP, path)
			if err != nil {
				panic(err)
			}
			toPath, err := filepath.Rel(dirsLCP, filepath.Join(dir2, relPath))
			if err != nil {
				panic(err)
			}
			diffs = append(diffs, gotextdiff.ToUnified(fromPath, toPath, str1, edits))
		}
		return nil
	})
	if err != nil {
		return diffs, err
	}

	// Walk through all of the files in dir2 that were not compared with files in dir1.
	err = filepath.Walk(dir2, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if isException(path) {
			return nil
		}

		relPath, err := filepath.Rel(dir2, path)
		if err != nil {
			return err
		}
		if _, ok := marked[relPath]; ok {
			return nil
		}

		str1 := ""

		file2, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		str2 := string(file2)

		edits := myers.ComputeEdits(span.URIFromPath(relPath), str1, str2)
		fromPath, err := filepath.Rel(dirsLCP, filepath.Join(dir1, relPath))
		if err != nil {
			panic(err)
		}
		toPath, err := filepath.Rel(dirsLCP, path)
		if err != nil {
			panic(err)
		}
		diffs = append(diffs, gotextdiff.ToUnified(fromPath, toPath, str1, edits))
		return nil
	})
	return diffs, nil
}

func longestCommonPrefix(strs ...string) string {
	longest := strings.Split(strs[0], string(filepath.Separator))

	cmp := func(a []string) {
		if len(a) < len(longest) {
			longest = longest[:len(a)]
		}
		for i := 0; i < len(longest); i++ {
			if a[i] != longest[i] {
				longest = longest[:i]
				return
			}
		}
	}

	for i := 1; i < len(strs); i++ {
		r := strings.Split(strs[i], string(filepath.Separator))
		cmp(r)
	}
	return "/" + filepath.Join(longest...)
}

func subtractBaseDiffs(diffMap map[string][]gotextdiff.Unified) error {
	// Remove chain and module diffs from other diffs
	chainDiffs := diffMap["chain"]
	moduleDiffs := diffMap["module"]
	for name, diffs := range diffMap {
		if name == "module" {
			diffs = subtractDiffs(diffs, chainDiffs)
		} else if name != "chain" {
			diffs = subtractDiffs(diffs, moduleDiffs)
		}
		diffMap[name] = diffs
	}

	return nil
}

func subtractDiffs(src []gotextdiff.Unified, base []gotextdiff.Unified) []gotextdiff.Unified {
	dst := make([]gotextdiff.Unified, 0, len(src))
	for i := 0; i < len(src); i++ {
		edited := false
		for j := 0; j < len(base); j++ {
			if equalScaffoldPaths(src[i].From, base[j].From) && equalScaffoldPaths(src[i].To, base[j].To) {
				if hs := subtractHunks(src[i].Hunks, base[j].Hunks); len(hs) > 0 {
					dst = append(dst, gotextdiff.Unified{
						From:  src[i].From,
						To:    src[i].To,
						Hunks: subtractHunks(src[i].Hunks, base[j].Hunks),
					})
				}
				edited = true
			}
		}

		if !edited {
			dst = append(dst, src[i])
		}
	}

	return dst
}

func subtractHunks(src []*gotextdiff.Hunk, base []*gotextdiff.Hunk) []*gotextdiff.Hunk {
	dst := make([]*gotextdiff.Hunk, 0, len(src))
	for i := 0; i < len(src); i++ {
		edited := false
		for j := 0; j < len(base); j++ {
			if src[i].FromLine <= base[j].FromLine && src[i].ToLine >= base[j].ToLine {
				if h := subtractHunk(src[i], base[j]); h != nil {
					dst = append(dst, h)
				}
				edited = true
			}
		}

		if !edited {
			dst = append(dst, src[i])
		}
	}

	return dst
}

func subtractHunk(src, base *gotextdiff.Hunk) *gotextdiff.Hunk {
	newLines := make([]gotextdiff.Line, 0, len(src.Lines))
	equals := 0
	for i := 0; i < len(src.Lines); i++ {
		rep := false
		for j := 0; j < len(base.Lines); j++ {
			if src.Lines[i].Kind != gotextdiff.Equal && src.Lines[i].Kind == base.Lines[j].Kind && src.Lines[i].Content == base.Lines[j].Content {
				rep = true
				break
			}
		}

		if !rep {
			newLines = append(newLines, src.Lines[i])
		}

		if src.Lines[i].Kind == gotextdiff.Equal {
			equals++
		}
	}

	// If all the lines in the hunk are equal or there's no line left, then return nil
	if equals == len(newLines) || len(newLines) == 0 {
		return nil
	}

	return &gotextdiff.Hunk{
		FromLine: src.FromLine,
		ToLine:   src.ToLine,
		Lines:    newLines,
	}
}

func equalScaffoldPaths(a, b string) bool {
	// Remove the first two directories from the path (version/scaffold_type) and compare the rest
	a = strings.Join(strings.Split(a, string(filepath.Separator))[2:], string(filepath.Separator))
	b = strings.Join(strings.Split(b, string(filepath.Separator))[2:], string(filepath.Separator))

	return a == b
}

func saveDiffMap(diffMap map[string][]gotextdiff.Unified, outputPath string) error {
	for name, diffs := range diffMap {
		outf, err := os.Create(filepath.Join(outputPath, name+".diff"))
		if err != nil {
			return err
		}
		defer outf.Close()
		for _, diff := range diffs {
			outf.WriteString(fmt.Sprint(diff))
			outf.WriteString("\n")
		}
	}

	return nil
}
