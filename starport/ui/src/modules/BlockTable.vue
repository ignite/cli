<template>
  <div>

    <div class="table-wrapper">
      <TableWrapper 
        :tableHeads="setRowCellsProps(['Height', 'Txs', 'Block Hash', 'Age'])"
        :tableId="blocksExplorerTableId"        
        :containsInnerSheet="true"
        :isTableEmpty="isBlocksTableEmpty"
        :isTableLoading="states.isLoading"
        :tableEmptyMsg="blockTableEmptyText"
        @sheet-closed="handleSheetClose"
        @scrolled-top="handleScrollTop"
        @scrolled-bottom="handleScrollBottom"
      >
        <div slot="utils" class="table-wrapper__utils">
          <button 
            class="table-wrapper__utils-btn"
            @click="handleFilterClick"
          >{{blockFilterText}}</button> 
        </div>

        <BlockSheet slot="innerSheet" :blockData="localHighlightedBlock"/>

        <div slot="tableContent" v-if="fmtBlockData">
          <TableRowWrapper 
            v-for="block in fmtBlockData"
            :key="block.blockMsg.blockHash"  
            :rowData="block"
            :rowId="block.blockMsg.blockHash"   
            :isRowActive="block.blockMsg.blockHash === localHighlightedBlock.id"   
            :isWithInnerSheet="true" 
            @row-clicked="handleRowClick"
          >   
            <TableRowCellsGroup 
              :tableCells="setRowCellsProps([
                block.blockMsg.height,
                block.blockMsg.txs,
                block.blockMsg.blockHash_sliced,
                block.blockMsg.time_formatted,
              ])"
            />     
          </TableRowWrapper>  
        </div>   
        
      </TableWrapper>
    </div>

  </div>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import blockHelpers from '@/mixins/blocks/helpers'

import axios from "axios"
import ReconnectingWebSocket from "reconnecting-websocket"
import _ from 'lodash'
import { uuid } from 'vue-uuid'

import TableWrapper from '@/components/table/TableWrapper'
import TableRowWrapper from '@/components/table/RowWrapper'
import TableRowCellsGroup from '@/components/table/RowCellsGroup'
import BlockSheet from '@/modules/BlockSheet'

export default {
  components: {
    TableWrapper,
    TableRowWrapper,    
    TableRowCellsGroup,
    BlockSheet,
  },
  data() {
    return {
      blockFormatter: blockHelpers.blockFormatter(),
      states: {
        isHidingBlocksWithoutTxs: false,
        isScrolledInTopHalf: true,
        isLoading: false
      },
      localHighlightedBlock: null
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
    setRowCellsProps(cells) {
      return cells.map(cell => ({ content: cell, id: uuid.v1() }))
    },
    handleRowClick(rowId, rowData) {
      const setTableRowStore = (isToActive=false, payload) => {
        const highlightBlockPayload = isToActive ? {
          id: payload.rowId,
          data: payload.rowData
        } : null
        
        this.setHighlightedBlock({
          block: highlightBlockPayload
        })
      }

      const isActiveRowClicked = this.highlightedBlock.id === rowId
      
      if (this.fmtIsTableSheetActive) {
        if (isActiveRowClicked) {
          this.setTableSheetState({
            tableId: this.blocksExplorerTableId,
            sheetState: false
          })
          setTableRowStore()
        } else {
          setTableRowStore(true, { rowId: rowId, rowData: rowData })
        }
      } else {
        this.setTableSheetState({
          tableId: this.blocksExplorerTableId,
          sheetState: true
        })
        setTableRowStore(true, { rowId: rowId, rowData: rowData })
      }
    },
    handleSheetClose() {
      this.setHighlightedBlock({ block: null })
    },
    handleFilterClick() {
      this.setHighlightedBlock({ block: null })
      this.setTableSheetState({
        tableId: this.blocksExplorerTableId,
        sheetState: false
      })      
      this.states.isHidingBlocksWithoutTxs = !this.states.isHidingBlocksWithoutTxs
    },
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
  watch: {
    highlightedBlock() {
      this.localHighlightedBlock = this.highlightedBlock
    },
  },
  created() {
    this.localHighlightedBlock = this.highlightedBlock
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