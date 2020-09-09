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
export default {
  props: {
    rowData: { type: Object },
    rowId: { type: String },
    isWithInnerSheet: { type: Boolean, default: false }
  },
  data() {
    return {
      isActive: false
    }
  },
  computed: {
    isRowActive() {
      if (this.$parent) {
        const $table = this.getParentTableNode(this.$parent)
        const activeRowId = $table.rowStore.activeRowId
        return activeRowId === this.rowId
      }

      return false
    }
  },
  methods: {
    setTableRowStore(tableNode, isToActive=false, payload=null) {
      if (isToActive) {
        tableNode.rowStore.activeRowId = payload.rowId
        tableNode.rowStore.activeRowData = payload.rowData
      } else {
        tableNode.rowStore.activeRowId = null
        tableNode.rowStore.activeRowData = null
      }
    },
    getParentTableNode(parentNode) {
      if (parentNode) {
        if (parentNode.$refs.table === undefined) this.getParentTableNode(parentNode.$parent)
        return parentNode
      }

      return null
    },
    handleClick() {
      const $table = this.getParentTableNode(this.$parent)
      const activeRowId = $table.rowStore.activeRowId
      const isSheetActive = $table.sheetStore.isActive
      const isActiveRowClicked = activeRowId === this.rowId
      
      if (isSheetActive) {
        if (isActiveRowClicked) {
          $table.sheetStore.isActive = false
          this.setTableRowStore($table)
        } else {
          this.setTableRowStore($table, true, { rowId: this.rowId, rowData: this.rowData })
        }
      } else {
        $table.sheetStore.isActive = true
        this.setTableRowStore($table, true, { rowId: this.rowId, rowData: this.rowData })
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