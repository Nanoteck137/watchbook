{
  description = "Devshell for watchbook";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    gitignore.url = "github:hercules-ci/gitignore.nix";
    gitignore.inputs.nixpkgs.follows = "nixpkgs";

    devtools.url     = "github:nanoteck137/devtools";
    devtools.inputs.nixpkgs.follows = "nixpkgs";

    tagopus.url      = "github:nanoteck137/tagopus/v0.1.1";
    tagopus.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gitignore, devtools, tagopus, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        backend = pkgs.buildGoModule {
          pname = "watchbook";
          version = fullVersion;
          src = ./.;
          subPackages = ["cmd/watchbook"];

          ldflags = [
            "-X github.com/nanoteck137/watchbook.Version=${version}"
            "-X github.com/nanoteck137/watchbook.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          tags = [];

          vendorHash = "";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
            wrapProgram $out/bin/watchbook --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
          '';
        };

        frontend = pkgs.buildNpmPackage {
          name = "watchbook-web";
          version = fullVersion;

          src = gitignore.lib.gitignoreSource ./web;
          npmDepsHash = "";

          PUBLIC_VERSION=version;
          PUBLIC_COMMIT=self.dirtyRev or self.rev or "no-commit";

          installPhase = ''
            runHook preInstall
            cp -r build $out/
            echo '{ "type": "module" }' > $out/package.json

            mkdir $out/bin
            echo -e "#!${pkgs.runtimeShell}\n${pkgs.nodejs}/bin/node $out\n" > $out/bin/watchbook-web
            chmod +x $out/bin/watchbook-web

            runHook postInstall
          '';
        };

        tools = devtools.packages.${system};
      in
      {
        packages = {
          default = backend;
          inherit backend frontend;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            air
            go
            gopls
            nodejs
            imagemagick
            ffmpeg

            tagopus.packages.${system}.default
            tools.publishVersion
          ];
        };
      }
    ) // {
      nixosModules.backend = import ./nix/backend.nix { inherit self; };
      nixosModules.frontend = import ./nix/frontend.nix { inherit self; };
      nixosModules.default = { ... }: {
        imports = [
          self.nixosModules.backend
          self.nixosModules.frontend
        ];
      };
    };
}
