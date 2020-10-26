import axios from 'axios'

const {
  VUE_APP_CUSTOM_URL,
  VUE_APP_API_COSMOS,
  VUE_APP_API_TENDERMINT,
  VUE_APP_WS_TENDERMINT,
  VUE_APP_ADDRESS_PREFIX
} = process.env

const state = {
  timer: null,
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
    prevStates: {
      frontend: null,
      rpc: null,
      api: null,
    }    
  }
}

const getters = {
  appEnv: state =>
    state.APP_ENV,
  backendEnv: state =>
    state.backend.env,
  backendRunningStates: state =>
    state.backend.running,
  wasAppRestarted: state => status => {
    return (state.backend.prevStates.rpc !== null && state.backend.prevStates.api !== null) &&
      (!state.backend.prevStates.rpc && !state.backend.prevStates.api) &&
      (status.is_consensus_engine_alive && status.is_my_app_backend_alive)
  },    
}

export default {
  namespaced: true,
  state,
  getters,
  mutations: {
    setAppEnv(state, {
      customUrl
    }) {
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
    },
    setPrevStates(state, {
      status
    }) {
      state.backend.prevStates = {
        frontend: status ? status.is_my_app_frontend_alive : false,
        rpc: status ? status.is_consensus_engine_alive : false,
        api: status ? status.is_my_app_backend_alive : false,
      }
    },    
    setTimer(state, {
      timer
    }) {
      state.timer = timer
    },
    clearTimer(state) {
      clearInterval(state.timer)
    }
  },
  actions: {
    async setStatusState({ getters, commit }) {
      try {
        const { data } = await axios.get(`${getters.appEnv.STARPORT_APP}/status`)
        const { status, env } = data

        commit('setAppEnv', { 
          customUrl: env.vue_app_custom_url
        })
        commit('setBackendRunningStates', {
          frontend: status.is_my_app_frontend_alive,
          rpc: status.is_consensus_engine_alive,
          api: status.is_my_app_backend_alive,
        })
        commit('setBackendEnv', {
          node_js: env.node_js,
          vue_app_custom_url: env.vue_app_custom_url
        })

        /**
         * 
         // If backend was down, but alive now,
         // it indicates the app is restarting.
         // Forcing browser to reload in this case to reset blockchain data.
         * 
         */        
        if (getters.wasAppRestarted(status)) {
          window.location.reload(false)
        }
        commit('setPrevStates', {
          status
        })
      } catch {
        commit('setBackendRunningStates', {
          frontend: false,
          rpc: false,
          api: false,
        })    
        
        commit('setPrevStates', {
          status: null
        })
      }
    }    
  }
}