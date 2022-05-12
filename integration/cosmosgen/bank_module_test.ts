import { describe, expect, it } from 'vitest'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'
import { isDeliverTxSuccess } from '@cosmjs/stargate'

describe('bank module', async () => {
  const { txClient, queryClient } = await import('cosmos-bank-v1beta1-js/module')

  it('should transfer to two different addresses', async () => {
    const { account1, account2, account3 } = global.accounts

    const mnemonic = account1['Mnemonic']
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic)
    const [account] = await wallet.getAccounts();
    const tx = await txClient(wallet, { addr: global.txApi })

    const denom = 'token'
    const toAddresses = [
      account2['Address'],
      account3['Address'],
    ]

    // Both accounts start with 100token before the transfer
    const result = await tx.signAndBroadcast([
      tx.msgSend({
        from_address: account.address,
        to_address: toAddresses[0],
        amount: [{ denom, amount: '100' }],
      }),
      tx.msgSend({
        from_address: account.address,
        to_address: toAddresses[1],
        amount: [{ denom, amount: '200' }],
      }),
    ])

    expect(isDeliverTxSuccess(result)).toEqual(true)

    // Check that the transfers were successful
    const query = await queryClient({ addr: global.queryApi })
    const cases = [
      { address: toAddresses[0], wantAmount: '200' },
      { address: toAddresses[1], wantAmount: '300' },
    ]

    for (let tc of cases) {
      let response = await query.queryBalance(tc.address, { denom })

      expect(response.statusText).toEqual('OK')
      expect(response.data.balance.amount).toEqual(tc.wantAmount)
    }
  })
})
