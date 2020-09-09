<template>
  <div 
    :class="['table', sheetStore.isActive ? '-is-collapsed' : '']"
    ref="table"
  >

    <div 
      v-if="containsInnerSheet"
      :class="['table__sheet', sheetStore.isActive ? '-is-active' : '']"
    >
      <TableSheet :blockData="rowStore.activeRowData" />
    </div>    

    <div 
      :class="['table__head']"
    >
      <RowCells
        :isTableHead="true"
        :tableCells="tableHeads"
      />
    </div>

    <div :class="['table__rows-wrapper']">
      <div>
        <slot/>
      </div>
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
      sheetStore: {
        isActive: false,
      },
      rowStore: {
        activeRowId: null,
        activeRowData: null
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
  height: inherit;
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
  /* position: sticky; */
  top: 0;
  right: 0;
  width: 70%;
  /* width: calc(100% - (var(--tc-w-1) + var(--tc-w-2) + var(--tc-w-3))); */
  /* width: 100%; */
  height: 100%;  
  /* height: 80vh; */
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

.table__rows-wrapper {
  min-height: inherit;
  max-height: inherit;
  overflow-y: scroll;
}

/* .table__sheet >>> .sheet {
  width: calc(100% - (var(--tc-w-1) + var(--tc-w-2) + var(--tc-w-3)));
  margin-left: auto;
} */
.table.-is-collapsed {
  --tc-w-1: 40%;
  --tc-w-2: 20%;
  --tc-w-3: 40%;
}
.table.-is-collapsed .table__rows-wrapper,
.table.-is-collapsed .table__head {
  width: 30%;
  overflow-x: hidden;
}
.table.-is-collapsed .table__rows-wrapper >>> .table__row .table__cells .table__col:nth-last-child(-n+2),
.table.-is-collapsed .table__head >>> .table__cells .table__col:nth-last-child(-n+2) {
  display: none;
}


</style>