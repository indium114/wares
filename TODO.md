# TODO

## Needed for 0.3.1 release

- [x] Add `tar` and `unzip` command availability check to `wares doctor`
- [x] `.zip` extraction

## Needed for 0.3.3

- [x] Fix: allow users to download multiple artifacts if `multiple: true`

## Needed for 0.3.4

- [ ] Add dirty git repo warning when changes are unstaged
	- Use `git status --porcelain`, ignore error if it matches "fatal: not a git repository (or any of the parent directories): .git"

## Needed for 0.4.0

- [ ] Allow users to manage their distro's package manager in wares config
  - [ ] Configure how their package manager handles installation, removal, and upgrading packages
  - [ ] Allow configuration of multiple package managers

## Needed for 0.5.0

- [ ] `wares add` command to add pre-made ware configurations for packages from a centralised repository
	- waresHub

## Not needed currently

- [ ] Explore configuring in `pkl` rather than `yaml`
- [x] Logo
