[![Build Status](https://travis-ci.com/dhamidi/leader.svg?branch=master)](https://travis-ci.com/dhamidi/leader)

![](./assets/logo.png)

VIM's leader key for your terminal!  Using leader you can launch predefined commands using a short sequence of key presses instead of having to type out the whole command.

For example, using Leader you could map pressing `g` followed by `c` to running `git commit`.

[![asciicast](https://asciinema.org/a/j1SIx0G9cQ5q8M8wf3SZeu5OW.png)](https://asciinema.org/a/j1SIx0G9cQ5q8M8wf3SZeu5OW)

# Installation

Download the `leader` binary from [here](https://github.com/dhamidi/leader/releases) and put it somewhere on your `$PATH`.

## bash and zsh

Add the following to your `~/.bashrc` or `~/.zshrc`:

```
eval "$(leader init)"
```

This installs leader and binds it to <kbd>\\</kbd>.

## fish

Add the following to your `~/.config/fish/config.fish`:

```
leader init | source
```

This installs leader into fish and binds it to <kbd>\\</kbd>.
Additionally it binds <kbd>Ctrl+V</kbd> to switch to a new mode in which every key is bound to `self-insert`.
This means you can press <kbd>Ctrl+V</kbd><kbd>\\</kbd> to insert a literal `\`.

# Configuration

Leader is configured through JSON files in various directories.

To get started, put the JSON listed below into `~/.leaderrc`.  This example configuration file contains shortcuts useful when developing with Golang:

```
{
  "keys": {
    "g": {
      "name": "go",
      "keys": {
        "b": "go build .",
        "t": {
          "name": "test",
          "loopingKeys": ["."],
          "keys": {
            ".": "go test .",
            "a": "go test ./..."
          }
        }
      }
    }
  }
}
```

This produces the following key bindings:

- `g b` is bound to running `go build .`
- `g t .` is bound to running `go test .`
- `g t a` is bound to running `go test ./...`

As this example shows, key maps can be nested to arbitrary depth.

A keymap's `name` is used to as a label to indicate which keymap the user is currently in when running `leader`.

New bindings can be added and removed through the `bind` subcommand as well:

```
$ leader bind --global d date
$ grep '"d"' ~/.leaderrc
  "d": "date",
$ leader @d
Sat Aug 25 21:06:04 EEST 2018
$ leader bind --unbind --global d
$ leader @d
$
```

## Looping keys

If a key occurs in the list given under a keymap's `loopingKeys` entry, this key can be pressed repeatedly to run the same command again.

## Load order

Leader tries to load a file called `.leaderrc` from your current working directory.  After trying to load that file it checks the parent directory for a `.leaderrc`, then that directory's parent directory etc until it has tried loading `$HOME/.leaderrc`.

The closer a file is to your working directory, the more important keybindings in that file are.  For example, binding `g b` to `go build .` in `~/.leaderrc` and to `gulp build` in `$HOME/projects/project-using-gulp` will make `leader` prefer running `gulp build` when in your frontend project's directory and `go build` elsewhere.

# Usage

```
leader                   # run commands in a new shell through an interactive menu
leader bind KEYS COMMAND # bind KEYS to run COMMAND in the current directory
leader print             # show interactive menu, but print commands instead of running them
leader list-keys         # list all key bindings
leader init              # print shell initialization code for $SHELL
leader help              # display leader's man page
leader version           # display the current version of leader
```

# Key bindings

The following key bindings are processed by `leader` itself and cannot be remapped:

| Key         | Function                     |
| ---         | --------                     |
| `Ctrl+C`    | Exit `leader`                |
| `Ctrl+B`    | Go back to the previous menu |
| `Up`        | Go back to the previous menu |
| `Left`      | Go back to the previous menu |
| `Backspace` | Go back to the previous menu |


# Execution environment

All commands triggered by leader are run in the context of the current shell.  This means that `cd`, `pushd` and other commands that modify the state of the current shell work without problems in your `.leaderrc`.
