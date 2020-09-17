<template>
  <div 
    :class="['table', fmtIsTableSheetActive ? '-is-collapsed' : '']"
    ref="table"
  >
    <div class="table__utils">
      <button v-if="containsInnerSheet" 
        @click="handleSheetClose"
        class="table__utils-sheet-btn"
      ></button>
    </div>

    <div class="table__wrapper">
      <div v-if="containsInnerSheet"
        :class="['table__sheet', fmtIsTableSheetActive ? '-is-active' : '']"
      >
        <slot name="innerSheet"/>
      </div>    

      <div class="table__main">
        <div :class="['table__head']">
          <RowCells
            :isTableHead="true"
            :tableCells="tableHeads"
          />
        </div>
        <div :class="['table__rows-wrapper']">
          <div><slot name="tableContent"/></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapMutations } from 'vuex'

import RowWrapper from './RowWrapper'
import RowCells from './RowCellsGroup'

export default {
  components: {
    RowWrapper,
    RowCells
  },
  props: {
    tableHeads: { type: Array, required: true },
    tableId: { type: String, required: true },
    containsInnerSheet: { type: Boolean, default: false }
  },
  computed: {
    ...mapGetters('cosmos', [ 'targetTable', 'isTableSheetActive' ]),
    fmtIsTableSheetActive() {
      return this.isTableSheetActive(this.tableId)
    }
  },
  methods: {
    ...mapMutations('cosmos', [
      'createTable',
      'setTableSheetState',
    ]),    
    handleSheetClose() {
      this.setTableSheetState({
        tableId: this.tableId,
        sheetState: false
      })

      this.$emit('sheet-closed')
    }
  },
  created() {
    if (!this.targetTable(this.tableId)) {
      this.createTable(this.tableId)
    }
  }
}
</script>

<style scoped>

.table {
  --tc-w-1: 10%;
  --tc-w-2: 5%;
  --tc-w-3: 20%;
  --tc-w-4: 1;
  --tc-w-5: 15%;
}
.table {
  position: relative;
  height: inherit;
  padding-bottom: 1.5rem;
}

.table__utils {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
.table__utils-sheet-btn {
  position: relative;
  padding: 0.5rem 0.25rem 1rem 0.5rem;
  opacity: 0;
}
.table__utils-sheet-btn:before,
.table__utils-sheet-btn:after {
  content: '';
  position: absolute;
  bottom: 1.125rem;
  right: 0;
  width: 1rem;
  height: 2px;
  background-color: var(--c-txt-grey);
}
.table__utils-sheet-btn:before {
  transform: rotateZ(45deg);
}
.table__utils-sheet-btn:after {
  transform: rotateZ(-45deg);
}

.table__wrapper {
  position: relative;
  height: inherit;
  /* border: 1px solid var(--c-theme-secondary);     */
  overflow: hidden;
  /* padding-left: 1rem; */
  /* padding-right: 1rem;   */
}
.table__wrapper .table__rows-wrapper {
  padding-right: 1rem;  
}
.table__wrapper >>> .table__cells.-header {
  padding-right: 2rem;  
}

.table__main {
  padding-left: 1rem;
  box-sizing: border-box;
  background-color: var(--c-bg-third);
  border-radius: 8px;  
  height: inherit;
  max-height: inherit;
}

.table >>> .table__cells {
  padding-top: 0.8rem;
  padding-bottom: 0.8rem;
}
.table >>> .table__cells {
  padding-left: 0.8rem;
  padding-right: 0.8rem;
}

/* temporary table styling */
.table >>> .table__cells .table__col:nth-child(1) {
  min-width: var(--tc-w-1);
  max-width: var(--tc-w-1);
}
.table >>> .table__cells .table__col:nth-child(2) {
  min-width: var(--tc-w-2);
  max-width: var(--tc-w-2);
}
.table >>> .table__cells .table__col:nth-child(3) {
  min-width: var(--tc-w-3);
  max-width: var(--tc-w-3);
}
.table >>> .table__cells .table__col:nth-child(4) {
  flex-grow: var(--tc-w-4);
}
.table >>> .table__cells .table__col:nth-child(5) {
  min-width: var(--tc-w-5);
  max-width: var(--tc-w-5);
}


.table__sheet {
  position: absolute;
  top: 0;
  right: 0;
  width: calc(100% - 23vw - 1rem);
  height: 100%;  
}
.table__sheet {  
  transform: translate3d(100%, 0, 0);
  opacity: 0;
  transition: transform ease-out .25s;  
  will-change: transform;
}
.table__sheet.-is-active {  
  transform: translate3d(0,0,0);
  opacity: 1;  
  transition: all ease-out .3s;
  will-change: transform;
}

.table__rows-wrapper {
  height: calc(100% - 1rem);
  min-height: inherit;
  max-height: inherit;
  overflow-y: scroll;
  padding-right: 1rem;
}
.table__rows-wrapper::-webkit-scrollbar { /* width */
  width: 6px;
}
.table__rows-wrapper::-webkit-scrollbar-track { /* Track */
  /* box-shadow: inset 0 0 1px var(--c-bg-grey);  */
  background: var(--c-bg-third); 
}
.table__rows-wrapper::-webkit-scrollbar-thumb { /* Handle */
  background-color: var(--c-bg-secondary); 
  border-radius: 10px;
}
.table__rows-wrapper::-webkit-scrollbar-thumb:hover { /* Handle on hover */
  background: var(--c-contrast-secondary); 
}


.table.-is-collapsed {
  --tc-w-1: 40%;
  --tc-w-2: 20%;
  --tc-w-3: 40%;
}
.table .table__main {
  width: 100%;
  transition: width .3s ease-in-out;
  will-change: width;
}
.table.-is-collapsed .table__main {
  width: 23vw;
  transition: width .3s ease-in-out;
  will-change: width;
}
/* .table.-is-collapsed .table__rows-wrapper,
.table.-is-collapsed .table__head {
  width: 30%;
  overflow-x: hidden;
} */
.table .table__rows-wrapper >>> .table__row .table__cells .table__col:nth-last-child(-n+2),
.table .table__head >>> .table__cells .table__col:nth-last-child(-n+2) {
  white-space: nowrap;
}
.table.-is-collapsed .table__rows-wrapper >>> .table__row .table__cells .table__col:nth-last-child(-n+2),
.table.-is-collapsed .table__head >>> .table__cells .table__col:nth-last-child(-n+2) {
  /* display: none; */
  opacity: 0;
  white-space: nowrap;
}
.table.-is-collapsed .table__utils button {
  opacity: 1;
  transition: opacity .3s;
}


</style>