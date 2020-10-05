import { mapGetters, mapActions, mapMutations } from 'vuex'

export default {
  computed: {
    ...mapGetters('cosmos', [ 'appEnv' ]),
  },
  methods: {
    ...mapActions('cosmos/blocks', [ 'initBlockConnection' ]),
  },
  created() {
    this.initBlockConnection({ appEnv: this.appEnv })
  }
}