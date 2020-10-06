import { mapGetters, mapMutations, mapActions } from 'vuex'

import blockConnection from './blocks/initConnection'
import txConnection from './transactions/init'
import backendConnection from './backend/initConnection'

export default {
  mixins: [backendConnection, blockConnection, txConnection],
  computed: {
    ...mapGetters('cosmos/ui', [ 'targetTable', 'blocksExplorerTableId' ]),
  },
  methods: {
    ...mapMutations('cosmos/ui', [ 'createTable' ])
  },
  created() {
    if (!this.targetTable(this.blocksExplorerTableId)) {
      this.createTable(this.blocksExplorerTableId)
    }
  }
}