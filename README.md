<p align="center">
<img src="https://github.com/szorfein/ryvst/blob/main/.github/assets/ryvst.png?raw=true" height="300" weight="510" />
<br />
<small align="center">Logo generated on gradually.ai</small>
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

We all like writing custom shell scripts, just to get wallpapers, fonts, cloning things we like. All these scripts do the same thing over and over, with the only variation being the "source" and the "destination".

Ryvst simplifies the installation of your "environment" by grabbing all the stuff you need so you can avoid (re)writing simple scripts using massivly `curl`, `wget`, `unzip`, `untar` or `git`. Instead you write a few lines in a YAML file. Ryvst does the job and checks for updates at the specified interval.

## Development

- Alpha
- Beta, (current stage)
- Release

## Features

- Downloads files very fast via Goroutines
- Can decompresses archive zip, tar.gz, tar.xz
- Time units now support Seconds, Minutes and Hours
- Supports Git (only `git clone` and `git pull`); you just need to install Git.
- Collections can execute command, e.g: `fc-cache -f` after installing fonts.

## Use from source

```sh
git clone https://github.com/szorfein/ryvst
cd ryvst
go build -o ryvst cmd/ryvst/main.go
```

## Usage

Fill your `$HOME/.config/ryvst` directory with YAML files defining your needs, then launch `ryvst`.

You can automate the startup of `ryvst` by including it in your `bash_profile`, `zprofile` or your `sway` configuration...

## Starting on Sway

After copying the binary in `$HOME/bin`, you can configure Sway like this:

```
# Run once at login
exec {
  ~/bin/ryvst
}
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
