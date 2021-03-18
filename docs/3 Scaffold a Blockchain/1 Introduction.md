# Scaffold a Blockchain

<!-- what is the general message we want to give? intro to the structure that is created for project and types or? we should say something about why we specify a github repo -->

The project directory of any Cosmos SDK blockchain contains many directories, source code files, configuration files, scripts, etc. Some of these files implement custom logic and are very specific to a particular project. Other files, however, are common between different Cosmos SDK projects and act as wiring between different parts of the project. Starport CLI automatically generates this common (boilerplate) code and helps in scaffolding custom functionality, so that developers can focus on application-specific logic.

One of the core features of Starport CLI is code scaffolding.

To create an entire blockchain from scratch run the following command:

```
starport app github.com/hello/planet
```

This command will create a directory called `planet`, which contains all the files for your project. The `github.com` URL in the argument is a string that will be used for Go module's path. The repository name (`planet`, in this case) will be used as the project's name. A git repository will be initialized locally.

starport type post title body

Scaffolds all the necessary files for create, read, update and delete (CRUD) actions for a specific new type. In this example, the type is `post`. The list of arguments that follow specify fields of the type, in this example: `title` and `body`. There can be any number of fields and fields can have specific types (by default fields are strings).

starport message

...

Starport CLI has replaced the deprecated `scaffold` program.

## Tutorials

[Tutorials](https://github.com/cosmos/sdk-tutorials/) help you get started with Starport and the Cosmos SDK.
