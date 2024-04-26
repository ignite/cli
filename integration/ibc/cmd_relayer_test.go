//go:build !relayer

package ibc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/v28/ignite/config/chain"
	"github.com/ignite/cli/v28/ignite/config/chain/base"
	v1 "github.com/ignite/cli/v28/ignite/config/chain/v1"
	"github.com/ignite/cli/v28/ignite/pkg/availableport"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/goanalysis"
	yamlmap "github.com/ignite/cli/v28/ignite/pkg/yaml"
	envtest "github.com/ignite/cli/v28/integration"
)

const (
	relayerMnemonic = "great immense still pill defense fetch pencil slow purchase symptom speed arm shoot fence have divorce cigar rapid hen vehicle pear evolve correct nerve"
)

var (
	bobName    = "bob"
	marsConfig = v1.Config{
		Config: base.Config{
			Version: 1,
			Build: base.Build{
				Proto: base.Proto{
					Path:            "proto",
					ThirdPartyPaths: []string{"third_party/proto", "proto_vendor"},
				},
			},
			Accounts: []base.Account{
				{
					Name:     "alice",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "slide moment original seven milk crawl help text kick fluid boring awkward doll wonder sure fragile plate grid hard next casual expire okay body",
				},
				{
					Name:     "bob",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "trap possible liquid elite embody host segment fantasy swim cable digital eager tiny broom burden diary earn hen grow engine pigeon fringe claim program",
				},
				{
					Name:     "relayer",
					Coins:    []string{"100000000000token", "1000000000000000000000stake"},
					Mnemonic: relayerMnemonic,
				},
			},
			Faucet: base.Faucet{
				Name:  &bobName,
				Coins: []string{"500token", "100000000stake"},
				Host:  ":4501",
			},
			Genesis: yamlmap.Map{"chain_id": "mars-1"},
		},
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
				Client: yamlmap.Map{"keyring-backend": keyring.BackendTest},
				App: yamlmap.Map{
					"api":      yamlmap.Map{"address": ":1318"},
					"grpc":     yamlmap.Map{"address": ":9092"},
					"grpc-web": yamlmap.Map{"address": ":9093"},
				},
				Config: yamlmap.Map{
					"p2p": yamlmap.Map{"laddr": ":26658"},
					"rpc": yamlmap.Map{"laddr": ":26658", "pprof_laddr": ":6061"},
				},
				Home: "$HOME/.mars",
			},
		},
	}
	earthConfig = v1.Config{
		Config: base.Config{
			Version: 1,
			Build: base.Build{
				Proto: base.Proto{
					Path:            "proto",
					ThirdPartyPaths: []string{"third_party/proto", "proto_vendor"},
				},
			},
			Accounts: []base.Account{
				{
					Name:     "alice",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "slide moment original seven milk crawl help text kick fluid boring awkward doll wonder sure fragile plate grid hard next casual expire okay body",
				},
				{
					Name:     "bob",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "trap possible liquid elite embody host segment fantasy swim cable digital eager tiny broom burden diary earn hen grow engine pigeon fringe claim program",
				},
				{
					Name:     "relayer",
					Coins:    []string{"100000000000token", "1000000000000000000000stake"},
					Mnemonic: relayerMnemonic,
				},
			},
			Faucet: base.Faucet{
				Name:  &bobName,
				Coins: []string{"500token", "100000000stake"},
				Host:  ":4500",
			},
			Genesis: yamlmap.Map{"chain_id": "earth-1"},
		},
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
				Client: yamlmap.Map{"keyring-backend": keyring.BackendTest},
				App: yamlmap.Map{
					"api":      yamlmap.Map{"address": ":1317"},
					"grpc":     yamlmap.Map{"address": ":9090"},
					"grpc-web": yamlmap.Map{"address": ":9091"},
				},
				Config: yamlmap.Map{
					"p2p": yamlmap.Map{"laddr": ":26656"},
					"rpc": yamlmap.Map{"laddr": ":26656", "pprof_laddr": ":6060"},
				},
				Home: "$HOME/.earth",
			},
		},
	}

	nameSendIbcPost = "SendIbcPost"
	funcSendIbcPost = `package keeper
func (k msgServer) SendIbcPost(goCtx context.Context, msg *types.MsgSendIbcPost) (*types.MsgSendIbcPostResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    // Construct the packet
    var packet types.IbcPostPacketData
    packet.Title = msg.Title
    packet.Content = msg.Content
    // Transmit the packet
    _, err := k.TransmitIbcPostPacket(
        ctx,
        packet,
        msg.Port,
        msg.ChannelID,
        clienttypes.ZeroHeight(),
        msg.TimeoutTimestamp,
    )
    return &types.MsgSendIbcPostResponse{}, err
}`

	nameOnRecvIbcPostPacket = "OnRecvIbcPostPacket"
	funcOnRecvIbcPostPacket = `package keeper
func (k Keeper) OnRecvIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) (packetAck types.IbcPostPacketAck, err error) {
    // validate packet data upon receiving
    if err := data.ValidateBasic(); err != nil {
        return packetAck, err
    }
    packetAck.PostId = k.AppendPost(ctx, types.Post{Title: data.Title, Content: data.Content})
    return packetAck, nil
}`

	nameOnAcknowledgementIbcPostPacket = "OnAcknowledgementIbcPostPacket"
	funcOnAcknowledgementIbcPostPacket = `package keeper
func (k Keeper) OnAcknowledgementIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData, ack channeltypes.Acknowledgement) error {
    switch dispatchedAck := ack.Response.(type) {
    case *channeltypes.Acknowledgement_Error:
        // We will not treat acknowledgment error in this tutorial
        return nil
    case *channeltypes.Acknowledgement_Result:
        // Decode the packet acknowledgment
        var packetAck types.IbcPostPacketAck
        if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
            // The counter-party module doesn't implement the correct acknowledgment format
            return errors.New("cannot unmarshal acknowledgment")
        }

        k.AppendSentPost(ctx,
            types.SentPost{
                PostId:  packetAck.PostId,
                Title:   data.Title,
                Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
            },
        )
        return nil
    default:
        return errors.New("the counter-party module does not implement the correct acknowledgment format")
    }
}`

	nameOnTimeoutIbcPostPacket = "OnTimeoutIbcPostPacket"
	funcOnTimeoutIbcPostPacket = `package keeper
func (k Keeper) OnTimeoutIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) error {
    k.AppendTimeoutPost(ctx,
        types.TimeoutPost{
            Title:   data.Title,
            Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
        },
    )
    return nil
}`
)

