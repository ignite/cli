<template>
  <div class="table" ref="table">

    <RowCells
      :isTableHead="true"
      :tableCells="tableHeads"
    />

    <slot/>
    
    <div 
      v-if="containsInnerSheet"
      :class="['table__sheet', sheetState.isActive ? '-is-active' : '']"
    >
      <TableSheet />
    </div>

  </div>
</template>

<script>
import RowWrapper from './RowWrapper'
import RowCells from './RowCellsGroup'
import TableSheet from './InnerSheet'
import Accordion from '@/components/accordion/Accordion'

export default {
  components: {
    RowWrapper,
    RowCells,
    TableSheet,
    Accordion
  },
  props: {
    tableHeads: { type: Array, required: true },
    containsInnerSheet: { type: Boolean, default: false }
  },
  data() {
    return {
      sheetState: {
        isActive: false
      },
      rowState: {
        activeRowId: null
      }
    }
  }
}
</script>

<style scoped>

.table {
  --tc-w-1: 10%;
  --tc-w-2: 5%;
  --tc-w-3: 15%;
  --tc-w-4: 1;
  --tc-w-5: 10%;
}
.table {
  position: relative;
}

.table >>> .table__cells {
  padding-top: 0.8rem;
  padding-bottom: 0.8rem;
}
.table >>> .table__cells:first-child {
  padding-top: 0.8rem;
}
.table >>> .table__cells {
  padding-left: 1rem;
  padding-right: 1rem;
}

/* temporary table styling */
.table >>> .table__cells .table__col:nth-child(1) {
  width: var(--tc-w-1);
}
.table >>> .table__cells .table__col:nth-child(2) {
  width: var(--tc-w-2);
}
.table >>> .table__cells .table__col:nth-child(3) {
  width: var(--tc-w-3);
}
.table >>> .table__cells .table__col:nth-child(4) {
  flex-grow: var(--tc-w-4);
}
.table >>> .table__cells .table__col:nth-child(5) {
  width: var(--tc-w-5);
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


.table__sheet {
  position: absolute;
  top: 0;
  right: 0;
  width: calc(100% - (var(--tc-w-1) + var(--tc-w-2) + var(--tc-w-3)));
  height: 100%;  
  background-color: var(--c-bg-primary);

  transform: translate3d(100%, 0, 0);
  transition: transform ease-out .25s;
  will-change: transform;
}
.table__sheet.-is-active {  
  transform: translate3d(0%,0,0);
  transition: transform ease-out .25s;
  will-change: transform;
}

</style>