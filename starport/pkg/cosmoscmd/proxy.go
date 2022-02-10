package cosmoscmd

import (
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/ctxticker"
	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

const TunnelRerunDelay = 5 * time.Second

func startProxyForTunneledPeers(clientCtx client.Context, cmd *cobra.Command) {
	if cmd.Name() != "start" {
		return
	}
	serverCtx := server.GetServerContextFromCmd(cmd)
	ctx := cmd.Context()
	if gitpod.IsOnGitpod() {
		go func() {
			_ = ctxticker.DoNow(ctx, TunnelRerunDelay, func() error {
				serverCtx.Logger.Info("Starting chisel server", "port", xchisel.DefaultServerPort)
				err := xchisel.StartServer(ctx, xchisel.DefaultServerPort)
				if err != nil {
					serverCtx.Logger.Error(
						"Failed to start chisel server",
						"port", xchisel.DefaultServerPort,
						"reason", err.Error(),
					)
				}
				return nil
			})
		}()
	}

	spnConfig, err := networkchain.GetSPNConfig(filepath.Join(clientCtx.HomeDir, cosmosutil.ChainConfigDir, networkchain.SPNConfigFile))
	if err != nil {
		serverCtx.Logger.Error("Failed to open spn config file", "reason", err.Error())
	}
	for _, peer := range spnConfig.TunneledPeers {
		if peer.Name == networkchain.HTTPTunnelChisel {
			peer := peer
			go func() {
				_ = ctxticker.DoNow(ctx, TunnelRerunDelay, func() error {
					serverCtx.Logger.Info("Starting chisel client", "tunnelAddress", peer.Address, "localPort", peer.LocalPort)
					err := xchisel.StartClient(ctx, peer.Address, peer.LocalPort, "26656")
					if err != nil {
						serverCtx.Logger.Error("Failed to start chisel client",
							"tunnelAddress", peer.Address,
							"localPort", peer.LocalPort,
							"reason", err.Error(),
						)
					}
					return nil
				})
			}()
		}
	}
}
