import axios from 'axios'
import { mapMutations } from 'vuex'

export default {
  data() {
    return {
      TEMP_LOCALHOST: 'http://localhost:12345', // ⚠️ CORS blocker
      env: {
        chain_id: '{{chain_id}}',
        node_js: false,
        vue_app_custom_url: '',
      },      
      timer: null,
    }
  },
  methods: {
    ...mapMutations('cosmos', [ 'setBackendRunningStates' ]),
    async setStatusState() {
      try {
        const { data } = await axios.get(`${this.TEMP_LOCALHOST}/status`)
        const { status, env } = data
        this.setBackendRunningStates({
          frontend: status.is_my_app_frontend_alive,
          rpc: status.is_consensus_engine_alive,
          api: status.is_my_app_backend_alive,
        })
        // this.running = {
        //   rpc: status.is_consensus_engine_alive,
        //   api: status.is_my_app_backend_alive,
        //   frontend: status.is_my_app_frontend_alive,
        // }
        this.env = env
      } catch {
        this.setBackendRunningStates({
          frontend: false,
          rpc: false,
          api: false,
        })        
        // this.running = {
        //   rpc: false,
        //   api: false,
        //   frontend: false,
        // }
      }
    },
  },
  async created() {
    // this.timer = setInterval(this.setStatusState.bind(this), 5000)
    // try {
    //   await this.setStatusState()
    // } catch {
    //   console.log(`Can't fetch /env`)
    // }
  },
  beforeDestroy() {
    clearInterval(this.timer)
  }
}