import { StdFee } from "@cosmjs/launchpad";
import { OfflineSigner, EncodeObject } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgChannelOpenConfirm } from "./types/ibc/core/channel/v1/tx";
import { MsgTimeoutOnClose } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelOpenTry } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelCloseInit } from "./types/ibc/core/channel/v1/tx";
import { MsgRecvPacket } from "./types/ibc/core/channel/v1/tx";
import { MsgAcknowledgement } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelOpenAck } from "./types/ibc/core/channel/v1/tx";
import { MsgTimeout } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelOpenInit } from "./types/ibc/core/channel/v1/tx";
import { MsgChannelCloseConfirm } from "./types/ibc/core/channel/v1/tx";
export declare const MissingWalletError: Error;
interface TxClientOptions {
    addr: string;
}
interface SignAndBroadcastOptions {
    fee: StdFee;
    memo?: string;
}
declare const txClient: (wallet: OfflineSigner, { addr: addr }?: TxClientOptions) => Promise<{
    signAndBroadcast: (msgs: EncodeObject[], { fee, memo }?: SignAndBroadcastOptions) => Promise<import("@cosmjs/stargate").BroadcastTxResponse>;
    msgChannelOpenConfirm: (data: MsgChannelOpenConfirm) => EncodeObject;
    msgTimeoutOnClose: (data: MsgTimeoutOnClose) => EncodeObject;
    msgChannelOpenTry: (data: MsgChannelOpenTry) => EncodeObject;
    msgChannelCloseInit: (data: MsgChannelCloseInit) => EncodeObject;
    msgRecvPacket: (data: MsgRecvPacket) => EncodeObject;
    msgAcknowledgement: (data: MsgAcknowledgement) => EncodeObject;
    msgChannelOpenAck: (data: MsgChannelOpenAck) => EncodeObject;
    msgTimeout: (data: MsgTimeout) => EncodeObject;
    msgChannelOpenInit: (data: MsgChannelOpenInit) => EncodeObject;
    msgChannelCloseConfirm: (data: MsgChannelCloseConfirm) => EncodeObject;
}>;
interface QueryClientOptions {
    addr: string;
}
declare const queryClient: ({ addr: addr }?: QueryClientOptions) => Promise<Api<unknown>>;
export { txClient, queryClient, };
