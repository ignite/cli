import axios from 'axios'
import { mapGetters, mapMutations } from 'vuex'

export default {
  data() {
    return {
      timer: null,
      prevStates: {
        frontend: null,
        rpc: null,
        api: null,
      }
    }
  },
  computed: {
    ...mapGetters('cosmos', [ 'appEnv', 'backendEnv' ]),
  },
  methods: {
    ...mapMutations('cosmos', [ 'setBackendRunningStates', 'setBackendEnv', 'setAppEnv' ]),
    wasAppRestarted(status) {
      return ((this.prevStates.rpc !== null && this.prevStates.api !== null) &&
        (!this.prevStates.rpc && !this.prevStates.api) &&
        (status.is_consensus_engine_alive && status.is_my_app_backend_alive))
    },
    setPrevStates(status) {
      this.prevStates = {
        frontend: status ? status.is_my_app_frontend_alive : false,
        rpc: status ? status.is_consensus_engine_alive : false,
        api: status ? status.is_my_app_backend_alive : false,
      }
    },
    async setStatusState() {
      try {
        const { data } = await axios.get(`${this.appEnv.STARPORT_APP}/status`)
        const { status, env } = data

        this.setAppEnv({
          customUrl: env.vue_app_custom_url 
        })
        this.setBackendRunningStates({
          frontend: status.is_my_app_frontend_alive,
          rpc: status.is_consensus_engine_alive,
          api: status.is_my_app_backend_alive,
        })
        this.setBackendEnv({
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
        if (this.wasAppRestarted(status)) {
          window.location.reload(false)
        }
        this.setPrevStates(status)
      } catch {
        this.setBackendRunningStates({
          frontend: false,
          rpc: false,
          api: false,
        })    
        
        this.setPrevStates(null)        
      }
    }
  },
  async created() {
    /*
     *
     // Fetch backend status regularly
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