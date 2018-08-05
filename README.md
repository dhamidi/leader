# Description

VIM's leader key for your terminal!  Using leader you can launch predefined commands using a short sequence of key presses instead of having to type out the whole command.

For example, using Leader you could map pressing `g` followed by `c` to running `git commit`.

# Features

- generates a keymap from your Makefile, Rakefile and package.json

# BASH integration

To trigger `leader` when pressing `\` in bash, run the following command and add it to your bash initialization file:

    bind -x '"\\":leader'

Now every time you press `\`, `leader` will be started.

# ZSH integration

To trigger `leader` when pressing `\` in zsh, run the following command and add it to your zsh initialization file:

    bindkey -s '\\' "$(which leader)"

Now every time you press `\`, `leader` will be started.
