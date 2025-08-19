//go:build !relayer

package relayer_test

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
	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/config/chain/base"
	v1 "github.com/ignite/cli/v29/ignite/config/chain/v1"
	"github.com/ignite/cli/v29/ignite/pkg/availableport"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/randstr"
	"github.com/ignite/cli/v29/ignite/pkg/xyaml"
	envtest "github.com/ignite/cli/v29/integration"
)

const (
	relayerMnemonic = "great immense still pill defense fetch pencil slow purchase symptom speed arm shoot fence have divorce cigar rapid hen vehicle pear evolve correct nerve"
)

var (
	bobName        = "bob"
	refChainConfig = v1.Config{
		Config: base.Config{
			Version: 1,
			Accounts: []base.Account{
				{
					Name:     "alice",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "slide moment original seven milk crawl help text kick fluid boring awkward doll wonder sure fragile plate grid hard next casual expire okay body",
				},
				{
					Name:     bobName,
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
			Genesis: xyaml.Map{"chain_id": randstr.Runes(12)},
		},
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
				Client: xyaml.Map{"keyring-backend": keyring.BackendTest},
				App: xyaml.Map{
					"api":      xyaml.Map{"address": ":1318"},
					"grpc":     xyaml.Map{"address": ":9092"},
					"grpc-web": xyaml.Map{"address": ":9093"},
				},
				Config: xyaml.Map{
					"p2p": xyaml.Map{"laddr": ":26658"},
					"rpc": xyaml.Map{"laddr": ":26658", "pprof_laddr": ":6061"},
				},
				Home: filepath.Join(os.TempDir(), randstr.Runes(5)),
			},
		},
	}
	hostChainConfig = v1.Config{
		Config: base.Config{
			Version: 1,
			Accounts: []base.Account{
				{
					Name:     "alice",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "slide moment original seven milk crawl help text kick fluid boring awkward doll wonder sure fragile plate grid hard next casual expire okay body",
				},
				{
					Name:     bobName,
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
			Genesis: xyaml.Map{"chain_id": randstr.Runes(12)},
		},
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
				Client: xyaml.Map{"keyring-backend": keyring.BackendTest},
				App: xyaml.Map{
					"api":      xyaml.Map{"address": ":1317"},
					"grpc":     xyaml.Map{"address": ":9090"},
					"grpc-web": xyaml.Map{"address": ":9091"},
				},
				Config: xyaml.Map{
					"p2p": xyaml.Map{"laddr": ":26656"},
					"rpc": xyaml.Map{"laddr": ":26656", "pprof_laddr": ":6060"},
				},
				Home: filepath.Join(os.TempDir(), randstr.Runes(5)),
			},
		},
	}

	nameOnRecvIbcPostPacket = "OnRecvIbcPostPacket"
	funcOnRecvIbcPostPacket = `package keeper
func (k Keeper) OnRecvIbcPostPacket(ctx context.Context, packet channeltypes.Packet, data types.IbcPostPacketData) (packetAck types.IbcPostPacketAck, err error) {
	packetAck.PostId, err = k.PostSeq.Next(ctx)
	if err != nil {
		return packetAck, err
	}
	return packetAck, k.Post.Set(ctx, packetAck.PostId, types.Post{Title: data.Title, Content: data.Content})
}`

	nameOnAcknowledgementIbcPostPacket = "OnAcknowledgementIbcPostPacket"
	funcOnAcknowledgementIbcPostPacket = `package keeper
func (k Keeper) OnAcknowledgementIbcPostPacket(ctx context.Context, packet channeltypes.Packet, data types.IbcPostPacketData, ack channeltypes.Acknowledgement) error {
    switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		// We will not treat acknowledgment error in this tutorial
		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.IbcPostPacketAck
		if err := k.cdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		seq, err := k.SentPostSeq.Next(ctx)
		if err != nil {
			return err
		}

		return k.SentPost.Set(ctx, seq,
			types.SentPost{
				PostId: packetAck.PostId,
				Title:  data.Title,
				Chain:  packet.DestinationPort + "-" + packet.DestinationChannel,
			},
		)
	default:
		return errors.New("the counter-party module does not implement the correct acknowledgment format")
	}
}`

	nameOnTimeoutIbcPostPacket = "OnTimeoutIbcPostPacket"
	funcOnTimeoutIbcPostPacket = `package keeper
func (k Keeper) OnTimeoutIbcPostPacket(ctx context.Context, packet channeltypes.Packet, data types.IbcPostPacketData) error {
	seq, err := k.TimeoutPostSeq.Next(ctx)
	if err != nil {
		return err
	}

	return k.TimeoutPost.Set(ctx, seq,
		types.TimeoutPost{
			Title: data.Title,
			Chain: packet.DestinationPort + "-" + packet.DestinationChannel,
		},
	)
}`
)

