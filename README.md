# Copr Manager
 
A command line tool for managing Copr repositories, written in Go.
.

## Usage
`cpm [OPTION] [REPO(s)...]`

```
Usage:
  cpm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  disable     Disable one or more Copr repositories without uninstalling them.
  enable      Enable or add one or more Copr repositories.
  help        Help about any command
  list        List installed Copr repositories
    --enabled     List all enabled repositories (default)
    --disabled    List all disabled repositories
    --all         List both disabled and enabled repositories
  prune       Remove duplicate repository configurations.
  remove      Uninstall one or more Copr repositories.
    --all         Remove all installed Copr repositories

Arguments:
  [REPO(s)...]  One or more repository names formatted as `author/repo`

Examples:
  cpm enable kylegospo/bazzite yalter/niri
  cpm disable kylegospo/bazzite 
  cpm remove kylegospo/bazzite
  cpm list --all
```

## Building
```shell
go build -o cpm main.go
```
