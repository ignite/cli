import { beforeAll, describe, expect, it } from 'vitest'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'
import { isDeliverTxSuccess } from '@cosmjs/stargate'

import { txClient, queryClient } from 'cosmos-bank-v1beta1-js/module'

describe('Bank', () => {
  let txApi: string
  let queryApi: string

  beforeAll(() => {
    txApi = process.env.TEST_TX_API || ''
    queryApi = process.env.TEST_QUERY_API || ''

    expect(txApi, 'TEST_TX_API is required').not.toEqual('')
    expect(queryApi, 'TEST_QUERY_API is required').not.toEqual('')
  })

  it('transfers to two different addresses', async () => {
    const denom = 'token'
    const toAddresses = [
      'cosmos19yy9sf00k00cjcwh532haeq8s63uhdy7qs5m2n',
      'cosmos10957ee377t2xpwyt4mlpedjldp592h0ylt8uz7',
    ]

    // TODO: should we send values from the chain integration test (mnemonic, addresses, ...) ?
    const mnemonic = 'toe mail light plug pact length excess predict real artwork laundry when steel online adapt clutch debate vehicle dash alter rifle virtual season almost'
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic)
    const [account] = await wallet.getAccounts();
    const tx = await txClient(wallet, { addr: txApi })

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

    const query = await queryClient({ addr: queryApi })
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
