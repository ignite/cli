---
order: 2
description: Install the Go programming language
---

# Install Go

In this tutorial, you will install the Go programming language (Golang) on your local computer. Follow the instructions for your operating system:

## macOS

* Download the latest macOS installer package from [Go downloads](https://golang.org/dl/).
* Open the downloaded package and follow the prompts to install Go.
* By default, the package installs the Go distribution to `/usr/local/go`, however it is always best to define the path explicitly:

    1. Open or create a `~/.profile` file with your favorite command-line text editor.
    2. Add the following line:

        ```sh
        export PATH=$PATH:$(go env GOPATH)/bin
        ```

    3. To make sure these changes are applied to `~/.profile`, run the following command:

        ```sh
        source ~/.profile
        ```

* Verify that you have installed Go by running the following command to check the version:

    ```sh
    go version
    ```

## Linux

* Download the latest Linux distribution package from [Go downloads](https://golang.org/dl/).
* Extract the archive you downloaded into `/usr/local` using the following command:

```sh
sudo tar -C /usr/local -xzf gox.xx.x.linux-amd64.tar.gz
```

For example, to install version 1.16.6:

```sh
sudo tar -C /usr/local -xzf go1.16.6.linux-amd64.tar.gz
```

**Note:** Make sure to install the latest supported version.

* To add Go to your `.profile` file:

    1. Open `~/.profile` file with your favorite command-line text editor and add the following line.

        ```sh
        export PATH=$PATH:$(go env GOPATH)/bin:/usr/local/go/bin
        ```

* To make sure these changes are applied to `~/.profile`, run the following command:

```sh
source ~/.profile
```

* Verify that you have installed Go by running the following command to check the version:

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

* To add Go to your `.profile` file:

    1. Open `~/.profile` file with your favorite command-line text editor and add the following line.

        ```sh
        export PATH=$PATH:$(go env GOPATH)/bin
        ```

* To make sure these changes are applied to `~/.profile`, run the following command:

```sh
source ~/.profile
```

* Verify that you have installed Go by running the following command to check the version:

```sh
go version
```
