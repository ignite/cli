import { describe, expect, it } from 'vitest'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'
import { isDeliverTxSuccess } from '@cosmjs/stargate'

describe('custom module', async () => {
  const { Client } = await import('client')

  it('should create a list entry', async () => {
    const { account1 } = globalThis.accounts

    const mnemonic = account1['Mnemonic']
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic)
    const [account] = await wallet.getAccounts();

    const denom = 'token'
    const env = {
      denom,
	    rpcURL: globalThis.txApi,
	    apiURL: globalThis.queryApi,
    }
    const client = new Client(env, wallet)

    const entry = {
      id: '0',
      creator: account.address,
      name: "test",
    }

    // Create a new list entry
    const result = await client.ChainDisco.tx.sendMsgCreateEntry({ value: entry })

    expect(isDeliverTxSuccess(result)).toEqual(true)

    // Check that the list entry is created
    const response = await client.ChainDisco.query.queryEntryAll()

    expect(response.statusText).toEqual('OK')
    expect(response.data['Entry']).toHaveLength(1)
    expect(response.data['Entry'][0]).toEqual(entry)
  })
})
