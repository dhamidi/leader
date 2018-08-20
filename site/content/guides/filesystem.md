---
Title: Jumping around in the filesystem
---

# Jumping around in the filesystem

Since commands issued through `leader` are evaluated in the context of the current shell, `leader` can be used to change the shell's current working directory and other shell-local settings (such as key bindings).

This can be used for quickly accessing common directories and navigating through the filesystem with single key presses.

Here is an example configuration showing some of the possibilities:

```
{
  "keys": {
    "j": {
      "name": "jump",
      "loopingKeys": ["j", "k"],
      "keys": {
        "k": "pushd .. && pwd",
        "j": "popd && pwd",
        "\\": "cd $GOPATH/src/github.com/dhamidi/leader"
      }
    }
  }
}
```

The example above includes a shortcut to jump into the `leader`'s repository on disk to make changes there, since it's a often used directory.
The keys <kbd>k</kbd> and <kbd>j</kbd> are used to navigate up and down through the directory hierarchy.
Since both of these keys are listed as looping keys, they can be pressed repeatedly to quickly navigate through the filesystem.

This could be further extend with a helper script that lists the alphabetically next/previous directory in the current directory, which would allow for "lateral" movement between the children of the current directory and could be bound to <kbd>h</kbd> and <kbd>l</kbd>.
