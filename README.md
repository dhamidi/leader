# Description

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
  "bindings": {
    "q": ["<quit>"],
    "g": {
      "name": "go",
      "keys": {
        "b": ["go", "build", "."],
        "t": {
          "name": "test",
          "loopingKeys": ["."],
          "keys": {
            ".": ["go", "test", "."],
            "a": ["go", "test", "./..."]
          }
        }
      }
    }
  }
}
```

This produces the following key bindings:

- `q` is bound to the builtin command `quit`.  The `<` and `>` mark the command as a builtin command.
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

# BASH integration

To trigger `leader` when pressing `\` in bash, run the following command and add it to your bash initialization file:

    bind -x '"\\":leader'

Now every time you press `\`, `leader` will be started.

# ZSH integration

To trigger `leader` when pressing `\` in zsh, run the following command and add it to your zsh initialization file:

    bindkey -s '\\' "$(which leader)"

Now every time you press `\`, `leader` will be started.
