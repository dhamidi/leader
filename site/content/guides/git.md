---
Title: Git
---

# Integrating Git

Using `git` from the command line can be quite tedious, especially when frequently having to run the same commands.
Users of `emacs` often use [magit] and the setup presented here can best be thought of as an approximation of magit's functionality outside of `emacs`.

## Configuration

Add the following to your `~/.leaderrc`:

```json
{
  "keys": {
    "g": {
      "name": "git",
      "loopingKeys": ["s"],
      "keys": {
        "p": "git pull",
        "P": "git push",
        "c": "git commit -m",
        "s": "git status",
        "a": "git add -p .",
        "f": "git fetch",
        "b": "git select-branch",
        "B": "read -e -p 'Create branch: '; git checkout -b \"$REPLY\"",
        "g": "read -e -p 'Pattern: '; git grep \"$REPLY\"",
        "z": {
          "name": "stash",
          "keys": {
            "z": "git stash save",
            "p": "git stash pop"
          }
        }
      }
    }
  }
}
```

Most of the bindings should be obvious.
A few bindings might need some explaining words:

- <kbd>g B</kbd> prompts for a branch name and then creates a branch with the provided name
- <kbd>g g</kbd> prompts for a pattern and then searches the repository using `git grep`
- <kbd>g s</kbd> runs `git status`.  Since `s` is listed under `loopingKeys`, you can press <kbd>s</kbd> repeatedly to update the status, without leaving the current menu.

[magit]: https://magit.vc/
