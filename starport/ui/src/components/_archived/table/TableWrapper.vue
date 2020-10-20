<template>
  <div 
    :class="['table', {'-is-collapsed': fmtIsTableSheetActive, '-is-loading': isTableLoading}]"
    ref="table"
  >
    <div class="table__utils">
      <div class="table__utils-wrapper">
        <slot name="utils" />
      </div>
      <div class="table__utils-wrapper">
        <button v-if="containsInnerSheet" 
          @click="handleSheetClose"
          class="table__utils-sheet-btn"
        ></button>
      </div>      
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
            :cellWidths="colWidths"
          />
        </div>
        <div 
          :class="['table__rows-wrapper']"
          @scroll="[handleTableScroll($event), updateScrollValue()]"
        >
          <div v-if="!isTableEmpty"><slot name="tableContent"/></div>
          <div v-else class="table__rows-wrapper-empty-view"><p>{{tableEmptyMsg}}</p></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapMutations } from 'vuex'

import _ from 'lodash'

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
    containsInnerSheet: { type: Boolean, default: false },
    isTableEmpty: { type: Boolean, default: true },
    isTableLoading: { type: Boolean, default: false },
    tableEmptyMsg: { type: String, default: 'Waiting for blocks' },
    colWidths: {
      type: Array,
      validator(value) {
        return value.filter(val => typeof val === 'string').length === value.length
      }
    }
  },
  data() {
    return {
      lastScrolledHeight: 0,
      lastScrolledTop: 0,
      lastTimestamp: undefined
    }
  },
  computed: {
    ...mapGetters('cosmos/ui', [ 'targetTable', 'isTableSheetActive' ]),
    fmtIsTableSheetActive() {
      return this.isTableSheetActive(this.tableId)
    }
  },
  methods: {
    ...mapMutations('cosmos/ui', [
      'createTable',
      'setTableSheetState',
    ]),    
    handleSheetClose() {
      this.setTableSheetState({
        tableId: this.tableId,
        sheetState: false
      })

      this.$emit('sheet-closed')
    },
    handleTableScroll: _.debounce(function(event) {
      const $table = event.target
      const scrolledHeight = $table.scrollTop + $table.offsetHeight
      const tableScrollHeight = $table.scrollHeight

      const isScrolledToTop = scrolledHeight <= $table.offsetHeight
      const isScrolledToBottom = scrolledHeight + 100 >= tableScrollHeight
      const isOnTopHalf = $table.scrollTop < (tableScrollHeight-$table.offsetHeight) / 2

      const isCallableScrolledDistance = 
        $table.offsetHeight / Math.abs(scrolledHeight-this.lastScrolledHeight) > 25
      
      if (isCallableScrolledDistance) {
        if (isScrolledToBottom) this.$emit('scrolled-bottom')
        if (isScrolledToTop) {
          this.$emit('scrolled-top')
          /*
           *
           // Scroll down the table a bit to prevent staying on top
           *
           */
          $table.scrollBy({
            top: 5,
            left: 0,
            behavior: 'smooth'
          })
        }
      }
    }, 250),
    updateScrollValue() {
      const $table = event.target
      const scrolledHeight = $table.scrollTop + $table.offsetHeight 
      this.lastScrolledHeight = scrolledHeight     
      this.lastScrolledTop = $table.scrollTop     
    }
  },
  created() {
    if (!this.targetTable(this.tableId)) {
      this.createTable(this.tableId)
    }
  
    function checkScrollPosition(timestamp) {
      if (this.lastTimestamp === undefined) {
        this.lastTimestamp = timestamp
      }
      if (timestamp - this.lastTimestamp > 500) {
        if (this.lastScrolledTop === 0) this.$emit('scrolled-top')
        this.lastTimestamp = timestamp            
      }
      window.requestAnimationFrame(checkScrollPosition.bind(this))
    }
    // window.requestAnimationFrame(checkScrollPosition.bind(this))    
  }
}
</script>

