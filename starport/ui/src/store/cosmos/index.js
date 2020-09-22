import blocks from './blocks'
import ui from './ui'

export default {
  namespaced: true,
  state: {
    LOCAL_ENV: {
      // COSMOS_RPC: 'rpc.nylira.net',
      SCAFFOLD: 'localhost:8080',
      COSMOS_RPC: 'localhost:26657',
      LCD: 'localhost:1317',
      STARPORT_APP: 'localhost:12345'
    },
    backend: {
      env: {
        node_js: false,
        vue_app_custom_url: '',
      },      
      running: { 
        frontend: false,        
        rpc: false,
        api: false,
      },      
    }
  },
  getters: {
    localEnv: state => state.LOCAL_ENV,
    backendEnv: state => state.backend.env,
    backendRunningStates: state => state.backend.running
  },
  mutations: {
    /**
     * 
     * 
     * @param {object} states
     * @param {boolean} states[].frontend
     * @param {boolean} states[].rpc
     * @param {boolean} states[].api
     * 
     * 
     */      
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
    },
    /**
     * 
     * 
     * @param {object} env
     * @param {boolean} states[].node_js
     * @param {string} states[].vue_app_custom_url
     * 
     * 
     */     
    setBackendEnv(state, {
      node_js,
      vue_app_custom_url
    }) {
      state.backend.env = {
        node_js,
        vue_app_custom_url
      }
    }
  },
  actions: {},
  modules: {
    blocks,
    ui
  }
}