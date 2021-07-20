---
order: 2
description: Install the Go programming language
---

# Install Go

In this tutorial, you will install the Go programming language (Golang) on your local computer. Follow the instructions for your operating system:

## macOS

* Download the latest MacOS installer package from <https://golang.org/dl/>.
* Open the downloaded package and follow the prompts through to completion.
* By default, the package installs the Go distribution to `/usr/local/go`, however it is always best to define the path explicitly:

    1. Open or create a `~/.bash_profile` file with your favorite command-line text editor.
    2. Add the following lines:

        ```sh
        export PATH=$PATH:$(go env GOPATH)/bin
        ```

    3. To make sure these changes execute, run the following command:

        ```sh
        source ~/.bash_profile
        ```

* Verify the installation of go by checking its version:

    ```sh
    go version
    ```

## Linux

* Download the latest Linux distribution package from <https://golang.org/dl/>.
* Extract the archive you downloaded into `/usr/local` using the following command:

```sh
sudo tar -C /usr/local -xzf go1.16.6.linux-amd64.tar.gz
```

**Note - Make sure to install the latest version available at that time.**

* To add Go to your `.bash_profile` file:

    1. Open `~/.bash_profile` file with your favorite command-line text editor.
    2. Add the following lines in your `.bash_profile` file:

        ```sh
        export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin
        ```

* To make sure the above changes are applied, run the following command:

```sh
source ~/.bash_profile
```

* Verify the installation of go using:

```sh
go version
```

## Windows Subsystem for Linux (WSL)

* Download the latest Linux distribution package from <https://golang.org/dl/>.
* Extract the archive you downloaded and move into the `/usr/local` directory by running:

```sh
sudo tar -xvf go1.16.6.linux-amd64.tar.gz && sudo mv go /usr/local
```

**Note - Make sure to install the latest version available at that time.**

* To add Go to your `.bash_profile` file:

    1. Open `~/.bash_profile` file with your favorite command-line text editor.
    2. Add the following lines in your `.bash_profile` profile:

        ```sh
        export GOPATH=$HOME/go
        export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
        ```

* To make sure the above changes are applied, run the following command:

```sh
~/.bash_profile
```

* Verify the installation of go using:

```sh
go version
```
