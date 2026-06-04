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

        packages.wares = pkgs.buildGoModule rec {
          pname = "wares";
          version = "0.8.7";

          src = self;

          vendorHash = "sha256-m5KteI/VNM2/EjzyLmNw/5Amvbwe8g2Bk0JB1D5JwpY=";

          subPackages = [ "." ];
          ldflags = [ "-s" "-w" "-X 'github.com/indium114/wares/cmd.Version=${version}'" ];

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
