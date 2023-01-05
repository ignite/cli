// Package sta provides access to swagger-typescript-api CLI.
package sta

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/nodetime"
)

// Option configures Generate configs.
type Option func(*configs)

// Configs holds Generate configs.
type configs struct {
	command Cmd
}

const routeNameTemplate = `<%
const { routeInfo, utils } = it;
const {
  operationId,
  method,
  route,
  moduleName,
  responsesTypes,
  description,
  tags,
  summary,
  pathArgs,
} = routeInfo;
const { _, fmtToJSDocLine, require } = utils;

const methodAliases = {
  get: (pathName, hasPathInserts) =>
    _.camelCase(` + "`" + `${pathName}_${hasPathInserts ? "detail" : "list"}` + "`" + `),
  post: (pathName, hasPathInserts) => _.camelCase(` + "`" + `${pathName}_create` + "`" + `),
  put: (pathName, hasPathInserts) => _.camelCase(` + "`" + `${pathName}_update` + "`" + `),
  patch: (pathName, hasPathInserts) => _.camelCase(` + "`" + `${pathName}_partial_update` + "`" + `),
  delete: (pathName, hasPathInserts) => _.camelCase(` + "`" + `${pathName}_delete` + "`" + `),
};

const createCustomOperationId = (method, route, moduleName) => {
  const hasPathInserts = /\{(\w){1,}\}/g.test(route);
  const splitedRouteBySlash = _.compact(_.replace(route, /\{(\w){1,}\}/g, "").split("/"));
  const routeParts = (splitedRouteBySlash.length > 1
    ? splitedRouteBySlash.splice(1)
    : splitedRouteBySlash
  ).join("_");
  return routeParts.length > 3 && methodAliases[method]
    ? methodAliases[method](routeParts, hasPathInserts)
    : _.camelCase(_.lowerCase(method) + "_" + [moduleName].join("_")) || "index";
};

if (operationId) {
	let routeName = operationId.replace('_','');	
  return routeName[0].toLowerCase() + routeName.slice(1);
}
if (route === "/")
  return _.camelCase(` + "`" + `${_.lowerCase(method)}Root` + "`" + `);

return createCustomOperationId(method, route, moduleName);
%>`

// WithCommand assigns a typescript API generator command to use for code generation.
// This allows to use a single nodetime STA generator binary in multiple code generation
// calls. Otherwise, `Generate` creates a new generator binary each time it is called.
func WithCommand(command Cmd) Option {
	return func(c *configs) {
		c.command = command
	}
}

// Cmd contains the information necessary to execute the typescript API generator command.
type Cmd struct {
	command []string
}

// Command returns the strings to execute the typescript API generator command.
func (c Cmd) Command() []string {
	return c.command
}

// Command sets the typescript API generator binary up and returns the command needed to execute it.
func Command() (command Cmd, cleanup func(), err error) {
	c, cleanup, err := nodetime.Command(nodetime.CommandSTA)
	command = Cmd{c}
	return
}

// Generate generates client code and TS types to outPath from an OpenAPI spec that resides at specPath.
func Generate(ctx context.Context, outPath, specPath string, options ...Option) error {
	c := configs{}

	for _, o := range options {
		o(&c)
	}

	command := c.command.Command()
	if command == nil {
		cmd, cleanup, err := Command()
		if err != nil {
			return err
		}

		defer cleanup()

		command = cmd.Command()
	}

	dir := filepath.Dir(outPath)
	file := filepath.Base(outPath)

	// generate temp template directory
	templateTmpPath, err := os.MkdirTemp("", "gen-js-sta-templates")
	if err != nil {
		return err
	}

	outTemplate := filepath.Join(templateTmpPath, "route-name.eta")
	err = os.WriteFile(outTemplate, []byte(routeNameTemplate), 0o644)
	if err != nil {
		return err
	}

	defer os.RemoveAll(templateTmpPath)

	// command constructs the sta command.
	command = append(command, []string{
		"--axios",
		"--module-name-index",
		"-1", // -1 removes the route namespace
		"-p",
		specPath,
		"--templates",
		templateTmpPath,
		"-o",
		dir,
		"-n",
		file,
	}...)

	// execute the command.
	return exec.Exec(ctx, command, exec.IncludeStdLogsToError())
}
