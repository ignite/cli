import blocks from './blocks'
import transactions from './transactions'
import ui from './ui'

const {
  VUE_APP_CUSTOM_URL,
  VUE_APP_API_COSMOS,
  VUE_APP_API_TENDERMINT,
  VUE_APP_WS_TENDERMINT,
  VUE_APP_ADDRESS_PREFIX
} = process.env

export default {
  namespaced: true,
  state: {
    APP_ENV: {      
      GITPOD: '',
      STARPORT_APP: VUE_APP_CUSTOM_URL ? '' : 'http://localhost:12345',
      FRONTEND: '',
      RPC: '',
      API: '',
      WS: '',
      ADDR_PREFIX: ''
    },
    backend: {
      env: {
        node_js: false,
        vue_app_custom_url: process.env.VUE_APP_CUSTOM_URL,
      },      
      running: { 
        frontend: false,        
        rpc: false,
        api: false,
      },      
    }
  },
  getters: {
    appEnv: state => state.APP_ENV,
    backendEnv: state => state.backend.env,
    backendRunningStates: state => state.backend.running
  },
  mutations: {
    setAppEnv(state, { customUrl }) {
      const GITPOD = customUrl && new URL(customUrl)

      state.APP_ENV.STARPORT_APP =
        (GITPOD && `${GITPOD.protocol}//12345-${GITPOD.hostname}`) ||
        'http://localhost:12345'

      state.APP_ENV.FRONTEND =
        (GITPOD && `${GITPOD.protocol}//8080-${GITPOD.hostname}`) ||
        'http://localhost:8080'

      state.APP_ENV.API =
        VUE_APP_API_COSMOS ||
        (GITPOD && `${GITPOD.protocol}//1317-${GITPOD.hostname}`) ||
        'http://localhost:1317'

      state.APP_ENV.RPC =
        VUE_APP_API_TENDERMINT ||
        (GITPOD && `${GITPOD.protocol}//26657-${GITPOD.hostname}`) ||
        'http://localhost:26657'

      state.APP_ENV.WS =
        VUE_APP_WS_TENDERMINT ||
        (GITPOD && `wss://26657-${GITPOD.hostname}/websocket`) ||
        'ws://localhost:26657/websocket'

      state.APP_ENV.ADDR_PREFIX = VUE_APP_ADDRESS_PREFIX || 'cosmos'      
    },
    /**
     * 
     * 
     * @param {object} state
     * @param {object} states
     * @param {boolean} states.frontend
     * @param {boolean} states.rpc
     * @param {boolean} states.api
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
     * @param {object} state
     * @param {object} env
     * @param {boolean} states.node_js
     * @param {string} states.vue_app_custom_url
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
  modules: {
    blocks,
    transactions,
    ui
  }
}