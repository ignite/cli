<template>
  <div class="explorer">
    <FullWidthContainer>
      <div slot="sideSheet" class="explorer__block">
        <BlockDetailSheet :block="highlightedBlock"/>
      </div>
      <div slot="mainContent" class="explorer__chain">
        <div class="explorer__chain-header">Blocks</div>
        <div class="explorer__chain-main">
          <BlockChain :blocks="fmtBlockData" />
        </div>
      </div>      
    </FullWidthContainer>
  </div>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import blockHelpers from '@/mixins/blocks/helpers'

import _ from 'lodash'
import { uuid } from 'vue-uuid'

import FullWidthContainer from '@/components/containers/FullWidthContainer'
import BlockDetailSheet from '@/modules/BlockDetailSheet'
import BlockChain from '@/modules/BlockChain'

export default {
  components: {
    FullWidthContainer,
    BlockChain,
    BlockDetailSheet,
  },
  data() {
    return {
      blockFormatter: blockHelpers.blockFormatter(),
      states: {
        isHidingBlocksWithoutTxs: false,
        isScrolledInTopHalf: true,
        isLoading: false
      },
    }
  },
  computed: {
    /*
     *
     * Vuex 
     *
     */
    ...mapGetters('cosmos', [ 'appEnv' ]),
    ...mapGetters('cosmos/ui', [ 'targetTable', 'isTableSheetActive', 'blocksExplorerTableId' ]),
    ...mapGetters('cosmos/blocks', [ 'highlightedBlock', 'blocksStack', 'lastBlock', 'stackChainRange', 'latestBlock' ]),
    ...mapGetters('cosmos/transactions', [ 'txsStack' ]),
    /*
     *
     * Local
     * 
     */    
    fmtIsTableSheetActive() {
      return this.isTableSheetActive(this.blocksExplorerTableId)
    },    
    fmtTargetTable() {
      return this.targetTable(this.blocksExplorerTableId)
    },
    fmtBlockData() {
      const fmtBlockForTable = this.blockFormatter.blockForTable(this.blocksStack)

      if (!fmtBlockForTable) return null

      if (this.states.isHidingBlocksWithoutTxs) {
        return this.blockFormatter.filterBlock(fmtBlockForTable).hideBlocksWithoutTxs()
      }

      return fmtBlockForTable
    },
    isBlocksTableEmpty() {
      return this.blocksStack.length<=0 ||
        !this.fmtBlockData || 
        this.fmtBlockData?.length<=0
    },
    blockFilterText() {
      return !this.states.isHidingBlocksWithoutTxs
        ? 'Hide blocks without txs'
        : 'Show blocks without txs'
    },
    blockTableEmptyText() {
      return (this.blocksStack.length>=0 && this.fmtBlockData?.length<=0 && this.states.isHidingBlocksWithoutTxs) 
        ? 'Waiting for blocks with txs'
        : 'Waiting for blocks'
    }
  },  
  methods: {
    /*
     *
     * Vuex 
     *
     */    
    ...mapMutations('cosmos/ui', [ 'setTableSheetState' ]),
    ...mapMutations('cosmos/blocks', [ 'popOverloadBlocks', 'sortBlocksStack' ]),
    ...mapActions('cosmos/blocks', [ 'addBlockEntry', 'getBlockchain', 'setHighlightedBlock' ]),
    /*
     *
     * Local 
     *
     */      
    /*
     *
     // Pop overloaded blocks (over maxStackCount)
     // only when scrolling to upperhalf of the table
     *
     */         
    async handleScrollTop() {
      this.states.isScrolledInTopHalf=true      
      if (!this.latestBlock) return 
      
      const isShowingLatestBlock = (parseInt(this.latestBlock.height) === this.stackChainRange.highestHeight)

      if (!isShowingLatestBlock && !this.states.isLoading) {
        this.states.isLoading=true

        await this.getBlockchain({ 
          blockHeight: this.stackChainRange.highestHeight,
          toGetLowerBlocks: false
        }).then(() => this.states.isLoading=false)        
      }
    },
    /*
     *
     // Load extra 20 blocks
     // only when scrolling to bottom of the table
     *
     */          
    async handleScrollBottom() {
      this.states.isScrolledInTopHalf=false
      this.states.isLoading=true

      await this.getBlockchain({ 
        blockHeight: this.lastBlock.height,
        toGetLowerBlocks: true
      }).then(() => this.states.isLoading=false)
    }   
  },
  beforeDestroy() {
    if (this.latestBlock) {
      this.getBlockchain({ 
        blockHeight: this.latestBlock.height,
        toReset: true,
        toGetLowerBlocks: true
      })
    }
  }
}
</script>

<style scoped>

.explorer {
  height: calc(100vh - var(--header-height) - 1px - 2.25rem);
  padding-top: 2.25rem;
}

.explorer__chain {
  height: inherit;
}
.explorer__chain-main {
  height: calc(100% - 40px);
}
.explorer__chain-header {
  font-size: 3.1875rem;
  font-weight: var(--f-w-bold);
  margin-bottom: 1.5rem;
  padding-left: calc(var(--g-offset-side) - 4px);
}

.explorer__block {
  height: 100%;
}

.table-wrapper {
  --table-height: 86vh;
}

.empty-view {
  width: 100%;
  height: 100%;
  max-height: var(--table-height);
  height: var(--table-height);
  background-color: var(--c-bg-secondary);
  border-radius: var(--bd-radius-primary);
  display: flex;
  justify-content: center;
  align-items: center;
}

.table-wrapper {
  max-height: var(--table-height);
  height: var(--table-height);
}

.table-wrapper__utils {
  margin-left: 4px;
}
.table-wrapper__utils-btn {
  font-size: 0.875rem;
  color: var(--c-txt-grey);
}

@media only screen and (max-width: 1200px) {
  .table-wrapper {
    --table-height: 80vh;
  }
}


</style>