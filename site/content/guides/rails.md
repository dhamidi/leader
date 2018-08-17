---
Title: Ruby on Rails
---

# Ruby on Rails

When developing with Ruby on Rails there are many common commands that are invoked frequently.
`leader` can remove a lot of friction from the development process by providing quick access to many of these commands.

## Configuration

Add the following to your `~/.leaderrc`:

```json
{
  "keys": {
    "r": {
      "name": "ruby",
      "keys": {
        "d": {
          "name": "database",
          "keys": {
            "m": "bundle rake db:migrate",
            "r": "bundle rake db:reset"
          }
        },
        "s": "bin/rails server",
        "c": "bin/rails console",
        "t": "bundle exec rake test"
      }
    }
  }
}
```

Don't forget that you can add project local bindings for project-specific rake tasks by placing an additional `.leaderrc` in your project's root directory.
