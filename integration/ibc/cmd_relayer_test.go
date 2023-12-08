//go:build !relayer

package ibc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v28/ignite/pkg/xurl"
	envtest "github.com/ignite/cli/v28/integration"
)

var (
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

func TestBlogIBC(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/planet")
		servers = app.RandomizeServerPorts()
		ctx     = env.Ctx()
	)

	nodeAddr, err := xurl.TCP(servers.RPC)
	if err != nil {
		t.Fatalf("cant read nodeAddr from host.RPC %v: %v", servers.RPC, err)
	}

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

	// serve both chains
	ctxEarth, cancelEarth := context.WithCancel(ctx)
	go func() {
		defer cancelEarth()
		env.Must(app.Serve("should serve earth", envtest.ExecCtx(ctxEarth)))
	}()
	ctxMars, cancelMars := context.WithCancel(ctx)
	go func() {
		defer cancelMars()
		env.Must(app.Serve("should serve mars", envtest.ExecCtx(ctxMars)))
	}()

	// configure and run the ts relayer.
	env.Must(env.Exec("configure the ts relayer",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"relayer",
				"configure", "-a",
				"--source-rpc", "http://0.0.0.0:26657",
				"--source-faucet", "http://0.0.0.0:4500",
				"--source-port", "planet",
				"--source-version", "earth-1",
				"--source-gasprice", "0.0000025stake",
				"--source-prefix", "cosmos",
				"--source-gaslimit", "300000",
				"--target-rpc", "http://0.0.0.0:26659",
				"--target-faucet", "http://0.0.0.0:4501",
				"--target-port", "planet",
				"--target-version", "mars-1",
				"--target-gasprice", "0.0000025stake",
				"--target-prefix", "cosmos",
				"--target-gaslimit", "300000",
			),
			step.Workdir(app.SourcePath()),
		)),
	))
	go func() {
		env.Must(env.Exec("run the ts relayer",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, "relayer", "connect"),
				step.Workdir(app.SourcePath()),
			)),
		))
	}()

	// check the chains is up
	stepsCheck := step.NewSteps(
		step.New(
			step.Exec(
				app.Binary(),
				"config",
				"output", "json",
			),
			step.PreExec(func() error {
				// todo set chain configs
				if err := env.IsAppServed(ctx, servers.API); err != nil {
					return err
				}
				return env.IsAppServed(ctx, servers.API)
			}),
		),
	)
	env.Exec("waiting the chain is up", stepsCheck, envtest.ExecRetry())

	var (
		output     = &bytes.Buffer{}
		txResponse struct {
			Code   int
			RawLog string `json:"raw_log"`
		}
	)
	// sign tx to add an item to the list.
	stepsTx := step.NewSteps(
		step.New(
			step.Stdout(output),
			step.PreExec(func() error {
				err := env.IsAppServed(ctx, servers.API)
				return err
			}),
			step.Exec(
				app.Binary(),
				"tx",
				"blog",
				"send-ibc-post",
				"channel-0",
				"Hello",
				"Hello Mars, I'm Alice from Earth",
				"--chain-id", "blog",
				"--from", "alice",
				"--node", nodeAddr,
				"--output", "json",
				"--log_format", "json",
				"--yes",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				err := json.Unmarshal(output.Bytes(), &txResponse)
				if err != nil {
					return fmt.Errorf("unmarshling tx response: %w", err)
				}
				return nil
			}),
		),
	)
	if !env.Exec("sign a tx", stepsTx, envtest.ExecRetry()) {
		t.FailNow()
	}
	require.Equal(t, 0, txResponse.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)
}
