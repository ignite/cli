<template>
  <div>

    <TableWrapper :tableHeads="['Time', 'Height', 'Proposer', 'Block Hash', 'Txs']">
      
      <Accordion :id="tableGroupId"> 
        <TableRowWrapper
          v-for="item in exampleDataTwo"
          :key="item.id"                    
        >   
          <AccordionItem
            :itemData="item"
            :groupId="tableGroupId"
          >
            <TableRowCellsGroup 
              slot="trigger" 
              :tableCells="['16s ago', '3111975', 'stake.fish', 'BBE79CF80CD...70A9F9F436A', '3']"
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
      exampleDataTwo: [
        { id: 1, isActive: false },
        { id: 2, isActive: false }
      ]
    }
  },
  // created() {
  //   let ws = new ReconnectingWebSocket("wss://rpc.nylira.net:443/websocket", [], { WebSocket: WebSocket });
  //   ws.onopen = function() {
  //     ws.send(
  //       JSON.stringify({
  //         jsonrpc: "2.0",
  //         method: "subscribe",
  //         id: "1",
  //         params: ["tm.event = 'NewBlock'"]
  //       })
  //     );
  //   };
  //   ws.onmessage = function(msg) {
  //     let msgData = JSON.parse(msg.data);
  //     console.log(msgData);
  //   };    
  // }
}
</script>

<style scoped>

</style>