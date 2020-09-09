<template>
  <div class="table">

    <RowCells
      :isTableHead="true"
      :tableCells="tableHeads"
    />

    <slot/>

  </div>
</template>

<script>
import RowWrapper from './RowWrapper'
import RowCells from './RowCellsGroup'
import Accordion from '@/components/accordion/Accordion'

export default {
  components: {
    RowWrapper,
    RowCells,
    Accordion
  },
  props: {
    tableHeads: { type: Array, required: true },
  }
}
</script>

<style scoped>

.table >>> .table__cells {
  padding-top: 0.8rem;
  padding-bottom: 0.8rem;
}
.table >>> .table__cells{
  padding-left: 1rem;
  padding-right: 1rem;
}

/* temporary table styling */
.table >>> .table__cells .table__col:nth-child(1) {
  width: 15%;
}
.table >>> .table__cells .table__col:nth-child(2) {
  width: 15%;
}
.table >>> .table__cells .table__col:nth-child(3) {
  width: 15%;
}
.table >>> .table__cells .table__col:nth-child(4) {
  flex-grow: 1;
}
.table >>> .table__cells .table__col:nth-child(5) {
  width: 5%;
}

/* temporary table styling */
.table >>> .table__cells.-panel .table__col:nth-child(1) {
  flex-grow: 1;  
}
.table >>> .table__cells.-panel .table__col:nth-child(2) {
  width: 15%;
}
.table >>> .table__cells.-panel .table__col:nth-child(3) {
  width: 20%;
}
.table >>> .table__cells.-panel .table__col:nth-child(4) {
  width: 15%;
}
.table >>> .table__cells.-panel .table__col:nth-child(5) {
  width: 5%;
}


</style>