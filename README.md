# go-git-tk

## What is this?
`Go Git Toolkit` is a frontend for bare Git repositories. Its purpose is to serve as a TUI Manager for a Gitserver. See [SETUP.md](SETUP.md) for assistance

You have the options to:
- Create new Repositories
- Shortcuts to edit their git hooks and the `~/.ssh/authorized_keys` file

## Caveats
- You can grant users only the whole database or nothing

## Building
Ensure you have `make` and `golang >= 1.22.3` installed

```bash
make all
```