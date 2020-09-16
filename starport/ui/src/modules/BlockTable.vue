<template>
  <div>

    <div class="temp" v-if="blockEntries.length == 0">Waiting for blocks...</div>

    <div v-else class="table-wrapper">
      <TableWrapper 
        :tableHeads="['Height', 'Txs', 'Proposer', 'Block Hash', 'Age']"
        :tableId="tableId"        
        :containsInnerSheet="true"
        @sheet-closed="handleSheetClose"
      >
        <BlockSheet slot="innerSheet" :blockData="highlightedBlock.data"/>
        <TableRowWrapper
          slot="tableContent"
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
  mixins: [blockHelpers],
  components: {
    TableWrapper,
    TableRowWrapper,    
    TableRowCellsGroup,
    BlockSheet
  },
  data() {
    return {
      tableId: 'cosmosBlocksExplorer'
    }
  },
  computed: {
    ...mapGetters('cosmos', [ 'targetTable', 'isTableSheetActive' ]),
    ...mapGetters('cosmos/blocks', [ 'highlightedBlock', 'blockEntries' ]),
    fmtIsTableSheetActive() {
      return this.isTableSheetActive(this.tableId)
    },    
    fmtTargetTable() {
      return this.targetTable(this.tableId)
    },
    messagesForTable() {
      return this.$_blockFormatter().blockForTable(this.blockEntries)
    }
  },  
  methods: {
    ...mapMutations('cosmos', [ 'setTableSheetState' ]),
    ...mapMutations('cosmos/blocks', [ 'setHighlightedBlock' ]),
    ...mapActions('cosmos/blocks', [ 'addBlockEntry' ]),
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
    }
  }
}
</script>

<style scoped>

.table-wrapper {
  max-height: 80vh;
  height: 80vh;
}

</style>