type (
	QueryChannels struct {
		Channels []struct {
			ChannelId      string   `json:"channel_id"`
			ConnectionHops []string `json:"connection_hops"`
			Counterparty   struct {
				ChannelId string `json:"channel_id"`
				PortId    string `json:"port_id"`
			} `json:"counterparty"`
			Ordering string `json:"ordering"`
			PortId   string `json:"port_id"`
			State    string `json:"state"`
			Version  string `json:"version"`
		} `json:"channels"`
	}

	QueryBalances struct {
		Balances sdk.Coins `json:"balances"`
	}
)

func runChain(
	t *testing.T,
	ctx context.Context,
	env envtest.Env,
	app envtest.App,
	cfg v1.Config,
	tmpDir string,
	ports []uint,
) (api, rpc, grpc, faucet string) {
	t.Helper()
	if len(ports) < 7 {
		t.Fatalf("invalid number of ports %d", len(ports))
	}

	var (
		chainID   = cfg.Genesis["chain_id"].(string)
		chainPath = filepath.Join(tmpDir, chainID)
		homePath  = filepath.Join(chainPath, "home")
		cfgPath   = filepath.Join(chainPath, chain.ConfigFilenames[0])
	)
	require.NoError(t, os.MkdirAll(chainPath, os.ModePerm))

	genAddr := func(port uint) string {
		return fmt.Sprintf(":%d", port)
	}

	cfg.Validators[0].Home = homePath

	cfg.Faucet.Host = genAddr(ports[0])
	cfg.Validators[0].App["api"] = yamlmap.Map{"address": genAddr(ports[1])}
	cfg.Validators[0].App["grpc"] = yamlmap.Map{"address": genAddr(ports[2])}
	cfg.Validators[0].App["grpc-web"] = yamlmap.Map{"address": genAddr(ports[3])}
	cfg.Validators[0].Config["p2p"] = yamlmap.Map{"laddr": genAddr(ports[4])}
	cfg.Validators[0].Config["rpc"] = yamlmap.Map{
		"laddr":       genAddr(ports[5]),
		"pprof_laddr": genAddr(ports[6]),
	}

	file, err := os.Create(cfgPath)
	require.NoError(t, err)
	require.NoError(t, yaml.NewEncoder(file).Encode(cfg))
	require.NoError(t, file.Close())

	app.SetConfigPath(cfgPath)
	app.SetHomePath(homePath)
	go func() {
		env.Must(app.Serve("should serve chain", envtest.ExecCtx(ctx)))
	}()

	genHTTPAddr := func(port uint) string {
		return fmt.Sprintf("http://127.0.0.1:%d", port)
	}
	return genHTTPAddr(ports[1]), genHTTPAddr(ports[5]), genHTTPAddr(ports[2]), genHTTPAddr(ports[0])
}

