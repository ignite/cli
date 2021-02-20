import { coins } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import Types from "{{ .TypesPath }}";
import { Api } from "{{ .RESTPath }}";

const types = [
  {{ range .Module.Msgs }}["/{{ .URI }}", Types.{{ .URI }}],
  {{ end }}
];

const registry = new Registry(types);

const txClient = async (wallet, { addr: addr } = { addr: "http://localhost:26657" }) => {
  if (!wallet) throw new Error("wallet is required");

  const client = await SigningStargateClient.connectWithWallet(addr, wallet, { registry });
  const { address } = wallet;
  const fee = {
    amount: coins(0, "token"),
    gas: "200000",
  };

  return {
    signAndBroadcast: (msgs) => client.signAndBroadcast(address, msgs, fee),
    {{ range .Module.Msgs }}{{ camelCase .Name }}: (data) => ({ typeUrl: "/{{ .URI }}", value: data }),
    {{ end }}
  };
};

const queryClient = async ({ addr: addr } = { addr: "http://localhost:1317" }) => {
  return new Api({ baseUrl: addr }).{{ .Module.Name }};
};

export {
  txClient,
  queryClient,
};
