import { mapGetters, mapActions, mapMutations } from 'vuex'

export default {
  computed: {
    ...mapGetters('cosmos', [ 'localEnv' ]),
  },
  methods: {
    ...mapActions('cosmos/blocks', [ 'initBlockConnection' ]),
  },
  created() {
    this.initBlockConnection({ localEnv: this.localEnv })
  }
}