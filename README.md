[![Build Status](https://travis-ci.com/dhamidi/leader.svg?branch=master)](https://travis-ci.com/dhamidi/leader)

# Description

![](./assets/logo.png)

VIM's leader key for your terminal!  Using leader you can launch predefined commands using a short sequence of key presses instead of having to type out the whole command.

For example, using Leader you could map pressing `g` followed by `c` to running `git commit`.

# Features

- generates a keymap from your Makefile, Rakefile and package.json

# Key bindings

The following key bindings are processed by `leader` itself and cannot be remapped:

| Key         | Function                     |
| ---         | --------                     |
| `Ctrl+C`    | Exit `leader`                |
| `Ctrl+B`    | Go back to the previous menu |
| `Up`        | Go back to the previous menu |
| `Left`      | Go back to the previous menu |
| `Backspace` | Go back to the previous menu |

# Configuration

Here is an example configuration file, containing shortcuts useful when developing with Golang:

```
{
  "keys": {
    "g": {
      "name": "go",
      "keys": {
        "b": "go build ."
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


## Looping keys

If a key occurs in the list given under a keymap's `loopingKeys` entry, this key can be pressed repeatedly to rnu the same command again.

## Load order

Leader tries to load a file called `.leaderrc` from your current working directory.  After trying to load that file it checks the parent directory for a `.leaderrc`, then that directory's parent directory etc until it has tried loading `$HOME/.leaderrc`.

The closer a file is to your working directory, the more important keybindings in that file are.  For example, binding `g b` to `go build .` in `~/.leaderrc` and to `gulp build` in `$HOME/projects/project-using-gulp` will make `leader` prefer running `gulp build` when in your frontend project's directory and `go build` elsewhere.


# Installation

Download the `leader` binary from [here](https://github.com/dhamidi/leader/releases) and put it somewhere on your `$PATH`.

Add the following to your `~/.bashrc` or `~/.zshrc`:

```
eval "$(leader init)"
```

This installs leader and binds it to `\`.

# Execution environment

All commands triggered by leader are run in the context of the current shell.  This means that `cd`, `pushd` and other commands that modify the state of the current shell work without problems in your `.leaderrc`.
