# Starport Plugin

TBD

## 3rd party plugin

Any developers can implement and provide plugins for Starport.
If someone has good idea for the plugin, implemented code -

1. can be build by Golang's plugin package.
2. should be public repository.
3. `manifest.json` file must be located on root directory of repository.


Single plugin repository can include multiple plugins.
To provide details about plugins included, `manifest.yml` file is required.

Sample `manifest.yml` file described below.

```
plugins:
- name: plugin-1
  path: ./plugin-1

- name: plugin-2
  path: ./plugin-2
```

- name: unique name of plugin which is included on repo.
- path: directory path where the plugin located.

Directory architectures seems like below.

At the top of repository, there is `manifest.yml` file to provide details.
And each sub directory has plugin code.

```
 .
 ├── manifest.yml
 ├── plugin-1
 │   └── plugin.go
 └── plugin-2
     └── plugin.go
```

Sample implementation of plugin.

```
package main

type plugin struct {
}

func (p *plugin) DoSomething() error {
    return nil
}

func (p *plugin) GetTransactions(blockNum int) error {
    return nil
}

// Plugin is exporting symbol.
var Plugin plugin
```

In the above, `var Plugin plugin` line is most important part because that provides entry point of
plugin.

Variable name MUST be `Plugin` but type of variable does not matter.


## Starport's configurations for plugin.

config.yml which is default config file of Starport should include `plugin` section like below.

```
plugins:
- id: demo
  repo-url: https://github.com/dsrvlabs/hackatom-plugins
```

- `id` is unique name of plugin which can be selected by Starport user.
This value is used for distinguishing plugins loaded on local machine.

- `repo-url` is repository url where plugin code uploaded.


## Commands

### List

List plugins described on `config.yml`.

```
~$ starport plugin list

Plugin list.
1: demo installed
```

### Install

```
# starport plugin install [plugin id]
~$ starport plugin install demo
```

### Running plugin

```
# starport scaffold [plugin id]/[plugin name] [function name to call] [arguments...]
~$ starport scaffold demo/plugin-1 DoSomething

~$ starport scaffold demo/plugin-1 GetTransactions 1234
```

---

## TODOs

- How to test.
- mockery
- testify

```
mockery --name Plugin
```
