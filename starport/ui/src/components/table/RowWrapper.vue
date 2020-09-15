<template>
  <div 
    :class="[
      'table__row', 
      isRowActive ? '-is-active' : '',
      isWithInnerSheet ? '-is-button' : ''
    ]"
    :role="isWithInnerSheet ? 'button' : ''"
    @click="handleClick"
  >
    <slot/>
  </div>
</template>

<script>
import { mapGetters, mapMutations } from 'vuex'

export default {
  props: {
    rowData: { type: Object },
    rowId: { type: String },
    isWithInnerSheet: { type: Boolean, default: false }
  },
  computed: {
    ...mapGetters('cosmos/blocks', [
      'highlightedBlock',
      'isTableSheetActive'
    ]),    
    isRowActive() {
      return this.highlightedBlock.id === this.rowId
    }
  },
  methods: {
    ...mapMutations('cosmos/blocks', [
      'setHighlightedBlock',
      'setTableSheetState'
    ]),
    setTableRowStore(isToActive=false, payload) {
      const highlightBlockPayload = isToActive ? {
        id: payload.rowId,
        data: payload.rowData
      } : null
      
      this.setHighlightedBlock(highlightBlockPayload)
    },
    handleClick() {
      const isActiveRowClicked = this.isRowActive
      
      if (this.isTableSheetActive) {
        if (isActiveRowClicked) {
          this.setTableSheetState(false)
          this.setTableRowStore()
        } else {
          this.setTableRowStore(true, { rowId: this.rowId, rowData: this.rowData })
        }
      } else {
        this.setTableSheetState(true)
        this.setTableRowStore(true, { rowId: this.rowId, rowData: this.rowData })
      }
    }
  }
}
</script>

<style scoped>

.table__row.-is-active {
  background-color: var(--c-bg-secondary);
}
.table__row.-is-button {
  cursor: pointer;
}
.table__row:hover {
  background-color: var(--c-bg-secondary);
  transition: background-color .3s;
}
.table__row { transition: background-color .3s; }

.table__row >>> .accord-item__contents .side-tab-list {
  margin-top: 1rem;
  padding-bottom: 1rem;
}

</style>