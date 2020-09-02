<template>
  <div>

    <table-wrapper :tableHeads="['Time', 'Height', 'Proposer', 'Block Hash', 'Txs']">
      <table-row-wrapper>
        <Accordion
          id="accordion-example-1"
          :content="exampleData2"
        > 
          <TableRowCellsGroup slot="trigger" :tableCells="['16s ago', '3111975', 'stake.fish', 'BBE79CF80CD...70A9F9F436A', '3']" />        
        </Accordion>
      </table-row-wrapper>
    </table-wrapper>

    <!-- <div class="table">
      <div class="table__row -row-group -header">
        <div class="table__col">Time</div>
        <div class="table__col">Height</div>
        <div class="table__col">Proposer</div>
        <div class="table__col">Block Hash</div>
        <div class="table__col">Txs</div>
      </div>
      <div 
        class="table__row"
        v-for="n in 5"
        :key="n"
      >
        <div class="-row-group">
          <div class="table__col">16s ago</div>
          <div class="table__col">3111975</div>
          <div class="table__col">stake.fish</div>
          <div class="table__col">BBE79CF80CD...70A9F9F436A</div>
          <div class="table__col">3</div>
        </div>
      </div>

      <Accordion id="accordion-example-1" :content="exampleData2" /> 
      
      <div class="table__row">
        <div class="-row-group">
          <div class="table__col">16s ago</div>
          <div class="table__col">3111975</div>
          <div class="table__col">stake.fish</div>
          <div class="table__col">BBE79CF80CD...70A9F9F436A</div>
          <div class="table__col">3</div>
        </div>

        <div class="table__panel panel">
          <div class="panel__top">
            <p>Transactions</p>
          </div>
          <div class="panel__table">
            <div class="table">
              <div class="table__row -row-group -panel -header">
                <div class="table__col">TxHash</div>
                <div class="table__col">Result</div>
                <div class="table__col">Fee</div>
                <div class="table__col">Gas</div>
                <div class="table__col">Msgs</div>
              </div>        
              <div 
                class="table__row"
                v-for="n in 5"
                :key="n"
              >
                <div class="-row-group -panel">
                  <div class="table__col">BBEXTDADF...25SDFS...</div>
                  <div class="table__col">Success</div>
                  <div class="table__col">0.000001 ATOM</div>
                  <div class="table__col">106,019 / 200,000</div>
                  <div class="table__col">2</div>
                </div>
              </div>    
            </div>
          </div>
        </div>
      </div>

    </div> -->

  </div>
</template>

<script>
import TableWrapper from '@/components/table/TableWrapper'
import TableRowWrapper from '@/components/table/RowWrapper'
import TableRowCellsGroup from '@/components/table/RowCellsGroup'

import Accordion from '@/components/accordion/Accordion'
import { accordionDataTemplate } from '@/components/accordion/Accordion'

export default {
  components: {
    TableWrapper,
    TableRowWrapper,    
    TableRowCellsGroup,
    Accordion
  },
  data() {
    return {
      exampleData2: [
        accordionDataTemplate({
          id: 1,
          active: false,
          title: `<TableRowCellsGroup :tableCells="['16s ago', '3111975', 'stake.fish', 'BBE79CF80CD...70A9F9F436A', '3']" />`,
          content: `
            <div class="table__panel panel">
              <div class="panel__top">
                <p>Transactions</p>
              </div>
              <div class="panel__table">
                <div class="table">
                  <div class="table__row -row-group -panel -header">
                    <div class="table__col">TxHash</div>
                    <div class="table__col">Result</div>
                    <div class="table__col">Fee</div>
                    <div class="table__col">Gas</div>
                    <div class="table__col">Msgs</div>
                  </div>        
                  <div 
                    class="table__row"
                    v-for="n in 5"
                    :key="n"
                  >
                    <div class="-row-group -panel">
                      <div class="table__col">BBEXTDADF...25SDFS...</div>
                      <div class="table__col">Success</div>
                      <div class="table__col">0.000001 ATOM</div>
                      <div class="table__col">106,019 / 200,000</div>
                      <div class="table__col">2</div>
                    </div>
                  </div>    
                </div>
              </div>
            </div>          
          `
        })
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