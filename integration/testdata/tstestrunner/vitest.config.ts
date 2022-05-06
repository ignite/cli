import path from 'path'
import { readdirSync, Dirent } from 'fs'
import { defineConfig } from 'vitest/config'

type Alias = {
  [key: string]: string
}

const collectModuleAliases = (alias: Alias, modulesPath: string) => {
  readdirSync(modulesPath, { withFileTypes: true }).forEach((item: Dirent) => {
    if (item.name.startsWith('.') || item.isFile()) {
      return
    }

    alias[item.name] = path.join(modulesPath, item.name)
  });
}

// Absolute path to the blockchain app directory
const chainPath = process.env.TEST_CHAIN_PATH

// The module aliases are used to be able to import generated code within the tests
const alias: Alias = {}

// Collect the module aliases for the chain app
collectModuleAliases(alias, `${chainPath}/vue/node_modules`)

// Collect the module aliases for the generated Vuex store
collectModuleAliases(alias, `${chainPath}/vue/src/store/generated/cosmos/cosmos-sdk`)
collectModuleAliases(alias, `${chainPath}/vue/src/store/generated/cosmos/ibc-go`)

export default defineConfig({
  test: {
    include: ['**/*_test.ts'],
    setupFiles: 'setupTest.ts'
  },
  resolve: {
    alias
  }
})
