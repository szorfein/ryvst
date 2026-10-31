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

A download manager for your dotfiles, declare your need in YAML.

This is a rewrite of [Reaver](https://github.com/szorfein/reaver) from Ruby to Go.

</div>

## Development

- Alpha, (current stage)
- Beta
- First Release

## Features

...

## Use from source

```sh
git clone https://github.com/szorfein/ryvst
cd ryvst
go build -o ryvst cmd/ryvst/main.go
```

## Examples

#### Need ohmyzsh?

[ohmyzsh](https://github.com/ohmyzsh) can be installed in `$HOME/.ohmyzsh`, plugins in `$HOME/.ohmyzsh/custom`, so create a new collections `zsh.yml`: $EDITOR ~/.config/ryvst/zsh.yml

```yml
every_time: 1d
collections:
  - name: ohmyzsh.zip
    url: https://github.com/ohmyzsh/ohmyzsh/archive/refs/heads/master.zip
    to_dest: .ohmyzsh
  - name: zsh-autosuggestions.tar.gz
    url: https://github.com/zsh-users/zsh-autosuggestions/archive/refs/heads/master.tar.gz
    to_dest: .oh-my-zsh/custom/plugins/zsh-autosuggestions
  - name: zsh-syntax-highlighting.tar.gz
    url: https://github.com/zsh-users/zsh-syntax-highlighting/archive/refs/heads/master.tar.gz
    to_dest: .oh-my-zsh/custom/plugins/zsh-syntax-highlighting
```

`to_dest` is relative to `$HOME`, `name` is the name of the file after the download. All the files are extracted in the given order, so place `ohmyzsh` first.
