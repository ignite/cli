//go:build !relayer

package ibc_test

import (
	"fmt"
	"github.com/ignite/cli/v28/ignite/pkg/protoanalysis/protoutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/goanalysis"
	envtest "github.com/ignite/cli/v28/integration"
)

var (
	nameSendIbcPost = "SendIbcPost"
	funcSendIbcPost = `package keeper
func (k msgServer) SendIbcPost(goCtx context.Context, msg *types.MsgSendIbcPost) (*types.MsgSendIbcPostResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: logic before transmitting the packet

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
    if err != nil {
        return nil, err
    }

    return &types.MsgSendIbcPostResponse{}, nil
}`

	nameOnRecvIbcPostPacket = "OnRecvIbcPostPacket"
	funcOnRecvIbcPostPacket = `package keeper
func (k Keeper) OnRecvIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) (packetAck types.IbcPostPacketAck, err error) {
    // validate packet data upon receiving
    if err := data.ValidateBasic(); err != nil {
        return packetAck, err
    }

    id := k.AppendPost(
        ctx,
        types.Post{
            Title:   data.Title,
            Content: data.Content,
        },
    )

    packetAck.PostId = strconv.FormatUint(id, 10)

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

        k.AppendSentPost(
            ctx,
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
    k.AppendTimedoutPost(
        ctx,
        types.TimedoutPost{
            Title:   data.Title,
            Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
        },
    )

    return nil
}`
)

func TestBlogIBC(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/planet")
	)

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
				"creator",
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
				"postID",
				"title",
				"chain",
				"creator",
				"--no-message",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a timedoutPost type list in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"timedoutPost",
				"title",
				"chain",
				"creator",
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
				"postID",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	//protoPacketPath := filepath.Join(app.SourcePath(), "proto/planet/blog/packet.proto")
	//require.NoError(t, addCreatorToPacket(protoPacketPath))

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

	app.EnsureSteady()
}

func addCreatorToPacket(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	protoFile, err := protoutil.ParseProtoFile(f)
	if err != nil {
		return err
	}

	IbcPostPacketData, err := protoutil.GetMessageByName(protoFile, "IbcPostPacketData")
	if err != nil {
		return fmt.Errorf("couldn't find message 'IbcPostPacketData' in %s: %w", path, err)
	}
	field := protoutil.NewField("creator", "string", protoutil.NextUniqueID(IbcPostPacketData))
	protoutil.Append(IbcPostPacketData, field)

	p := protoutil.Print(protoFile)
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f, "%d", len(p))
	if err != nil {
		return err
	}
	return nil
}
