import { defineConfig } from 'vitest/config'

import aliases from './aliases'

// TODO: add .env file support for a better developer experience ? It would allow
// writting new tests agains a running blockchain without the need of scaffolding
// a new one for each test run.

// Collect the module aliases for the generated code, including cosmos, tendermint
// or user generated modules, and also the dependencies for the frontend client.
// Module aliases are used to be able to import auto generated code within the tests.
const alias = aliases.collect(process.env.TEST_CHAIN_PATH)

export default defineConfig({
  test: {
    include: ['**/*_test.ts'],
    globals: true,
    setupFiles: 'setupTest.ts'
  },
  resolve: {
    alias
  }
})
