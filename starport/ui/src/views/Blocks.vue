<template>
  <div>

    <TableWrapper :tableHeads="['Time', 'Height', 'Proposer', 'Block Hash', 'Txs']">
      
      <Accordion
        :id="'accordion-example-1'"
        :content="exampleDataTwo"
      > 
        <TableRowWrapper
          v-for="item in exampleDataTwo"
          :key="item.id"                    
        >   
          <AccordionItem
            :itemData="item"
            :groupId="'accordion-example-1'"
          >
            <TableRowCellsGroup 
              slot="trigger" 
              :tableCells="['16s ago', '3111975', 'stake.fish', 'BBE79CF80CD...70A9F9F436A', '3']"
            />     
            <div slot="contents">
              <div class="table__panel panel">
                <div class="panel__top">
                  <p>Transactions</p>
                </div>
                <div class="panel__table">
                  <TableWrapper :tableHeads="['TxHash', 'Result', 'Fee', 'Gas', 'Msgs']">
                    <TableRowWrapper v-for="n in 5" :key="n">
                      <TableRowCellsGroup            
                        :tableCells="['16s ago', '3111975', 'stake.fish', 'BBE79CF80CD...70A9F9F436A', '3']"
                      />                            
                    </TableRowWrapper>
                  </TableWrapper>
                </div>
              </div>               
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

import Accordion from '@/components/accordion/Accordion'
import AccordionItem from '@/components/accordion/AccordionItem'
import { accordionDataTemplate } from '@/components/accordion/Accordion'

export default {
  components: {
    TableWrapper,
    TableRowWrapper,    
    TableRowCellsGroup,
    Accordion,
    AccordionItem
  },
  data() {
    return {
      exampleData: [[
        accordionDataTemplate({
          id: 1,
          active: false,
        })
      ],[
        accordionDataTemplate({
          id: 1,
          active: false,
        })
      ]],
      exampleDataTwo: [
        { id: 1, isActive: false },
        { id: 2, isActive: false }
      ]
    }
  }
}
</script>

<style scoped>

.table .-row-group {
  display: flex;
  padding: 0 8px;
}
.table .table__row:not(:last-child) {
  margin-bottom: 1rem;
}
.table .-row-group.-header {
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid black;
}

/* temporary table styling */
.table .-row-group .table__col:nth-child(1) {
  width: 15%;
}
.table .-row-group .table__col:nth-child(2) {
  width: 15%;
}
.table .-row-group .table__col:nth-child(3) {
  width: 15%;
}
.table .-row-group .table__col:nth-child(4) {
  flex-grow: 1;
}
.table .-row-group .table__col:nth-child(5) {
  width: 5%;
}

/* temporary table styling */
.table .-row-group.-panel .table__col:nth-child(1) {
  flex-grow: 1;  
}
.table .-row-group.-panel .table__col:nth-child(2) {
  width: 15%;
}
.table .-row-group.-panel .table__col:nth-child(3) {
  width: 20%;
}
.table .-row-group.-panel .table__col:nth-child(4) {
  width: 15%;
}
.table .-row-group.-panel .table__col:nth-child(5) {
  width: 5%;
}

.table__panel {

}

.panel {
  background-color: rgba(0,0,0, 0.1);
}

</style>