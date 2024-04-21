# Copr Manager

CLI app for managing Copr repos, written in Go.

## Usage
`cpm [OPTION] [REPO(s)...]`

```
Options:
  enable    Add or enable one or more Copr repositories
  remove    Remove one or more Copr repositories
    --all         Remove all installed Copr repositories
  list      List Copr repositories in your repo folder
    --enabled     List all enabled repositories (default)
    --disabled    List all disabled repositories
    --all         List both disabled and enabled repositories
  disable   Disable one or more Copr repositories without deleting the repository files
  help      Display help text

Arguments:
  [REPO(s)...]  One or more repository names formatted as `author/repo`

Examples:
  cpm enable kylegospo/bazzite yalter/niri
  cpm disable kylegospo/bazzite 
  cpm remove kylegospo/bazzite
  cpm list --all
```