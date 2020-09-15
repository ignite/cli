<template>
  <div>

    <div class="temp" v-if="messages.length == 0">Waiting for blocks...</div>

    <div v-else class="table-wrapper">
      <TableWrapper 
        :tableHeads="['Height', 'Txs', 'Proposer', 'Block Hash', 'Age']"
        :tableId="tableId"        
        :containsInnerSheet="true"
      >
        <!-- <BlockSheet slot="innerSheet" :blockData=""/> -->
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
import { mapGetters, mapMutations } from 'vuex'
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
      tableGroupId: 'blocks-table',
      tendermintRootUrl: 'rpc.nylira.net',
      cosmosRootUrl: 'localhost:1317',
      messages: [],
      tableId: 'cosmosBlocksExplorer'
    }
  },
  computed: {
    ...mapGetters('cosmos', [ 'targetTable', 'isTableSheetActive' ]),
    ...mapGetters('cosmos/blocks', [ 'highlightedBlock' ]),
    fmtIsTableSheetActive() {
      return this.isTableSheetActive(this.tableId)
    },    
    fmtTargetTable() {
      return this.targetTable(this.tableId)
    },
    messagesForTable() {
      if (this.messages.length > 0) {
        return this.messages.map((message) => {
          const {
            time,
            height,
            proposer_address,
            num_txs
          } = message.header

          const {
            hash
          } = message.blockMeta.block_id

          return {
            blockMsg: {
              time_formatted: time.slice(0,5),
              time: time,
              height,
              proposer: proposer_address.slice(0,5),
              blockHash_sliced: `${hash.slice(0,30)}...`,
              blockHash: hash,
              txs: num_txs          
            },
            tableData: {
              id: height,
              isActive: false
            },
            txs: message.txsDecoded
          }          
        })        
      }
    }
  },  
  methods: {
    ...mapMutations('cosmos', [ 'setTableSheetState' ]),
    ...mapMutations('cosmos/blocks', [ 'setHighlightedBlock' ]),
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
    }
  },
  created() {
    let ws = new ReconnectingWebSocket(`wss://${this.tendermintRootUrl}:443/websocket`, [], { WebSocket: WebSocket });
    ws.onopen = function() {
      ws.send(
        JSON.stringify({
          jsonrpc: "2.0",
          method: "subscribe",
          id: "1",
          params: ["tm.event = 'NewBlock'"]
        })
      );
    };
    ws.onmessage = (msg) => {
      const { result } = JSON.parse(msg.data)

      // console.log(this.tendermintRootUrl)

      if (result.data && result.events) {
        const { data, events } = result        
        const { data: txsData, header } = data.value.block

        async function fetchBlockMeta() {
          try {
            return await axios.get(`https://rpc.nylira.net/block?${header.height}`)
          } catch (err) {
            console.error(err)
          }
        }
        async function fetchDecodedTx(txEncoded) {
          try {
            return await axios.post(`http://localhost:1317/txs/decode`, { tx: txEncoded }) 
          } catch (err) {
            console.error(txEncoded, err)
          }        
        }   
        
        /* TODO: Proposer address is in HEX format? Decoding API is required? */
        // async function fetchValidator() {
        //   try {
        //     console.log(header)
        //     return await axios.get(`https://lcd.nylira.net/staking/validators/${header.proposer_address}`)
        //   } catch (err) {
        //     console.error(err)
        //   }   
        // }
        // fetchValidator().then(validator => console.log(validator))

        const messageHolder = {
          header,
          txs: txsData.txs,
          blockMeta: null,
          txsDecoded: []
        }


        fetchBlockMeta()
          .then(blockMeta => {
            messageHolder.blockMeta = blockMeta.data.result.block_meta

            if (txsData.txs && txsData.txs.length > 0) {
              const txsDecoded = txsData.txs.map(txEncoded => fetchDecodedTx(txEncoded))
              
              txsDecoded.forEach(txRes => txRes.then(txResolved => {
                messageHolder.txsDecoded.push(txResolved.data.result)
              }))
            }    

            console.log(messageHolder)

            this.messages.unshift(messageHolder)                  
          })
   
      }         
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