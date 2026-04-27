{
  description = "go devshell and package, created by scaffolder";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          name = "go-devshell";

          packages = with pkgs; [
            go
            gopls
            gotools
            delve
            just
            goreleaser
          ];
        };

        packages.wares = pkgs.buildGoModule {
          pname = "wares";
          version = "0.2.1";

          src = self;

          vendorHash = "sha256-tRsa4osUQUKQ+QrYJO5kTTT43w8gfnbFbxpd3edkRSE=";

          subPackages = [ "." ];
          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "A declarative AppImage/binary package manager";
            license = licenses.mit;
            platforms = platforms.all;
          };
        };

        apps.wares = {
          type = "app";
          program = "${self.packages.${pkgs.stdenv.hostPlatform.system}.wares}/bin/wares";
        };
      });
}
