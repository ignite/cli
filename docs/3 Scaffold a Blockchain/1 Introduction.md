# Scaffold a Blockchain

The project directory of any Cosmos SDK blockchain contains many directories, source code files, configuration files, scripts, etc. Some of these files implement custom logic and are very specific to a particular project. Other files, however, are common between different Cosmos SDK projects and act as wiring between different parts of the project. Starport CLI automatically generates this common (boilerplate) code and helps in scaffolding custom functionality, so that developers can focus on application-specific logic.

## Code Scaffolding

One of the core features of Starport CLI is code scaffolding.

To create a blockchain from scratch, run the following command:

```
starport app github.com/hello/planet
```

This command creates a directory called `planet` that contains all of the files for your project.

- _github.com/hello/planet_ is used for the Go module path. A git repository is initialized locally.

- _planet_ in the string is the repository name that defines the project name.

## Generate Code for New Types

```
starport type post title body
```

The `starport type` command scaffolds all the necessary files for create, read, update and delete (CRUD) actions for a specific new type.

In this example, the type is `post`. The list of arguments that follow specify fields of the type, in this example: `title` and `body`. There can be any number of fields and fields can have specific types. By default, fields are strings.

Starport CLI has replaced the deprecated `scaffold` program.

## Step-by-Step Learning

For learn new skills by using step-by-step processes, head over to the [Cosmos SDK Tutorials](https://tutorials.cosmos.network/). These tutorials get you started with Starport and the Cosmos SDK.
