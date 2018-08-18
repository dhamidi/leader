# Features

## `leader import`

**Status**: *open*

**Goals**:
- Make it easier for new users to get a useful configuration
- Integrate with existing tools to operate them more efficiently

**Example applications**:
- In an existing Ruby on Rails project, run `leader import rake` to automatically generate sensible menus from all rake tasks in the project.  Rake namespaces will be converted to key maps.  Naming collisions caused by nested namespaces/tasks are resolved by creating a nested key map keyed on the first letter of the namespace/task until no more conflicts are found.
- In a project using `npm`, run `leader import npm` to create a task for each of the scripts contained in `package.json`'s scripts section.

## Support for [fish]

**Status**: *done*

**Goals**:
- Support all major shells (bash, zsh and fish)

**Example applications**:
- A fish user wants to use `leader`.

[fish]: https://fishshell.com/

## `leader install`

**Status**: *open*

**Goals**:
- Make it easier for people to integrate leader into their shell by eliminating the need to find documentation for the user's current shell.
- Provide a hook for automated updates in case the initialization process of leader ever changes

**Example applications**:
- A new user downloads leader and just wants to get started.  They should be able to run `leader install` and have leader working in their current shell.

**Notes**:
- `leader install` is the only syntax that works the same in `bash`, `zsh` and `fish`
- `leader install` would have to:
  - add `eval "$(leader init)"` (`bash`, `zsh`) or `leader init | source` (`fish`) to the shell's initialization file
  - spawn a new process of the currently running shell to force rereading of the shell's initialization file.
- if possible, `leader install` also installs the manual page for leader in a directory where `man` can find it
