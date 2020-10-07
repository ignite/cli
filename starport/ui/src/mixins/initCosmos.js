import { mapGetters, mapMutations, mapActions } from 'vuex'

import blockConnection from './blocks/initConnection'
import backendConnection from './backend/initConnection'

export default {
  mixins: [backendConnection, blockConnection],
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