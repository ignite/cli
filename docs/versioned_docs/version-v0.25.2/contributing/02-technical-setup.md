---
sidebar_position: 1
---
# Technical setup

To ensure you have a successful experience working with our Developer Tutorials content, Ignite recommends this technical setup.

## Setting up Visual Studio Code

1. Install [Visual Studio Code](https://vscode-docs.readthedocs.io/en/latest/editor/setup).
1. Click **Extensions** in the sidebar.
1. Install this extension:
    - Go for VS Code The official Go extension for Visual Studio Code
1. When prompted:
    - `go get -v golang.org/x/tools/gopls`
    - Select `Install all` for all packages

Be sure to set up [Visual Studio Code](https://code.visualstudio.com/docs/setup/setup-overview) for your environment.

**Tip** On MacOS, install `code` in $PATH to enable [Launching Visual Studio Code from the command line](https://code.visualstudio.com/docs/setup/mac#_launching-from-the-command-line). Open the Command Palette (Cmd+Shift+P) and type 'shell command'.  

## GitHub integration

Click the GitHub icon in the sidebar for GitHub integration and follow the prompts.

## Clone the repos that you work in

- Fork or clone the <https://github.com/ignite/cli/> repository.

Internal Ignite users have different permissions. If you're not sure, fork the repo.

## Terminal tips

Master your terminal to be happy.

### iTerm2 terminal emulator

On macOS, install the [iTerm2](https://iterm2.com) OSS terminal emulator as a replacement for the default Terminal app. Installing iTerm2 as a replacement for Terminal provides an updated version of the Bash shell that supports useful features like programmable completion.

### Using ZSH as your default shell

The Z shell, also known as zsh, is a UNIX shell that is built on top of the macOS default Bourne shell.

1. If you want to set your default shell to zsh, install and set up [zsh](https://github.com/ohmyzsh/ohmyzsh/wiki/Installing-ZSH) as the default shell.

1. Install these plugins:
    - [zsh-auto-suggestions](https://github.com/zsh-users/zsh-autosuggestions/blob/master/INSTALL.md#oh-my-zsh)
    - [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting/blob/master/INSTALL.md#oh-my-zsh)

1. Edit your `~/.zshrc` file to add the plugins to load on startup:

    ```
    plugins=(
      git
      zsh-autosuggestions
      zsh-syntax-highlighting
    )
    ```

1. Log out and log back in to the terminal to use your new default zsh shell.

## Install Go

Follow the steps in [Install Ignite CLI](../guide/01-install.md) docs to install Ignite CLI and Go.
