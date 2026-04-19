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
          ];
        };

        packages.wares = pkgs.buildGoModule {
          pname = "wares";
          version = "0.1.0";

          src = self;

          vendorHash = pkgs.lib.fakeHash; # TODO: Change

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
