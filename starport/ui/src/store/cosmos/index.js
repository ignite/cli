import env from './env'
import blocks from './blocks'
import txs from './txs'

export default {
  namespaced: true,
  modules: {
    blocks,
    txs,
    env
  }
}