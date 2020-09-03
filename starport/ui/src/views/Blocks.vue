<template>
  <div>

    <div class="temp" v-if="messages.length == 0">Waiting for blocks...</div>

    <TableWrapper v-else :tableHeads="['Time', 'Height', 'Proposer', 'Block Hash', 'Txs']">
      
      <Accordion :id="tableGroupId"> 
        <TableRowWrapper
          v-for="msg in messagesForTable"
          :key="msg.tableData.id"                    
        >   
          <AccordionItem
            :itemData="msg.tableData"
            :groupId="tableGroupId"
          >
            <TableRowCellsGroup 
              slot="trigger" 
              :tableCells="[
                msg.blockMsg.time,
                msg.blockMsg.height,
                msg.blockMsg.proposer,
                msg.blockMsg.blockHash,
                msg.blockMsg.txs,
              ]"
            />     
            <div slot="contents">
              <InnerTable :parentGroupId="tableGroupId" />
            </div>
          </AccordionItem>     
        </TableRowWrapper>   
                
      </Accordion>

    </TableWrapper>

  </div>
</template>

<script>
import TableWrapper from '@/components/table/TableWrapper'
import TableRowWrapper from '@/components/table/RowWrapper'
import TableRowCellsGroup from '@/components/table/RowCellsGroup'
import InnerTable from '@/components/table/InnerTable'

import Accordion from '@/components/accordion/Accordion'
import AccordionItem from '@/components/accordion/AccordionItem'

import ReconnectingWebSocket from "reconnecting-websocket";

export default {
  components: {
    TableWrapper,
    TableRowWrapper,    
    TableRowCellsGroup,
    InnerTable,
    Accordion,
    AccordionItem
  },
  data() {
    return {
      tableGroupId: 'blocks-table',
      messages: [],
      exampleDataTwo: [
        { id: 1, isActive: false },
        { id: 2, isActive: false }
      ]
    }
  },
  computed: {
    messagesForTable() {
      if (this.messages.length > 0) {
        return this.messages.map((message) => {
          const {
            time,
            height,
            proposer_address,
            app_hash,
            num_txs
          } = message.header

          return {
            blockMsg: {
              time: time.slice(0,5),
              height,
              proposer: proposer_address.slice(0,5),
              blockHash: app_hash.slice(0,10),
              txs: num_txs          
            },
            tableData: {
              id: height,
              isActive: false
            }
          }          
        })        
      }
    }
  },  
  created() {
    let ws = new ReconnectingWebSocket("wss://rpc.nylira.net:443/websocket", [], { WebSocket: WebSocket });
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

      if (result.data && result.events) {
        const { data, events } = result        
        const { data: txsData, header } = data.value.block
        
        this.messages.push({
          header,          
          txs: txsData.txs
        })
      }

      // console.log(this.messages)
    }
  }
}
</script>

<style scoped>

</style>