{ self }: 
{ config, lib, pkgs, ... }:
with lib; let
  cfg = config.services.watchbook;

  watchbookConfig = pkgs.writeText "config.toml" ''
    listen_addr = "${cfg.host}:${toString cfg.port}"
    data_dir = "${cfg.dataDir}"
    username = "${cfg.username}"
    initial_password = "${cfg.initialPassword}"
    jwt_secret = "${cfg.jwtSecret}"
  '';
in
{
  options.services.watchbook = {
    enable = mkEnableOption "Enable the watchbook service";

    port = mkOption {
      type = types.port;
      default = 5424;
      description = "port to listen on";
    };

    host = mkOption {
      type = types.str;
      default = "";
      description = "hostname or address to listen on";
    };

    dataDir = mkOption {
      type = types.path;
      default = "/var/lib/watchbook";
      description = "path to the data directory";
    };

    username = mkOption {
      type = types.str;
      description = "username of the first user";
    };

    initialPassword = mkOption {
      type = types.str;
      description = "initial password of the first user (should change after the first login)";
    };

    jwtSecret = mkOption {
      type = types.str;
      description = "jwt secret";
    };

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.system}.backend;
      description = "package to use for this service (defaults to the one in the flake)";
    };

    user = mkOption {
      type = types.str;
      default = "watchbook";
      description = "user to use for this service";
    };

    group = mkOption {
      type = types.str;
      default = "watchbook";
      description = "group to use for this service";
    };

    openFirewall = mkOption {
      type = types.bool;
      default = false;
      description = "open the ports in the firewall";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.watchbook = {
      description = "watchbook";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];

      serviceConfig = mkMerge [
        {
          User = cfg.user;
          Group = cfg.group;

          ExecStart = "${cfg.package}/bin/watchbook serve -c '${watchbookConfig}'";

          Restart = "on-failure";
          RestartSec = "5s";

          PrivateTmp = true;
          ProtectHome = true;
          ProtectHostname = true;
          ProtectKernelLogs = true;
          ProtectKernelModules = true;
          ProtectKernelTunables = true;
          ProtectProc = "invisible";
          ProtectSystem = "strict";
          RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
          RestrictNamespaces = true;
          RestrictRealtime = true;
          RestrictSUIDSGID = true;
        }

        (mkIf (cfg.dataDir != "/var/lib/watchbook") {
          ReadWritePaths = [ cfg.dataDir ];
        })

        (mkIf (cfg.dataDir == "/var/lib/watchbook") {
          StateDirectory = "watchbook";
        })
      ];
    };

    networking.firewall = lib.mkIf cfg.openFirewall {
      allowedTCPPorts = [ cfg.port ];
    };

    users.users = mkIf (cfg.user == "watchbook") {
      watchbook = {
        group = cfg.group;
        isSystemUser = true;
        home = "${cfg.dataDir}";
      };
    };

    users.groups = mkIf (cfg.group == "watchbook") {
      watchbook = {};
    };
  };
}
