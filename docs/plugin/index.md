---
order: 1
title: Plugin System
parent:
  title: Starport Plugin System
  order: 5
description: Description of Starport plugin system.
---

# Starport Plugin

Starport plugin system provides additional environment for adding extra features.

Plugin developer can implement and publish their plugins.

Chain developers using Starport can use plugin implemented by others.

## How to use plugin system

### Step 1: Build Starport

Developer should build Starport before using Starport plugin system.
Checkout `develop` branch and build it.

```
~$ git checkout develop
~$ make build
~$ make install
```

### Step 2: Start new chain project

Before testing plugin system, you should start new chain project which is created by Starport.
If you already have it, go next section.

or create new chain project.

```
~$ starport scaffold chain github.com/cosmonaut/mars
```

Repository name might be changed.

### Step 3: Add plugin entry

Default `config.yml` file automatically created by Starport.
You need to add `plugins` section to use plugins as described below.

**config.yml**
```
accounts:
  - name: me
    coins: ["1000tokens", "10000stakes"]
  
  - name: you
    coins: ["1000tokens", "10000stakes"]

validator:
  name: dsrv
  staked: "10000000stakes"

plugins:
  - name: hello-starport-plugin
    description: Starport plugin system sample
    repository_url: "https://github.com/rootwarp/starport-plugin-sample"
```

`repository_url` is literally repository URL from Git. This repository should be public.

Because repository can contain several plugins, `name` field is required to distinguish plugins.

Developer can implement several plugins into one Git repository by separating directories.
`name` field at the above should be sub directory name of Git repository.

### Step 4: Install sample plugin.

A plugin which is appended on `config.yml` should be installed manually.

First list plugin entry.
If the plugin entry is added on `config.yml` appropriately, you can see new entry on list.
`INSTALLED` column maybe `false`.

```
~$ starport plugin list
+---+-----------------------+-----------+----------------------------------------------------+-------------------------------+
| # | NAME                  | INSTALLED | REPOSITORY URL                                     | DESCRIPTION                   |
+---+-----------------------+-----------+----------------------------------------------------+-------------------------------+
| 0 | hello-starport-plugin | false     | https://github.com/rootwarp/starport-plugin-sample | Starport plugin system sample |
+---+-----------------------+-----------+----------------------------------------------------+-------------------------------+
```

If you confirm the plugin list, try to install.

```
~$ starport plugin install hello-starport-plugin
```

After installed, you can check the plugin is installed or not.
If the plugin is installed normally, `INSTALLED` column should be `true`.

Installed plugins are stored into `$HOME/.starport/plugins`.

```
~$ starport plugin list

+---+-----------------------+-----------+----------------------------------------------------+-------------------------------+
| # | NAME                  | INSTALLED | REPOSITORY URL                                     | DESCRIPTION                   |
+---+-----------------------+-----------+----------------------------------------------------+-------------------------------+
| 0 | hello-starport-plugin | true      | https://github.com/rootwarp/starport-plugin-sample | Starport plugin system sample |
+---+-----------------------+-----------+----------------------------------------------------+-------------------------------+
```

And you can see new command `hello-starport-plugin` on `scaffold`'s sub-command.

```
~$ starport scaffold
Scaffold commands create and modify the source code files to add functionality.

CRUD stands for "create, read, update, delete".

Usage:
  starport scaffold [command]

Aliases:
  scaffold, s

Available Commands:
  chain                 Fully-featured Cosmos SDK blockchain
  module                Scaffold a Cosmos SDK module
  list                  CRUD for data stored as an array
  map                   CRUD for data stored as key-value pairs
  single                CRUD for data stored in a single location
  type                  Scaffold only a type definition
  message               Message to perform state transition on the blockchain
  query                 Query to get data from the blockchain
  packet                Message for sending an IBC packet
  band                  Scaffold an IBC BandChain query oracle to request real-time data
  vue                   Vue 3 web app template
  hello-starport-plugin Starport plugin system sample

Flags:
  -h, --help   help for scaffold

Use "starport scaffold [command] --help" for more information about a command.
```

### Step 5: Execute sample plugin.

You can list plugin's functions by running below.

```
~$ starport scaffold hello-starport-plugin 
Demo plugin

Usage:
  starport scaffold hello-starport-plugin [command]

Available Commands:
  Banner
  SayHello    Just say hello

Flags:
  -h, --help   help for hello-starport-plugin

Use "starport scaffold hello-starport-plugin [command] --help" for more information about a command.
```

And try to run one of the plugin's function.

```
~$ starport scaffold hello-starport-plugin Banner DSRV
2021/12/06 17:34:01 Initialize plugin

 ____   _____  _____  _____
|    \ |   __|| __  ||  |  |
|  |  ||__   ||    -||  |  |
|____/ |_____||__|__| \___/
```

## Plugin developer

- Directory architectures.
- Mandatory functions.
- Exporting symbol.

### Create repo and new directory

Single plugin repository can have several plugins, so these plugins should be separated by directory.

Directory name will be plugin's name.
For example, sample plugin project has two plugin, `arithmetic` and `hello-starport-plugin`.
This repository's directory structure seems like below.

```
~$ tree
.
├── arithmetic
│   ├── arithmetic.so
│   ├── go.mod
│   └── main.go
└── hello-starport-plugin
    ├── go.mod
    ├── go.sum
    ├── hello-starport-plugin.so
    ├── main.go
    └── main_test.go

2 directories, 8 files
```

### Implement mandatory functions.

There are two functions that Starport plugin should implement.

- `Init()`: Initialze function for plugin, this function will be called whenever the plugin system call plugin function.
- `Help(name string) string`: This function provide help text to Starport plugin system.
Starport will call this function with plugin function name as parameter if user call CLI with `--help` flag.

If those two functions are not implemented on plugin, Starport plugin system will return error when it try to build and load plugin.

Starport plugin project may like below.

```
package main

import (
	"log"
)

type hello struct {
}

func (p *hello) Init() {
	log.Println("Initialize plugin")
}

func (p *hello) Help(name string) string {
  return ""
}

var Plugin hello
```

Beware that variable named `Plugin` is defined at the last of the example.

When Starport plugin system loads plugin, it try to find `Plugin` variable from plugin symbol to load implemented functions.
So plugin developer MUST provide this variable for working with Starport plugin system.
