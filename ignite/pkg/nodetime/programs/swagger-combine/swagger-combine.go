package swaggercombine

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"regexp"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/nodetime"
)

// Config represent swagger-combine config.
type Config struct {
	Swagger string `json:"swagger"`
	Info    Info   `json:"info"`
	APIs    []API  `json:"apis"`
}

type Info struct {
	Title       string `json:"title"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type API struct {
	ID           string       `json:"-"`
	URL          string       `json:"url"`
	OperationIDs OperationIDs `json:"operationIds"`
	Dereference  struct {
		Circular string `json:"circular"`
	} `json:"dereference"`
}

type OperationIDs struct {
	Rename map[string]string `json:"rename"`
}

var opReg = regexp.MustCompile(`(?m)operationId.+?(\w+)`)

// AddSpec adds a new OpenAPI spec to Config by path in the fs and unique id of spec.
func (c *Config) AddSpec(id, path string) error {
	// make operationId fields unique.
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	ops := opReg.FindAllStringSubmatch(string(content), -1)
	rename := make(map[string]string, len(ops))

	for _, op := range ops {
		o := op[1]
		rename[o] = id + o
	}

	// add api with replaced operation ids.
	c.APIs = append(c.APIs, API{
		ID:           id,
		URL:          path,
		OperationIDs: OperationIDs{Rename: rename},
		// Added due to https://github.com/maxdome/swagger-combine/pull/110 after enabling more services to be generated in #835
		Dereference: struct {
			Circular string `json:"circular"`
		}(struct{ Circular string }{Circular: "ignore"}),
	})

	return nil
}

// Combine combines openapi specs into one and saves to out path.
// specs is a spec id-fs path pair.
func Combine(ctx context.Context, c Config, out string) error {
	command, cleanup, err := nodetime.Command(nodetime.CommandSwaggerCombine)
	if err != nil {
		return err
	}
	defer cleanup()

	f, err := os.CreateTemp("", "*.json")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if err := json.NewEncoder(f).Encode(c); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	command = append(command, []string{
		f.Name(),
		"-o", out,
		"-f", "yaml",
		"--continueOnConflictingPaths", "true",
		"--includeDefinitions", "true",
	}...)

	// execute the command.
	return exec.Exec(ctx, command, exec.IncludeStdLogsToError())
}
