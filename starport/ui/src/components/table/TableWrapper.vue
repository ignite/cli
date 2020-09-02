<template>
  <div class="table">

    <div class="table__row -row-group -header">
      <div
        v-for="(th, index) in tableHeads"
        :key="th+index"
        class="table__col"
      >{{th}}</div>
    </div>

    <slot/>

  </div>
</template>

<script>
import RowWrapper from './RowWrapper'
import Accordion from '@/components/accordion/Accordion'
import { accordionDataTemplate } from '@/components/accordion/Accordion'

export default {
  components: {
    RowWrapper,
    Accordion
  },
  props: {
    tableHeads: { type: Array, required: true },
    // tableRows: { type: Array, required: true },
  },  
  data() {
    return {
      exampleData2: [
        accordionDataTemplate({
          id: 1,
          active: false,
          title: "test",
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

.table >>> .-row-group {
  display: flex;
  padding: 0 8px;
}
.table >>> .table__row:not(:last-child) {
  margin-bottom: 1rem;
}
.table >>> .-row-group.-header {
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid black;
}

/* temporary table styling */
.table >>> .-row-group .table__col:nth-child(1) {
  width: 15%;
}
.table >>> .-row-group .table__col:nth-child(2) {
  width: 15%;
}
.table >>> .-row-group .table__col:nth-child(3) {
  width: 15%;
}
.table >>> .-row-group .table__col:nth-child(4) {
  flex-grow: 1;
}
.table >>> .-row-group .table__col:nth-child(5) {
  width: 5%;
}

/* temporary table styling */
.table >>> .-row-group.-panel .table__col:nth-child(1) {
  flex-grow: 1;  
}
.table >>> .-row-group.-panel .table__col:nth-child(2) {
  width: 15%;
}
.table >>> .-row-group.-panel .table__col:nth-child(3) {
  width: 20%;
}
.table >>> .-row-group.-panel .table__col:nth-child(4) {
  width: 15%;
}
.table >>> .-row-group.-panel .table__col:nth-child(5) {
  width: 5%;
}


</style>