<style scoped>

.table {
  --table-collapsed-width: 15vw;
}
.table {
  position: relative;
  height: inherit;
  padding-bottom: 1.5rem;
}

.table__utils {
  width: 100%;
  display: flex;
  justify-content: space-between;
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

.table__wrapper {
  position: relative;
  transition: opacity .3s ease-in-out;
}
.table__wrapper:before,
.table__wrapper:after {
  position: absolute;
  z-index: 1;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  opacity: 0;
  transition: all .3s ease-in-out;  
}
.table__wrapper:before {
  content: '';
  background-color: var(--c-bg-secondary);
}
.table__wrapper:after {
  content: 'Fetching blocks';
  display: flex;  
  justify-content: center;
  align-items: center;
  color: var(--c-txt-grey);
  animation: tempLoadingEffect 2s ease-in-out infinite;  
}
.table.-is-loading .table__wrapper:before {
  opacity: .8;
}
.table.-is-loading .table__wrapper:after {
  opacity: 1;
}
.table.-is-loading .table__wrapper:before,
.table.-is-loading .table__wrapper:after {
  pointer-events: initial;
  transition: all .3s ease-in-out;
}

.table__main {
  padding-left: 1rem;
  box-sizing: border-box;
  background-color: var(--c-bg-secondary);
  border-radius: var(--bd-radius-primary);  
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


.table__sheet {
  position: absolute;
  top: 0;
  right: 0;
  width: calc(100% - var(--table-collapsed-width) - 1rem);
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
  transition: transform ease-out .3s;
  will-change: transform;
}

.table__rows-wrapper {
  height: calc(100% - 5rem);
  min-height: inherit;
  max-height: inherit;
  overflow-y: scroll;
  overflow-x: hidden;
  padding-right: 1rem;
}
.table__rows-wrapper::-webkit-scrollbar { /* width */
  width: 6px;
}
.table__rows-wrapper::-webkit-scrollbar-track { /* Track */
  /* box-shadow: inset 0 0 1px var(--c-bg-grey);  */
  background: var(--c-bg-secondary); 
}
.table__rows-wrapper::-webkit-scrollbar-thumb { /* Handle */
  background-color: var(--c-bg-third); 
  border-radius: 10px;
}
.table__rows-wrapper::-webkit-scrollbar-thumb:hover { /* Handle on hover */
  background: var(--c-contrast-secondary); 
}

.table__rows-wrapper-empty-view {
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--c-txt-grey);
  animation: tempLoadingEffect 1.5s ease-in-out infinite;
}
@keyframes tempLoadingEffect {
  0% { color: var(--c-txt-grey); }
  50% { color: var(--c-txt-secondary); }
  100% { color: var(--c-txt-grey); }
}


.table.-is-collapsed {
  --tc-w-1: 50%;
  --tc-w-2: 50%;
  /* --tc-w-3: 40%; */
}
.table .table__main {
  width: 100%;
  transition: width .3s ease-in-out;
  will-change: width;
}
.table.-is-collapsed .table__main {
  width: var(--table-collapsed-width);
  transition: width .3s ease-in-out;
  will-change: width;
}

.table.-is-collapsed >>> .table__cells .table__col:nth-child(1) {
  min-width: var(--tc-w-1) !important;
  max-width: var(--tc-w-1) !important;
}
.table.-is-collapsed >>> .table__cells .table__col:nth-child(2) {
  min-width: var(--tc-w-2) !important;
  max-width: var(--tc-w-2) !important;
}
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

@media only screen and (max-width: 1400px) {
  .table {
    --table-collapsed-width: 20vw;
  }
}
@media only screen and (max-width: 992px) {
  .table {
    min-width: 850px;
  }
  .table.-is-collapsed .table__main {
    width: 320px;
  }
  .table__sheet {
    width: calc(100% - 320px - 1rem);  
  }  
}


</style>