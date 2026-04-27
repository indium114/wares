# TODO

## Needed for 0.3.0 release

- [x] BREAKING: Change config filename from `wares.yaml` to `config.yaml`
- [x] Use `name` field for symlinking rather than top-level entry name
  - [x] Allow symlinking files multiple levels deep inside archive (optional --strip-components=1)
- [ ] `clean` command to remove old versions of packages

## Not needed currently

- [ ] Explore configuring in `pkl` rather than `yaml`
- [ ] Allow users to manage their distro's package manager in wares config
  - [ ] Configure how their package manager handles installation, removal, and upgrading packages
  - [ ] Allow configuration of multiple package managers
