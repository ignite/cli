import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgSubmitMisbehaviour } from "./types/ibc/core/client/v1/tx";
import { MsgCreateClient } from "./types/ibc/core/client/v1/tx";
import { MsgUpdateClient } from "./types/ibc/core/client/v1/tx";
import { MsgUpgradeClient } from "./types/ibc/core/client/v1/tx";
const types = [
    ["/ibc.core.client.v1.MsgSubmitMisbehaviour", MsgSubmitMisbehaviour],
    ["/ibc.core.client.v1.MsgCreateClient", MsgCreateClient],
    ["/ibc.core.client.v1.MsgUpdateClient", MsgUpdateClient],
    ["/ibc.core.client.v1.MsgUpgradeClient", MsgUpgradeClient],
];
const registry = new Registry(types);
const defaultFee = {
    amount: [],
    gas: "200000",
};
const txClient = async (wallet, { addr: addr } = { addr: "http://localhost:26657" }) => {
    if (!wallet)
        throw new Error("wallet is required");
    const client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
    const { address } = (await wallet.getAccounts())[0];
    return {
        signAndBroadcast: (msgs, { fee: fee } = { fee: defaultFee }) => client.signAndBroadcast(address, msgs, fee),
        msgSubmitMisbehaviour: (data) => ({ typeUrl: "/ibc.core.client.v1.MsgSubmitMisbehaviour", value: data }),
        msgCreateClient: (data) => ({ typeUrl: "/ibc.core.client.v1.MsgCreateClient", value: data }),
        msgUpdateClient: (data) => ({ typeUrl: "/ibc.core.client.v1.MsgUpdateClient", value: data }),
        msgUpgradeClient: (data) => ({ typeUrl: "/ibc.core.client.v1.MsgUpgradeClient", value: data }),
    };
};
const queryClient = async ({ addr: addr } = { addr: "http://localhost:1317" }) => {
    return new Api({ baseUrl: addr });
};
export { txClient, queryClient, };
