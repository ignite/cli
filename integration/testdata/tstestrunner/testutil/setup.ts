import { beforeAll, expect } from 'vitest'

// Make sure that the tests have fetch API support
import 'isomorphic-unfetch'

type Account = {
  name: string;
  address: string;
  mnemonic: string;
  coins: string[];
}

type GlobalAccounts = {
  [name: string]: Account
}

beforeAll(() => {
    // Initialize required globals
    globalThis.txApi = process.env.TEST_TX_API || ''
    globalThis.queryApi = process.env.TEST_QUERY_API || ''

    expect(globalThis.txApi, 'TEST_TX_API is required').not.toEqual('')
    expect(globalThis.queryApi, 'TEST_QUERY_API is required').not.toEqual('')

    // Initialize the global accounts
    globalThis.accounts = <GlobalAccounts>{}

    JSON.parse(process.env.TEST_ACCOUNTS || '[]').forEach((account: Account) => {
      const name = account['Name']

      globalThis.accounts[name] = account
    })
})
