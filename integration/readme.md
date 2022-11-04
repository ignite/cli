# Starport Integration Tests

The Starport integration tests build a new application and run all Starport commands to check the Starport code integrity. The runners and helper methods are located in this current folder. The test commands are split into folders, for better concurrency, each folder is a parallel job into the CI workflow. To create a new one, we only need to create a new folder. This will be automatically detected and added into the PR CI checks, or we can only create new tests into an existing folder or file.

Running synchronously all integration tests can be very slow. The command below can run everything:

```shell
go test -v -timeout 120m ./integration
```

Or you can just run a specific test folder, like the `list` types test

```shell
go test -v -timeout 120m ./integration/list
```

# Usage

- Create a new env and scaffold an empty chain:

```go
var (
 env  = envtest.New(t)
 path = env.Scaffold("github.com/test/blog")
)
```

- Now, you can use the env to run the starport commands and check the success status:

```go
env.Must(env.Exec("create a list with bool",
    step.NewSteps(step.New(
        step.Exec(envtest.IgniteApp, "s", "list", "--yes", "document", "signed:bool"),
        step.Workdir(path),
    )),
))
env.EnsureAppIsSteady(path)
```

- To check if the command returns an error, you can add the `envtest.ExecShouldError()` step:

```go
env.Must(env.Exec("should prevent creating a list with duplicated fields",
    step.NewSteps(step.New(
        step.Exec(envtest.IgniteApp, "s", "list", "--yes", "company", "name", "name"),
        step.Workdir(path),
    )),
    envtest.ExecShouldError(),
))
env.EnsureAppIsSteady(path)
```
