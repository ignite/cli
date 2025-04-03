---
sidebar_position: 0
---

# Migration Guides

Welcome to the section on upgrading to a newer version of Ignite CLI! If you're
looking to update to the latest version, you'll want to start by checking the
documentation to see if there are any special considerations or instructions you
need to follow.

If there is no documentation for the latest version of Ignite CLI, it's
generally safe to assume that there were no breaking changes, and you can
proceed with using the latest version with your project.

## Create your own Migration Guide

The `gen-mig-diffs` tool helps developers manage and visualize code changes across multiple major versions of Ignite. With each major upgrade, the codebase might undergo significant changes, making it challenging for developers to track these differences after several updates. The `gen-mig-diffs` tool simplifies this process by scaffolding blockchains with both the old and new versions and displaying the differences.

It is located in the [Ignite CLI GitHub repository](https://github.com/ignite/cli/tree/main/ignite/internal/tools/gen-mig-diffs)
directory and has been made into a standalone project.

To set up this tool in your development environment:

```shell
gen-mig-diffs [flags]
```

This tool generates migration diff files for each of Ignite's scaffold commands. It compares two specified versions of Ignite and provides a clear, organized view of the changes.

## How to Get Started

1. Clone the Ignite CLI repository:

```shell
git clone https://github.com/ignite/cli.git --depth=1 && \
cd cli/ignite/internal/tools/gen-mig-diffs
```

2. Install and show usage:

```shell
go install . && gen-mig-diffs -h
```

### Example Migration

As an example, to generate migration diffs between versions 0.27.2 and 28.3.0, use the following command:

```shell
gen-mig-diffs --output temp/migration --from v0.27.2 --to v28.3.0
```

This command scaffolds blockchains with the specified versions and shows the differences, making it easier for developers to understand and apply necessary changes when upgrading their projects.

## Usage

```bash
This tool is used to generate migration diff files for each of ignites scaffold commands

Usage:
  gen-mig-diffs [flags]

Flags:
  -f, --from string              Version of Ignite or path to Ignite source code to generate the diff from
  -h, --help                     help for gen-mig-diffs
  -o, --output string            Output directory to save the migration document (default "docs/docs/06-migration")
      --repo-output string       Output path to clone the Ignite repository
  -s, --repo-source string       Path to Ignite source code repository. Set the source automatically set the cleanup to false
      --repo-url string          Git URL for the Ignite repository (default "https://github.com/ignite/cli.git")
      --scaffold-cache string    Path to cache directory
      --scaffold-output string   Output path to clone the Ignite repository
  -t, --to string                Version of Ignite or path to Ignite source code to generate the diff to
```