type (
	QueryChannels struct {
		Channels []struct {
			ChannelID      string   `json:"channel_id"`
			ConnectionHops []string `json:"connection_hops"`
			Counterparty   struct {
				ChannelID string `json:"channel_id"`
				PortID    string `json:"port_id"`
			} `json:"counterparty"`
			Ordering string `json:"ordering"`
			PortID   string `json:"port_id"`
			State    string `json:"state"`
			Version  string `json:"version"`
		} `json:"channels"`
	}

	QueryBalances struct {
		Balances sdk.Coins `json:"balances"`
	}
)

func runChain(
	ctx context.Context,
	t *testing.T,
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
	cfg.Validators[0].App["api"] = xyaml.Map{"address": genAddr(ports[1])}
	cfg.Validators[0].App["grpc"] = xyaml.Map{"address": genAddr(ports[2])}
	cfg.Validators[0].App["grpc-web"] = xyaml.Map{"address": genAddr(ports[3])}
	cfg.Validators[0].Config["p2p"] = xyaml.Map{"laddr": genAddr(ports[4])}
	cfg.Validators[0].Config["rpc"] = xyaml.Map{
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
		app.MustServe(ctx)
	}()

	genHTTPAddr := func(port uint) string {
		return fmt.Sprintf("http://127.0.0.1:%d", port)
	}
	return genHTTPAddr(ports[1]), genHTTPAddr(ports[5]), genHTTPAddr(ports[2]), genHTTPAddr(ports[0])
}

