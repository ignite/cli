import { mapGetters, mapActions, mapMutations } from 'vuex'

export default {
  methods: {
    ...mapActions('cosmos/blocks', [ 'initBlockConnection' ]),
  },
  created() {
    this.initBlockConnection()
  }
}