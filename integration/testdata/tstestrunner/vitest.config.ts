import { defineConfig } from 'vitest/config'

// TODO: add .env file support for a better developer experience ? It would allow
// writting new tests agains a running blockchain without the need of scaffolding
// a new one for each test run.

export default defineConfig({
  test: {
    include: ['**/*_test.ts'],
    globals: true,
    setupFiles: 'testutil/setup.ts'
  },
  resolve: {
    alias: {
      'client': process.env.TEST_TSCLIENT_DIR
    }
  }
})