func TestBlogIBC(t *testing.T) {
	var (
		env    = envtest.New(t)
		app    = env.Scaffold("github.com/test/planet")
		ctx    = env.Ctx()
		tmpDir = t.TempDir()
	)
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(func() {
		cancel()
		time.Sleep(5 * time.Second)
		require.NoError(t, os.RemoveAll(tmpDir))
	})

	env.Must(env.Exec("create an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"module",
				"blog",
				"--ibc",
				"--require-registration",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a post type list in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"post",
				"title",
				"content",
				"--no-message",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a sentPost type list in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"sentPost",
				"postID:uint",
				"title",
				"chain",
				"--no-message",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a timeoutPost type list in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"timeoutPost",
				"title",
				"chain",
				"--no-message",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a ibcPost package in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"packet",
				"ibcPost",
				"title",
				"content",
				"--ack",
				"postID:uint",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	blogKeeperPath := filepath.Join(app.SourcePath(), "x/blog/keeper")
	require.NoError(t, goanalysis.ReplaceCode(
		blogKeeperPath,
		nameSendIbcPost,
		funcSendIbcPost,
	))
	require.NoError(t, goanalysis.ReplaceCode(
		blogKeeperPath,
		nameOnRecvIbcPostPacket,
		funcOnRecvIbcPostPacket,
	))
	require.NoError(t, goanalysis.ReplaceCode(
		blogKeeperPath,
		nameOnAcknowledgementIbcPostPacket,
		funcOnAcknowledgementIbcPostPacket,
	))
	require.NoError(t, goanalysis.ReplaceCode(
		blogKeeperPath,
		nameOnTimeoutIbcPostPacket,
		funcOnTimeoutIbcPostPacket,
	))

	// serve both chains.
	ports, err := availableport.Find(
		14,
		availableport.WithMinPort(4000),
		availableport.WithMaxPort(5000),
	)
	require.NoError(t, err)
	earthAPI, earthRPC, earthGRPC, earthFaucet := runChain(t, ctx, env, app, earthConfig, tmpDir, ports[:7])
	earthChainID := earthConfig.Genesis["chain_id"].(string)
	earthHome := earthConfig.Validators[0].Home
	marsAPI, marsRPC, marsGRPC, marsFaucet := runChain(t, ctx, env, app, marsConfig, tmpDir, ports[7:])
	marsChainID := marsConfig.Genesis["chain_id"].(string)
	marsHome := marsConfig.Validators[0].Home

	// check the chains is up
	stepsCheckChains := step.NewSteps(
		step.New(
			step.Exec(
				app.Binary(),
				"config",
				"output", "json",
			),
			step.PreExec(func() error {
				if err := env.IsAppServed(ctx, earthAPI); err != nil {
					return err
				}
				return env.IsAppServed(ctx, marsAPI)
			}),
		),
	)
	env.Exec("waiting the chain is up", stepsCheckChains, envtest.ExecRetry())

	// ibc relayer.
	env.Must(env.Exec("install the hermes relayer app",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"app",
				"install",
				"-g",
				// filepath.Join(goenv.GoPath(), "src/github.com/ignite/apps/hermes"), // Local path for test proposals
				"github.com/ignite/apps/hermes@hermes/v0.2.1",
			),
		)),
	))

	env.Must(env.Exec("configure the hermes relayer app",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"relayer",
				"hermes",
				"configure",
				earthChainID,
				earthRPC,
				earthGRPC,
				marsChainID,
				marsRPC,
				marsGRPC,
				"--chain-a-faucet", earthFaucet,
				"--chain-b-faucet", marsFaucet,
				"--generate-wallets",
				"--overwrite-config",
			),
		)),
	))

	go func() {
		env.Must(env.Exec("run the hermes relayer",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, "relayer", "hermes", "start", earthChainID, marsChainID),
			)),
			envtest.ExecCtx(ctx),
		))
	}()
	time.Sleep(3 * time.Second)

	var (
		queryOutput   = &bytes.Buffer{}
		queryResponse QueryChannels
	)
	env.Must(env.Exec("verify if the channel was created", step.NewSteps(
		step.New(
			step.Stdout(queryOutput),
			step.Exec(
				app.Binary(),
				"q",
				"ibc",
				"channel",
				"channels",
				"--node", earthRPC,
				"--log_format", "json",
				"--output", "json",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(queryOutput.Bytes(), &queryResponse); err != nil {
					return fmt.Errorf("unmarshling tx response: %w", err)
				}
				if len(queryResponse.Channels) == 0 ||
					len(queryResponse.Channels[0].ConnectionHops) == 0 {
					return fmt.Errorf("channel not found")
				}
				if queryResponse.Channels[0].State != "STATE_OPEN" {
					return fmt.Errorf("channel is not open")
				}
				return nil
			}),
		),
	)))

	var (
		sender       = "alice"
		receiverAddr = "cosmos1nrksk5swk6lnmlq670a8kwxmsjnu0ezqts39sa"
		txOutput     = &bytes.Buffer{}
		txResponse   struct {
			Code   int
			RawLog string `json:"raw_log"`
			TxHash string `json:"txhash"`
		}
	)

	stepsTx := step.NewSteps(
		step.New(
			step.Stdout(txOutput),
			step.Exec(
				app.Binary(),
				"tx",
				"ibc-transfer",
				"transfer",
				"transfer",
				"channel-0",
				receiverAddr,
				"100000stake",
				"--from", sender,
				"--node", earthRPC,
				"--home", earthHome,
				"--chain-id", earthChainID,
				"--output", "json",
				"--log_format", "json",
				"--keyring-backend", "test",
				"--yes",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(txOutput.Bytes(), &txResponse); err != nil {
					return fmt.Errorf("unmarshling tx response: %w", err)
				}
				return cmdrunner.New().Run(ctx, step.New(
					step.Exec(
						app.Binary(),
						"q",
						"tx",
						txResponse.TxHash,
						"--node", earthRPC,
						"--home", earthHome,
						"--chain-id", earthChainID,
						"--output", "json",
						"--log_format", "json",
					),
					step.Stdout(txOutput),
					step.PreExec(func() error {
						txOutput.Reset()
						return nil
					}),
					step.PostExec(func(execErr error) error {
						if execErr != nil {
							return execErr
						}
						if err := json.Unmarshal(txOutput.Bytes(), &txResponse); err != nil {
							return err
						}
						return nil
					}),
				))
			}),
		),
	)
	if !env.Exec("send an IBC transfer", stepsTx, envtest.ExecRetry()) {
		t.FailNow()
	}
	require.Equal(t, 0, txResponse.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)

	var (
		balanceOutput   = &bytes.Buffer{}
		balanceResponse QueryBalances
	)
	env.Must(env.Exec("check ibc balance", step.NewSteps(
		step.New(
			step.Stdout(balanceOutput),
			step.Exec(
				app.Binary(),
				"q",
				"bank",
				"balances",
				receiverAddr,
				"--node", marsRPC,
				"--home", marsHome,
				"--log_format", "json",
				"--output", "json",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(balanceOutput.Bytes(), &balanceResponse); err != nil {
					return fmt.Errorf("unmarshalling tx response: %w", err)
				}
				if balanceResponse.Balances.Empty() {
					return fmt.Errorf("empty balances")
				}
				if !strings.HasPrefix(balanceResponse.Balances[0].Denom, "ibc/") {
					return fmt.Errorf("invalid ibc balance: %v", balanceResponse.Balances[0])
				}
				return nil
			}),
		),
	)))

	// TODO test ibc using the blog post methods:
	// step.Exec(app.Binary(), "tx", "blog", "send-ibc-post", "transfer", "channel-0", "Hello", "Hello_Mars-Alice_from_Earth", "--chain-id", earthChainID, "--from", "alice", "--node", earthGRPC, "--output", "json", "--log_format", "json", "--yes")
	// TODO test ibc using the hermes ft-transfer:
	// step.Exec(envtest.IgniteApp, "relayer", "hermes", "exec", "--", "--config", earthConfig, "tx", "ft-transfer", "--timeout-seconds", "1000", "--dst-chain", earthChainID, "--src-chain", marsChainID, "--src-port", "transfer", "--src-channel", "channel-0", "--amount", "100000", "--denom", "stake", "--output", "json", "--log_format", "json", "--yes")
}
