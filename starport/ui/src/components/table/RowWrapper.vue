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
        const activeRowId = $table.rowState.activeRowId
        return activeRowId === this.rowId
      }

      return false
    }
  },
  methods: {
    getParentTableNode(parentNode) {
      if (parentNode) {
        if (parentNode.$refs.table === undefined) this.getParentTableNode(parentNode.$parent)
        return parentNode
      }

      return null
    },
    handleClick() {
      const $table = this.getParentTableNode(this.$parent)
      const activeRowId = $table.rowState.activeRowId
      const isSheetActive = $table.sheetState.isActive
      const isActiveRowClicked = activeRowId === this.rowId
      
      if (isSheetActive) {
        if (isActiveRowClicked) {
          $table.sheetState.isActive = false
          $table.rowState.activeRowId = null
        } else {
          $table.rowState.activeRowId = this.rowId
        }
      } else {
        $table.sheetState.isActive = true
        $table.rowState.activeRowId = this.rowId
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