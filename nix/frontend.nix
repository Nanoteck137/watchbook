{ self }:
{ config, lib, pkgs, ... }:
with lib; let
  cfg = config.services.watchbook-web;
in
{
  options.services.watchbook-web = {
    enable = mkEnableOption "Enable the watchbook-web service";

    port = mkOption {
      type = types.port;
      default = 5425;
      description = "port to listen on";
    };

    host = mkOption {
      type = types.str;
      default = "";
      description = "hostname or address to listen on";
    };

    apiAddress = mkOption {
      type = types.str;
      description = "address to the api server";
    };

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.system}.frontend;
      description = "package to use for this service (defaults to the one in the flake)";
    };

    user = mkOption {
      type = types.str;
      default = "watchbook-web";
      description = lib.mdDoc "user to use for this service";
    };

    group = mkOption {
      type = types.str;
      default = "watchbook-web";
      description = lib.mdDoc "group to use for this service";
    };

    openFirewall = mkOption {
      type = types.bool;
      default = false;
      description = "open the ports in the firewall";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.watchbook-web = {
      description = "Frontend for watchbook";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];

      environment = {
        PORT = "${toString cfg.port}";
        HOST = "${cfg.host}";
        API_ADDRESS = "${cfg.apiAddress}";
        HOST_HEADER = "x-forwarded-host";
        BODY_SIZE_LIMIT = "Infinity";
      };

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;

        ExecStart = "${cfg.package}/bin/watchbook-web";

        Restart = "on-failure";
        RestartSec = "5s";

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
      };
    };

    networking.firewall = lib.mkIf cfg.openFirewall {
      allowedTCPPorts = [ cfg.port ];
    };

    users.users = mkIf (cfg.user == "watchbook-web") {
      watchbook-web = {
        group = cfg.group;
        isSystemUser = true;
      };
    };

    users.groups = mkIf (cfg.group == "watchbook-web") {
      watchbook-web = {};
    };
  };
}
