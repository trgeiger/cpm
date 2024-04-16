# copr-tool

CLI app for managing COPR repos, written in Go.

```shell
Usage: copr-tool [OPTION] [REPO...]

Options:
  enable    Add or enable one or more COPR repositories.
  remove    Remove one or more COPR repositories.
  list      List all (enabled and disabled) COPR repositories in your repo folder.
  disable   Disable one or more COPR repositories without deleting the repository files.
  help      Display help text.

Arguments:
  [REPO...]  One or more repository names formatted as `author/repo`

Examples:
  copr-tool enable kylegospo/bazzite
  copr-tool disable kylegospo/bazzite 
  copr-tool remove kylegospo/bazzite
  copr-tool list
```