import blocks from './blocks'
import ui from './ui'

export default {
  namespaced: true,
  state: {
    LOCAL_ENV: {
      // COSMOS_RPC: 'rpc.nylira.net',
      COSMOS_RPC: 'localhost:26657',
      LCD: 'localhost:1317'
    },
    backend: {
      running: { // ⚠️ temporarily set to true for local dev
        frontend: true,        
        rpc: true,
        api: true,
      },      
    }
  },
  getters: {
    localEnv: state => state.LOCAL_ENV,
    backendRunningStates: state => state.backend.running
  },
  mutations: {
    setBackendRunningStates(state, {
      frontend,
      rpc,
      api
    }) {
      state.backend.running = {
        frontend,
        rpc,
        api
      }
    }
  },
  actions: {},
  modules: {
    blocks,
    ui
  }
}