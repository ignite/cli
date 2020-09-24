import axios from 'axios'
import { mapGetters, mapMutations } from 'vuex'

export default {
  data() {
    return {
      timer: null,
    }
  },
  computed: {
    ...mapGetters('cosmos', [ 'appEnv', 'backendEnv' ]),
  },
  methods: {
    ...mapMutations('cosmos', [ 'setBackendRunningStates', 'setBackendEnv', 'setAppEnv' ]),
    async setStatusState() {
      try {
        const { data } = await axios.get(`${this.appEnv.STARPORT_APP}/status`)
        const { status, env } = data

        this.setBackendRunningStates({
          frontend: status.is_my_app_frontend_alive,
          rpc: status.is_consensus_engine_alive,
          api: status.is_my_app_backend_alive,
        })

        this.setBackendEnv({
          node_js: env.node_js,
          vue_app_custom_url: env.vue_app_custom_url
        })
      } catch {
        this.setBackendRunningStates({
          frontend: false,
          rpc: false,
          api: false,
        })        
      }
    },
    async getAppEnvs() {
      const { GITPOD, FRONTEND, RPC, API, WS, ADDR_PREFIX } = this.appEnv
      const { VUE_APP_API_COSMOS, VUE_APP_WS_TENDERMINT, VUE_APP_API_TENDERMINT, VUE_APP_ADDRESS_PREFIX } = process.env

      const fmtAPI =
        VUE_APP_API_COSMOS ||
        (GITPOD && `${GITPOD.protocol}//1317-${GITPOD.hostname}`) ||
        API
      const fmtRPC =
        VUE_APP_API_TENDERMINT ||
        (GITPOD && `${GITPOD.protocol}//26657-${GITPOD.hostname}`) ||
        RPC
      const fmtWS =
        VUE_APP_WS_TENDERMINT ||
        (GITPOD && `wss://26657-${GITPOD.hostname}/websocket`) ||
        WS
      const fmtADDR_PREFIX = 
        VUE_APP_ADDRESS_PREFIX || 
        ADDR_PREFIX
      
      return {
        RPC: fmtRPC,
        API: fmtAPI,
        WS: fmtWS,
        ADDR_PREFIX: fmtADDR_PREFIX
      }
    }
  },
  async created() {
    /*
     *
     // 1. Set global app variables
     *
     */
    await this.setAppEnv(this.getAppEnvs())    

    /*
     *
     // 2. Fetch backend status regularly
     *
     */
    this.timer = setInterval(this.setStatusState.bind(this), 5000)
    try {
      await this.setStatusState()
    } catch {
      console.log(`Can't fetch /env`)
    }
  },
  beforeDestroy() {
    clearInterval(this.timer)
  }
}