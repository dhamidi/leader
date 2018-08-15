# Features

## `leader import`

*Goals*:
- Make it easier for new users to get a useful configuration
- Integrate with existing tools to operate them more efficiently

*Example applications*:
- In an existing Ruby on Rails project, run `leader import rake` to automatically generate sensible menus from all rake tasks in the project.  Rake namespaces will be converted to key maps.  Naming collisions caused by nested namespaces/tasks are resolved by creating a nested key map keyed on the first letter of the namespace/task until no more conflicts are found.
- In a project using `npm`, run `leader import npm` to create a task for each of the scripts contained in `package.json`'s scripts section.

## Support for [fish]

*Goals*:
- Support all major shells (bash, zsh and fish)

*Example applications*:
- A fish user wants to use `leader`.

[fish]: https://fishshell.com/