func TestBlogIBC(t *testing.T) {
	var (
		env         = envtest.New(t)
		app         = env.ScaffoldApp("github.com/apps/blog", "--no-module")
		tmpDir      = t.TempDir()
		ctx, cancel = context.WithCancel(env.Ctx())
	)
	t.Cleanup(func() {
		cancel()
		time.Sleep(5 * time.Second)
		require.NoError(t, os.RemoveAll(tmpDir))
	})

	app.Scaffold(
		"create an IBC module",
		false,
		"module",
		"blog",
		"--ibc",
		"--require-registration",
	)

	app.Scaffold(
		"create a post type list in an IBC module",
		false,
		"list",
		"post",
		"title",
		"content",
		"--no-message",
		"--module",
		"blog",
	)

	app.Scaffold(
		"create a sentPost type list in an IBC module",
		false,
		"list",
		"sentPost",
		"postID:uint",
		"title",
		"chain",
		"--no-message",
		"--module",
		"blog",
	)

	app.Scaffold(
		"create a timeoutPost type list in an IBC module",
		false,
		"list",
		"timeoutPost",
		"title",
		"chain",
		"--no-message",
		"--module",
		"blog",
	)

	app.Scaffold(
		"create a ibcPost package in an IBC module",
		false,
		"packet",
		"ibcPost",
		"title",
		"content",
		"--ack",
		"postID:uint",
		"--module",
		"blog",
	)

	blogKeeperPath := filepath.Join(app.SourcePath(), "x/blog/keeper")
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
	hostChainAPI, hostChainRPC, hostChainGRPC, hostChainFaucet := runChain(ctx, t, app, hostChainConfig, tmpDir, ports[:7])
	hostChainChainID := hostChainConfig.Genesis["chain_id"].(string)
	hostChainHome := hostChainConfig.Validators[0].Home
	refChainAPI, refChainRPC, refChainGRPC, refChainFaucet := runChain(ctx, t, app, refChainConfig, tmpDir, ports[7:])
	refChainChainID := refChainConfig.Genesis["chain_id"].(string)
	refChainHome := refChainConfig.Validators[0].Home

	// check the chains is up
	app.WaitChainUp(ctx, hostChainAPI)
	app.WaitChainUp(ctx, refChainAPI)

	// ibc relayer.
	env.Must(env.Exec("install the hermes relayer app",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"app",
				"install",
				"-g",
				// filepath.Join(goenv.GoPath(), "src/github.com/ignite/apps/hermes"), // Local path for test proposals
				"github.com/ignite/apps/hermes@hermes/v0.2.8",
			),
		)),
	))

	env.Must(env.Exec("configure the hermes relayer app",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"relayer",
				"hermes",
				"configure",
				hostChainChainID,
				hostChainRPC,
				hostChainGRPC,
				refChainChainID,
				refChainRPC,
				refChainGRPC,
				"--chain-a-faucet", hostChainFaucet,
				"--chain-b-faucet", refChainFaucet,
				"--generate-wallets",
				"--overwrite-config",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	go func() {
		env.Must(env.Exec("run the hermes relayer",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp,
					"relayer",
					"hermes",
					"start",
					hostChainChainID,
					refChainChainID,
				),
				step.Workdir(app.SourcePath()),
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
			step.Stderr(queryOutput),
			step.Exec(
				app.Binary(),
				"q",
				"ibc",
				"channel",
				"channels",
				"--node", hostChainRPC,
				"--log_format", "json",
				"--output", "json",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(queryOutput.Bytes(), &queryResponse); err != nil {
					return errors.Errorf("unmarshling tx response: %w", err)
				}
				if len(queryResponse.Channels) == 0 ||
					len(queryResponse.Channels[0].ConnectionHops) == 0 {
					return errors.Errorf("channel not found")
				}
				if queryResponse.Channels[0].State != "STATE_OPEN" {
					return errors.Errorf("channel is not open")
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
			Code   int    `json:"code"`
			RawLog string `json:"raw_log"`
			TxHash string `json:"txhash"`
		}
	)

	stepsTx := step.NewSteps(
		step.New(
			step.Stdout(txOutput),
			step.Stderr(txOutput),
			step.PreExec(func() error {
				txOutput.Reset()
				return nil
			}),
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
				"--node", hostChainRPC,
				"--home", hostChainHome,
				"--chain-id", hostChainChainID,
				"--output", "json",
				"--log_format", "json",
				"--keyring-backend", "test",
				"--yes",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				output := txOutput.Bytes()
				if err := json.Unmarshal(txOutput.Bytes(), &txResponse); err != nil {
					return errors.Errorf("unmarshalling tx response error: %w, response: %s", err, string(output))
				}

				time.Sleep(4 * time.Second)

				return cmdrunner.New().Run(ctx, step.New(
					step.Exec(
						app.Binary(),
						"q",
						"tx",
						txResponse.TxHash,
						"--node", hostChainRPC,
						"--home", hostChainHome,
						"--output", "json",
						"--log_format", "json",
					),
					step.Stdout(txOutput),
					step.Stderr(txOutput),
					step.PreExec(func() error {
						txOutput.Reset()
						return nil
					}),
					step.PostExec(func(execErr error) error {
						if execErr != nil {
							return execErr
						}
						output := txOutput.Bytes()
						if err := json.Unmarshal(output, &txResponse); err != nil {
							return errors.Errorf("unmarshalling tx response error: %w, response: %s", err, string(output))
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
	steps := step.NewSteps(
		step.New(
			step.Stdout(balanceOutput),
			step.Stderr(balanceOutput),
			step.Exec(
				app.Binary(),
				"q",
				"bank",
				"balances",
				receiverAddr,
				"--node", refChainRPC,
				"--home", refChainHome,
				"--log_format", "json",
				"--output", "json",
			),
			step.PreExec(func() error {
				balanceOutput.Reset()
				return nil
			}),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}

				output := balanceOutput.Bytes()
				if err := json.Unmarshal(output, &balanceResponse); err != nil {
					return errors.Errorf("unmarshalling query response error: %w, response: %s", err, string(output))
				}
				if balanceResponse.Balances.Empty() {
					return errors.Errorf("empty balances")
				}
				if !strings.HasPrefix(balanceResponse.Balances[0].Denom, "ibc/") {
					return errors.Errorf("invalid ibc balance: %v", balanceResponse.Balances[0])
				}

				return nil
			}),
		),
	)
	env.Must(env.Exec("check ibc balance", steps, envtest.ExecRetry()))
}
