import blocks from './blocks'
import ui from './ui'

export default {
  namespaced: true,
  state: {
    APP_ENV: {
      GITPOD: process.env.VUE_APP_CUSTOM_URL && new URL(process.env.VUE_APP_CUSTOM_URL),
      FRONTEND: 'http://localhost:8080',
      RPC: 'http://localhost:26657',
      API: 'http://localhost:1317',
      WS: 'ws://localhost:26657/websocket',
      STARPORT_APP: 'http://localhost:12345',
      ADDR_PREFIX: 'cosmos'
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
    appEnv: state => state.APP_ENV,
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
    },
    /**
     * 
     * 
     * @param {object} state
     * @param {object} appEnv
     * 
     * 
     */     
    setAppEnv(state, appEnv) {
      state.APP_ENV = {
        ...state.APP_ENV,
        ...appEnv
      }
    }
  },
  actions: {},
  modules: {
    blocks,
    ui
  }
}