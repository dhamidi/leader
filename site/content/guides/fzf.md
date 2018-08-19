---
Title: Interactive selections with fzf
---

# Interactive selections with fzf

Sometimes commands that you bind to keys with leader need some extra information,
for example which git branch to check out or which manual page to open.

Most of the times that extra bit of information comes from a set of existing items:
the list of branches in your git repository or the list of manual pages installed on your system.

[fzf](https://github.com/junegunn/fzf) is a command line utility for presenting interactive selections in the terminal.
This guide examines a few possibilities of using `fzf` together with `leader`.

**Prerequisites**:

- have `fzf` installed
- have a directory on your `PATH` in which you can place scripts

## Changing git branches

Place the following script as `git-select-branch` in a directory on your `PATH`:

```
#!/bin/bash

BRANCH_NAME=$(git branch --all --list --format='%(refname:lstrip=2)' | fzf --height=10)
git checkout "$BRANCH_NAME"
```

Then add the following binding to your `~/.leaderrc`, whereever you already have bindings for `git`:

```
    "b": "git select-branch"
```

The end result:

[![asciicast](https://asciinema.org/a/O0NB8Qle6YkPuIHTGJBYDOr0x.png)](https://asciinema.org/a/O0NB8Qle6YkPuIHTGJBYDOr0x)

## Quickly opening man pages

Place the following script as `select-man-page` in a directory on your `PATH`.  The scripts uses the output of `apropos(1)` to show a summary of manual pages installed on your system
and then presents you with a selection of all those man pages.

```
#!/bin/bash
apropos . |
fzf --height=10 |
tr -d '()' |
awk '{print $2 " " $1}' |
xargs man
```

Then add the following binding to your `~/.leaderrc`:

```
  "m": "select-man-page"
```

The end result:

[![asciicast](https://asciinema.org/a/lcIaOUnazTMc0egXiFq6CESHD.png)](https://asciinema.org/a/lcIaOUnazTMc0egXiFq6CESHD)
