<p align="center">
<img src="https://github.com/szorfein/ryvst/blob/main/.github/assets/ryvst.png?raw=true" height="300" weight="500" />
</p>

<h1 align="center"></h1>
<div align="center" style="margin-top: 2rem;">
<p>
<a href="https://github.com/szorfein/ryvst/stargazers"><img src="https://img.shields.io/github/stars/szorfein/ryvst?style=for-the-badge&logo=starship&color=C9CBFF&logoColor=C9CBFF&labelColor=302D41" alt="stars"><a>&nbsp;&nbsp;
<a href="https://github.com/szorfein/ryvst/"><img src="https://img.shields.io/github/repo-size/szorfein/ryvst?style=for-the-badge&logo=hyprland&logoColor=f9e2af&label=Size&labelColor=302D41&color=f9e2af" alt="REPO SIZE"></a>&nbsp;&nbsp;
<a href="https://github.com/szorfein/ryvst/commits/main/"><img src="https://img.shields.io/github/last-commit/szorfein/ryvst?style=for-the-badge&logo=github&logoColor=eba0ac&label=Last%20Commit&labelColor=302D41&color=eba0ac" alt="Last Commit"></a>&nbsp;&nbsp;
<a href="https://github.com/szorfein/ryvst/blob/main/LICENSE"><img src="https://img.shields.io/github/license/szorfein/ryvst?style=for-the-badge&logo=&color=CBA6F7&logoColor=CBA6F7&labelColor=302D41" alt="LICENSE"></a>&nbsp;&nbsp;
</p>

A download manager for your dotfiles. Declare your needs in YAML.

This is a rewrite of [Reaver](https://github.com/szorfein/reaver) from Ruby to Go.

</div>

## Development

- Alpha, (current stage)
- Beta
- Release

## Features

- Downloads files very fast via Goroutines
- Uncompresses archive ZIP, TAR support coming soon
- Time units now support Seconds, Minutes and Hours
- Supports Git (only `git clone` and `git pull`); you just need to install Git.
- Collections can execute command, e.g: `fc-cache -f` after downloadings fonts.

## Use from source

```sh
git clone https://github.com/szorfein/ryvst
cd ryvst
go build -o ryvst cmd/ryvst/main.go
```

## Examples

#### Need ohmyzsh with plugins?

[ohmyzsh](https://github.com/ohmyzsh) can be installed in `$HOME/.ohmyzsh` with plugins in `$HOME/.ohmyzsh/custom/plugins`, so create a new collections `zsh.yml`:

    $EDITOR ~/.config/ryvst/zsh.yml

```yml
every_time: 24h
tasks:
  - name: ohmyzsh.zip
    url: https://github.com/ohmyzsh/ohmyzsh/archive/refs/heads/master.zip
    to_dest: .ohmyzsh
  - name: zsh-autosuggestions.zip
    url: https://github.com/zsh-users/zsh-autosuggestions/archive/refs/heads/master.zip
    to_dest: .ohmyzsh/custom/plugins/zsh-autosuggestions
  - name: zsh-syntax-highlighting.zip
    url: https://github.com/zsh-users/zsh-syntax-highlighting/archive/refs/heads/master.zip
    to_dest: .ohmyzsh/custom/plugins/zsh-syntax-highlighting
```

`to_dest` is relative to `$HOME`, `name` is the name of the file after the download. All the files are extracted/moved in the given order, so place `ohmyzsh` first.

#### Something with Git?

If you own a Git repository, you may prefer to use Git instead of an archive.

    $EDITOR ~/.config/ryvst/neovim.yml

```yml
every_time: 12h
tasks:
  - url: https://github.com/szorfein/vie.git
    to_dest: .config/nvim
```

You don't need to use `name` and the `url` should end with `.git`

#### Look in /examples

The [examples](https://github.com/szorfein/ryvst/tree/main/examples) directory should be enough to get you started writing your own file collection. fonts, neovim, ohmyzsh, tmux...
