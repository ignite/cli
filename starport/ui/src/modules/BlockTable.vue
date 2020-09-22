<template>
  <div>

    <div class="table-wrapper">
      <TableWrapper 
        :tableHeads="['Height', 'Txs', 'Proposer', 'Block Hash', 'Age']"
        :tableId="tableId"        
        :containsInnerSheet="true"
        :isTableEmpty="isBlocksTableEmpty"
        :tableEmptyMsg="blockTableEmptyText"
        @sheet-closed="handleSheetClose"
      >
        <div slot="utils" class="table-wrapper__utils">
          <!-- TODO: enhance UI -->
          <button 
            class="table-wrapper__utils-btn"
            @click="handleFilterClick"
          >{{blockFilterText}}</button> 
        </div>

        <BlockSheet slot="innerSheet" :blockData="highlightedBlock.data"/>

        <div slot="tableContent">
          <TableRowWrapper 
            v-for="msg in messagesForTable"
            :key="msg.tableData.id"  
            :rowData="msg"
            :rowId="msg.blockMsg.blockHash"   
            :isRowActive="msg.blockMsg.blockHash === highlightedBlock.id"   
            :isWithInnerSheet="true" 
            @row-clicked="handleRowClick"
          >   
            <TableRowCellsGroup 
              :tableCells="[
                msg.blockMsg.height,
                msg.blockMsg.txs,
                msg.blockMsg.proposer,
                msg.blockMsg.blockHash_sliced,
                msg.blockMsg.time_formatted,
              ]"
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

import TableWrapper from '@/components/table/TableWrapper'
import TableRowWrapper from '@/components/table/RowWrapper'
import TableRowCellsGroup from '@/components/table/RowCellsGroup'
import BlockSheet from '@/modules/BlockSheet'

export default {
  components: {
    TableWrapper,
    TableRowWrapper,    
    TableRowCellsGroup,
    BlockSheet
  },
  data() {
    return {
      tableId: 'cosmosBlocksExplorer',
      blockFormatter: blockHelpers.blockFormatter(),
      states: {
        isHidingBlocksWithoutTxs: false
      }
    }
  },
  computed: {
    /*
     *
     * Vuex 
     *
     */
    ...mapGetters('cosmos/ui', [ 'targetTable', 'isTableSheetActive' ]),
    ...mapGetters('cosmos/blocks', [ 'highlightedBlock', 'blockEntries' ]),
    /*
     *
     * Local
     * 
     */    
    fmtIsTableSheetActive() {
      return this.isTableSheetActive(this.tableId)
    },    
    fmtTargetTable() {
      return this.targetTable(this.tableId)
    },
    messagesForTable() {
      const fmtBlockForTable = this.blockFormatter.blockForTable(this.blockEntries)
      if (this.states.isHidingBlocksWithoutTxs) {
        return this.blockFormatter.filterBlock(fmtBlockForTable).hideBlocksWithoutTxs()
      }
      return fmtBlockForTable
    },
    isBlocksTableEmpty() {
      return this.blockEntries.length<=0 || !this.messagesForTable || this.messagesForTable?.length<=0
    },
    blockFilterText() {
      return !this.states.isHidingBlocksWithoutTxs
        ? 'Hide blocks without txs'
        : 'Show blocks without txs'
    },
    blockTableEmptyText() {
      return (this.blockEntries.length>=0 && this.messagesForTable?.length<=0 && this.states.isHidingBlocksWithoutTxs) 
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
    ...mapMutations('cosmos/blocks', [ 'setHighlightedBlock' ]),
    ...mapActions('cosmos/blocks', [ 'addBlockEntry' ]),
    /*
     *
     * Local 
     *
     */      
    handleRowClick(rowId, rowData) {
      const setTableRowStore = (isToActive=false, payload) => {
        const highlightBlockPayload = isToActive ? {
          id: payload.rowId,
          data: payload.rowData
        } : null
        
        this.setHighlightedBlock(highlightBlockPayload)
      }

      const isActiveRowClicked = this.highlightedBlock.id === rowId
      
      if (this.fmtIsTableSheetActive) {
        if (isActiveRowClicked) {
          this.setTableSheetState({
            tableId: this.tableId,
            sheetState: false
          })
          setTableRowStore()
        } else {
          setTableRowStore(true, { rowId: rowId, rowData: rowData })
        }
      } else {
        this.setTableSheetState({
          tableId: this.tableId,
          sheetState: true
        })
        setTableRowStore(true, { rowId: rowId, rowData: rowData })
      }
    },
    handleSheetClose() {
      this.setHighlightedBlock(null)
    },
    handleFilterClick() {
      this.setHighlightedBlock(null)
      this.setTableSheetState({
        tableId: this.tableId,
        sheetState: false
      })      
      this.states.isHidingBlocksWithoutTxs = !this.states.isHidingBlocksWithoutTxs
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
  max-height: 86vh;
  height: 86vh;
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