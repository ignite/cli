package scaffold

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

type (
	// Command represents a set of commands and prerequisites scaffold commands that are required to run before them.
	Command struct {
		// Name is the unique identifier of the command
		Name string
		// Prerequisites is the names of commands that need to be run before this command set
		Prerequisites []string
		// Commands is the list of scaffold commands that are going to be run
		// The commands will be prefixed with "ignite scaffold" and executed in order
		Commands []string
	}

	Scaffold struct {
		ignitePath string
		commands   Commands
	}

	Commands = map[string]Command
)

var DefaultCommands = Commands{
	"chain": Command{
		Commands: []string{"chain example --no-module"},
	},
	"module": Command{
		Prerequisites: []string{"chain"},
		Commands: []string{
			"module example --ibc",
		},
	},
	"list": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"list list1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	"map": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"map map1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --index i1:string --module example --yes",
		},
	},
	"single": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"single single1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	"type": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"type type1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	"message": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"message message1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	"query": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"query query1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints --module example --yes",
		},
	},
	"packet": Command{
		Prerequisites: []string{"module"},
		Commands: []string{
			"packet packet1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --ack f1:string,f2:strings,f3:bool,f4:int,f5:ints,f6:uint,f7:uints,f8:coin,f9:coins --module example --yes",
		},
	},
}

func New(ignitePath string, commands Commands) Scaffold {
	return Scaffold{
		ignitePath: ignitePath,
		commands:   commands,
	}
}

func (s Scaffold) Run(ver *semver.Version, out string) error {
	for _, c := range s.commands {
		if err := s.runCommand(c.Name, c.Prerequisites, c.Commands, ver, out); err != nil {
			return err
		}

		if err := applyPostScaffoldExceptions(ver, c.Name, out); err != nil {
			return err
		}
	}
	return nil
}

func (s Scaffold) runCommand(
	name string,
	prerequisites []string,
	cmds []string,
	ver *semver.Version,
	out string,
) error {
	// TODO add cache for duplicated commands.
	for _, p := range prerequisites {
		c, ok := s.commands[p]
		if !ok {
			return errors.Errorf("command %s not found", name)
		}
		if err := s.runCommand(name, c.Prerequisites, c.Commands, ver, out); err != nil {
			return err
		}
	}

	for _, cmd := range cmds {
		if err := s.executeScaffold(ver, name, cmd, out); err != nil {
			return err
		}
	}

	return nil
}

func (s Scaffold) executeScaffold(ver *semver.Version, name, cmd string, out string) error {
	args := []string{s.ignitePath, "scaffold"}
	args = append(args, strings.Fields(cmd)...)
	args = append(args, "--path", filepath.Join(out, name))
	args = applyPreExecuteExceptions(ver, args)

	if err := exec.Exec(context.Background(), args); err != nil {
		return errors.Wrapf(err, "failed to execute ignite scaffold command: %s", cmd)
	}

	return nil
}

// In this function we can manipulate command arguments before executing it in order to compensate for differences in versions.
func applyPreExecuteExceptions(ver *semver.Version, args []string) []string {
	// In versions <0.27.0, "scaffold chain" command always creates a new directory with the name of chain at the given --path
	// so we need to append "example" to the path if the command is not "chain"
	if ver.LessThan(semver.MustParse("v0.27.0")) && args[2] != "chain" {
		args[len(args)-1] = filepath.Join(args[len(args)-1], "example")
	}

	return args
}

// In this function we can manipulate the output of scaffold commands after they have been executed in order to compensate for differences in versions.
func applyPostScaffoldExceptions(ver *semver.Version, name string, out string) error {
	// In versions <0.27.0, "scaffold chain" command always creates a new directory with the name of chain at the given --path
	// so we need to move the directory to the parent directory.
	if ver.LessThan(semver.MustParse("v0.27.0")) {
		err := os.Rename(filepath.Join(out, name, "example"), filepath.Join(out, "example_tmp"))
		if err != nil {
			return errors.Wrapf(err, "failed to move %s directory to tmp directory", name)
		}

		err = os.RemoveAll(filepath.Join(out, name))
		if err != nil {
			return errors.Wrapf(err, "failed to remove %s directory", name)
		}

		err = os.Rename(filepath.Join(out, "example_tmp"), filepath.Join(out, name))
		if err != nil {
			return errors.Wrapf(err, "failed to move tmp directory to %s directory", name)
		}
	}

	return nil
}
