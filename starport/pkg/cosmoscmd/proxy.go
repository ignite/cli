package cosmoscmd

import (
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

func startProxyForTunneledPeers(clientCtx client.Context, cmd *cobra.Command) error {
	if cmd.Name() == "start" {
		serverCtx := server.GetServerContextFromCmd(cmd)
		cmdCtx := cmd.Context()
		if gitpod.IsOnGitpod() {
			serverCtx.Logger.Info("Starting chisel server", "port", xchisel.DefaultServerPort, "proxy", serverCtx.Config.P2P.ListenAddress)
			go func() {
				err := xchisel.StartServer(cmdCtx, xchisel.DefaultServerPort, serverCtx.Config.P2P.ListenAddress)
				if err != nil {
					serverCtx.Logger.Error("Failed to start chisel server", "port", xchisel.DefaultServerPort)
				}
			}()
		}

		tunneledPeersConfig, err := networkchain.GetSPNConfig(filepath.Join(clientCtx.HomeDir, cosmosutil.ChainConfigDir, networkchain.SPNConfigFile))
		if err == nil {
			for _, peer := range tunneledPeersConfig.TunneledPeers {
				if peer.Name == networkchain.HTTPTunnelChisel {
					peer := peer
					serverCtx.Logger.Info("Starting chisel client", "tunnelAddress", peer.Address, "localPort", peer.LocalPort)
					go func() {
						err := xchisel.StartClient(cmdCtx, peer.Address, peer.LocalPort, xchisel.DefaultServerPort)
						if err != nil {
							serverCtx.Logger.Error("Failed to start chisel client", "tunnelAddress", peer.Address, "localPort", peer.LocalPort)
							return
						}
					}()
				}
			}
		}
	}

	return nil
}
