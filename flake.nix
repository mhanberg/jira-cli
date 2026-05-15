{
  description = "jira-cli development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ { flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      perSystem = { pkgs, lib, ... }: {
        packages.default = pkgs.buildGoModule rec {
          pname = "jira-cli";
          version = "0.0.0-dev";

          src = lib.cleanSourceWith {
            src = ./.;
            filter = path: type:
              let baseName = baseNameOf (toString path);
              in !(lib.elem baseName [ ".direnv" ".go" ".gocache" "bin" "vendor" "dist" "coverage" ]);
          };

          vendorHash = "sha256-cl+Sfi9WSPy8qOtB13rRiKtQdDC+HC0+FMKpsWbtU2w=";

          subPackages = [ "cmd/jira" ];

          env.CGO_ENABLED = "0";

          ldflags = [
            "-s"
            "-w"
            "-X github.com/ankitpokhrel/jira-cli/internal/version.Version=${version}"
          ];

          doCheck = false;

          meta = {
            description = "Feature-rich interactive Jira command line";
            homepage = "https://github.com/ankitpokhrel/jira-cli";
            license = lib.licenses.mit;
            mainProgram = "jira";
            platforms = lib.platforms.unix;
          };
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_25
            gopls
            gotools
            go-tools
            golangci-lint
            delve
            gomodifytags
            gotests
            impl
            govulncheck
          ];

          shellHook = ''
            export GOBIN="$PWD/.go/bin"
            export GOCACHE="$PWD/.gocache"
            mkdir -p "$GOBIN"
            export PATH="$GOBIN:$PATH"
            echo "go $(go version | awk '{print $3}') dev shell"
          '';
        };
      };
    };
}
