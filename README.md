# wares

*wares* is a declarative AppImage/binary package manager!

## Installation

### Downloading wares

To install, just grab the binary for your operating system from the **Releases** section on the right.

### Setting up wares

Run the following to check that everything is in order:

```shell
/path/to/wares doctor
```

If it tells you that ~/Wares is not in your `$PATH`, please add it.

### Letting wares manage itself

Then, create `~/.config/wares` and paste the following into `~/.config/wares/wares.yaml`:

```yaml
wares:
  wares:
    name: wares
    repo: indium114/wares
    asset: "wares_Linux_x86_64"
```

Replace `Linux_x86_64` with `Darwin_aarch64` if you're on a Mac with Apple Silicon, or `Darwin_x86_64` if you're on an Intel Mac.

Then, run `/path/to/wares sync` to download Wares, and it will now manage itself.

## Usage

### Installing a package

To install a package, add it to the `wares` section of `wares.yaml`.

For example, here's me installing [Helix](https://github.com/helix-editor/helix) using Wares

```yaml
wares:
  hx: # This will be the name of the resulting binary
    name: hx                 # Doesn't currently do anything
    repo: helix-editor/helix # GitHub repo (without github.com)
    asset: "*.AppImage"      # Pattern which will match the downloaded asset you would like
                               # For example, using "*Linux-x86_64*" will match with any file containing the substring `Linux-x86_64` in its name
```

### Updating packages

To update packages, run the following command:

```shell
wares update
```

This will update the version in `pallet.lock`. Now just sync to install the new version:

```shell
wares sync
```
