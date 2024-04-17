# copr-tool

CLI app for managing Copr repos, written in Go.

```shell
Usage: copr-tool [OPTION] [REPO...]

Options:
  enable    Add or enable one or more Copr repositories.
  remove    Remove one or more Copr repositories.
  list      List all (enabled and disabled) Copr repositories in your repo folder.
  disable   Disable one or more Copr repositories without deleting the repository files.
  help      Display help text.

Arguments:
  [REPO...]  One or more repository names formatted as `author/repo`

Examples:
  copr-tool enable kylegospo/bazzite
  copr-tool disable kylegospo/bazzite 
  copr-tool remove kylegospo/bazzite
  copr-tool list
```