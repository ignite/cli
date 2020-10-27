import env from '@tendermint/vue/src/store/env'
import blocks from '@tendermint/vue/src/store/blocks'
import txs from '@tendermint/vue/src/store/txs'

export default {
  namespaced: true,
  modules: {
    blocks,
    txs,
    env
  }
}