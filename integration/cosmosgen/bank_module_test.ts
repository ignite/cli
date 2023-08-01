import { describe, expect, it } from "vitest";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { isDeliverTxSuccess } from "@cosmjs/stargate";

describe("bank module", async () => {
  const { Client } = await import("client");

  it("should transfer to two different addresses", async () => {
    const { account1, account2, account3 } = globalThis.accounts;

    const mnemonic = account1["Mnemonic"];
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);
    const [account] = await wallet.getAccounts();

    const denom = "token";
    const env = {
      denom,
      rpcURL: globalThis.txApi,
      apiURL: globalThis.queryApi,
    };
    const client = new Client(env, wallet);

    const toAddresses = [account2["Address"], account3["Address"]];

    // Both accounts start with 100token before the transfer
    const result = await client.signAndBroadcast([
      client.CosmosBankV1Beta1.tx.msgSend({
        value: {
          fromAddress: account.address,
          toAddress: toAddresses[0],
          amount: [{ denom, amount: "100" }],
        },
      }),
      client.CosmosBankV1Beta1.tx.msgSend({
        value: {
          fromAddress: account.address,
          toAddress: toAddresses[1],
          amount: [{ denom, amount: "200" }],
        },
      }),
    ]);

    expect(isDeliverTxSuccess(result)).toEqual(true);

    // Check that the transfers were successful
    const cases = [
      { address: toAddresses[0], wantAmount: "200" },
      { address: toAddresses[1], wantAmount: "300" },
    ];

    for (let tc of cases) {
      let response = await client.CosmosBankV1Beta1.query.queryBalance(
        tc.address,
        { denom }
      );

      expect(response.statusText).toEqual("OK");
      expect(response.data.balance.amount).toEqual(tc.wantAmount);
    }
  });
});
