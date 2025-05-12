<div align="center">
  <h1> Generate Ignite Migration Diffs </h1>
</div>

This repository hosts the Chain Scaffold Migration Tool for Ignite CLI, designed to help developers migrate their
projects from older versions of Ignite to the latest release.
This tool addresses compatibility and feature alignment as a detailed
in [Issue #3699](https://github.com/ignite/cli/issues/3699) and implemented
in [PR #3718](https://github.com/ignite/cli/pull/3718).

The migration tool aims to streamline the update process for projects built with Ignite CLI, ensuring they leverage the
latest improvements and SDK stack.

#### Features

- Automated migration of chain scaffold files.
- Detailed comparison and generation of migration differences.
- Support for multiple versions of chain scaffolds.

## Installation

It is located in the `ignite/internal/tools/gen-mig-diffs`
directory and made it a standalone project.

To set up this tool in your development environment:

1. Clone the Ignite CLI repository:

```shell
git clone https://github.com/ignite/cli.git && \
cd cli/ignite/internal/tools/gen-mig-diffs
```

2. Install and show usage:

```shell
go install . && gen-mig-diffs -h
```

3. Run migration diff tool:

```shell
gen-mig-diffs --output temp/migs --from v0.27.2 --to v28.3.0
```

4. In case of the issue `unable to authenticate, attempted methods [none publickey], no supported methods remain`.
   Make sure you have SSH keys set up for GitHub. If yes, try to add the SSH key to your SSH agent:

```shell
chmod 600 ~/.ssh/id_rsa
ssh-add ~/.ssh/id_rsa
```

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