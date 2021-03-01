import { StdFee } from "@cosmjs/launchpad";
import { OfflineSigner, EncodeObject } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgAcknowledgement } from "./types/ibc/core/channel/v1/tx";
import { MsgRecvPacket } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelOpenInit } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelCloseInit } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelCloseConfirm } from "./types/ibc/core/channel/v1/tx";
import { MsgTimeout } from "./types/ibc/core/channel/v1/tx";
import { MsgTimeoutOnClose } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelOpenAck } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelOpenTry } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelOpenConfirm } from "./types/ibc/core/channel/v1/tx";
interface TxClientOptions {
    addr: string;
}
interface SignAndBroadcastOptions {
    fee: StdFee;
}
declare const txClient: (wallet: OfflineSigner, { addr: addr }?: TxClientOptions) => Promise<{
    signAndBroadcast: (msgs: EncodeObject[], { fee: fee }?: SignAndBroadcastOptions) => Promise<import("@cosmjs/stargate").BroadcastTxResponse>;
    msgAcknowledgement: (data: MsgAcknowledgement) => EncodeObject;
    msgRecvPacket: (data: MsgRecvPacket) => EncodeObject;
    msgChannelOpenInit: (data: MsgChannelOpenInit) => EncodeObject;
    msgChannelCloseInit: (data: MsgChannelCloseInit) => EncodeObject;
    msgChannelCloseConfirm: (data: MsgChannelCloseConfirm) => EncodeObject;
    msgTimeout: (data: MsgTimeout) => EncodeObject;
    msgTimeoutOnClose: (data: MsgTimeoutOnClose) => EncodeObject;
    msgChannelOpenAck: (data: MsgChannelOpenAck) => EncodeObject;
    msgChannelOpenTry: (data: MsgChannelOpenTry) => EncodeObject;
    msgChannelOpenConfirm: (data: MsgChannelOpenConfirm) => EncodeObject;
}>;
interface QueryClientOptions {
    addr: string;
}
declare const queryClient: ({ addr: addr }?: QueryClientOptions) => Promise<Api<unknown>>;
export { txClient, queryClient, };
