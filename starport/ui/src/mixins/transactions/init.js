import { mapActions } from 'vuex'

export default {
  methods: {
    ...mapActions('cosmos/transactions', [ 'initTxsStack' ])
  },
  created() {
    this.initTxsStack()
  }
}