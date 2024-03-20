package scaffold

import "github.com/ignite/cli/v29/ignite/pkg/errors"

type (
	// Command represents a set of command and prerequisites scaffold command that are required to run before them.
	Command struct {
		// Name is the unique identifier of the command
		Name string
		// Prerequisite is the name of command that need to be run before this command set
		Prerequisite string
		// Commands is the list of scaffold command that are going to be run
		// The command will be prefixed with "ignite scaffold" and executed in order
		Commands []string
	}

	Commands []Command
)

func (c Commands) Get(name string) (Command, error) {
	for _, cmd := range c {
		if cmd.Name == name {
			return cmd, nil
		}
	}
	return Command{}, errors.Errorf("command %s not exist", name)
}

func (c Commands) Has(name string) bool {
	for _, cmd := range c {
		if cmd.Name == name {
			return true
		}
	}
	return false
}

func (c Commands) Validate() error {
	cmdMap := make(map[string]bool)
	for i, command := range c {
		if cmdMap[command.Name] {
			return errors.Errorf("duplicate command name found: %s", command.Name)
		}
		cmdMap[command.Name] = true
		if command.Name == "" {
			return errors.Errorf("empty command name at index %d: %v", i, command)
		}
		if len(command.Commands) == 0 {
			return errors.Errorf("empty command list at index %d: %v", i, command)
		}
	}
	for _, command := range c {
		if command.Prerequisite != "" && !cmdMap[command.Prerequisite] {
			return errors.Errorf("command %s pre-requisete %s not found", command.Name, command.Prerequisite)
		}
	}
	return nil
}

var defaultCommands = Commands{
	Command{
		Name:     "chain",
		Commands: []string{"chain example --no-module"},
	},
	Command{
		Name:         "module",
		Prerequisite: "chain",
		Commands:     []string{"module example --ibc"},
	},
	Command{
		Name:         "list",
		Prerequisite: "module",
		Commands: []string{
			"list list1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	Command{
		Name:         "map",
		Prerequisite: "module",
		Commands: []string{
			"map map1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --index i1:string --module example --yes",
		},
	},
	Command{
		Name:         "single",
		Prerequisite: "module",
		Commands: []string{
			"single single1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	Command{
		Name:         "type",
		Prerequisite: "module",
		Commands: []string{
			"type type1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	Command{
		Name:         "message",
		Prerequisite: "module",
		Commands: []string{
			"message message1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --module example --yes",
		},
	},
	Command{
		Name:         "query",
		Prerequisite: "module",
		Commands: []string{
			"query query1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints --module example --yes",
		},
	},
	Command{
		Name:         "packet",
		Prerequisite: "module",
		Commands: []string{
			"packet packet1 f1:string f2:strings f3:bool f4:int f5:ints f6:uint f7:uints f8:coin f9:coins --ack f1:string,f2:strings,f3:bool,f4:int,f5:ints,f6:uint,f7:uints,f8:coin,f9:coins --module example --yes",
		},
	},
}
