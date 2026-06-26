# 0.9.1

- shallow-clone blueprint repos by default
	- falls back to full clone if the remote doesn't support shallow cloning

# 0.9.0

- add a `wares build` command, which will build the `_self` blueprint defined in `waresfile.yaml` and symlink artifacts to `<project dir>/wares-result`

# 0.8.10

- (fix) resolve latest commit if there is no locked commit

# 0.8.9

- (fix) respect locked commit when building blueprints
	- probably the longest-standing bug in *wares* because i didn't notice it until now :P
	- what happened before was: even if the locked commit for a blueprint was older than the latest one, the `wares sync` command would overwrite the locked commit with the latest one

# 0.8.8

- update to new warehouse url (https://codeberg.org/indium114/warehouse)

# 0.8.7

- move store dir for shells to `~/.local/share/wares/_shells/<path/to/project>`
	- note: you will have to run `wares shell --clean` the first time you enter a shell after this update
	- for example, the store path for `~/Projects/indium114/wares` will be `~/.local/share/wares/_shells/home/indium114/Projects/wares`
	- this is to provide separation between stores when you use multiple wares shells

# 0.8.6

- don't reinstall wares that haven't been updated
- (fix) move `fmt.Println` call over to `slag`
- (fix) print `no update command` warning correctly

# 0.8.5

> [!note]
> this release is entirely under-the-hood changes

- move all `fmt.Printf` calls over to new `slag` library
- disable `fang`'s default error handler

# 0.8.4

- change log message statuses to all-lowercase
- (fix) make `wares --version` actually show the version

# 0.8.3

- support extracting `.bz2` archives

# 0.8.2

- add `--clean` flag to `sync` and `shell` commands to rebuild blueprints from source

# 0.8.1

> [!WARNING]
> breaking changes in this version!

- set `$WARES_SHELL_ACTIVE` to true when in a wares shell
- (fix/breaking) rename `blueprint` section (incorrect) of `waresfile.yaml` to `blueprints` (correct)

# 0.8.0

- add wares *shells* and accompanying `wares shell` command, allowing temporary environments with certain wares and blueprints installed, a la *Nix shells*
	- configured in `waresfile.yaml`

# 0.7.1

- `wares query` command to get information about installed wares and blueprints

# 0.7.0

> [!WARNING]
> breaking changes in this version!

- support downloading release artifacts from *Forgejo* instances like *Codeberg*
- add `system:` property to wares and blueprints to install to `/Wares`
- (fix) rebuild blueprint after update
- (fix/breaking) trim `https:` from blueprint source path
	- paths will go from `~/.local/share/wares/https:/example.com/foo/bar` to `~/.local/share/wares/example.com/foo/bar`
	- may require an extra sync

# 0.6.1

- show current and new version of packages when running `update` command

# 0.6.0

- `blueprints` system to install and manage programs built from *source*

# 0.5.0

- add the `warehouse` to easily add ware configurations (package definitions) with the `wares add` command

# 0.4.1

- (fix) mark files in `~/Wares` as executable even if sync fails
- (fix) only update symlinks *after* sync completes

# 0.4.0

- `managers` feature to install and manage programs with system package managers (`apt`, `flatpak`, etc.)

# 0.3.4

- show a warning if the user's `~/.config/wares` directory's git tree has uncommitted changes

# 0.3.3

- (fix) allow multiple *GitHub* release artifacts if `multiple` is `true`

# 0.3.2

- (fix) extract properly when `removetoplevel` is `false`

# 0.3.1

- support extracting `.zip` archives
- check for `tar` and `unzip` commands when running `doctor` command

# 0.3.0

> [!WARNING]
> breaking changes in this version!

- (breaking) rename config file from `wares.yaml` to `config.yaml`
- `wares clean` command to remove old versions of packages from `~/.local/share/wares`
- add `removetoplevel` for extracting archives with top-level directories

# 0.2.2

- (fix) only remove symlinks when *syncing*, not when running *update*

# 0.2.1

- add `multiple` option to allow symlinking multiple artifacts
- (fix/optimisation) don't check for packages that need to be uninstalled on every loop

# 0.2.0

- `.tar.gz` archive extraction
- (fix) don't exit after creating directories with `wares doctor` command

# 0.1.4

- delete packages when removed from config

# 0.1.3

- (fix) *actually* mark files in `~/Wares` as executable

# 0.1.2

- (fix) mark files in `~/Wares` as executable
- display download output from `gh` command

# 0.1.1

- change `[ERR]` to `[ERROR]` for error messages

# 0.1.0

- the first release!
