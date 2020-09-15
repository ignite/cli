<template>
  <div 
    :class="[
      'table__row', 
      isRowActive ? '-is-active' : '',
      isWithInnerSheet ? '-is-button' : ''
    ]"
    :role="isWithInnerSheet ? 'button' : ''"
    @click="$emit('row-clicked', rowId, rowData)"
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
    isRowActive: { type: Boolean, default: false },
    isWithInnerSheet: { type: Boolean, default: false }
  },
  computed: {
    ...mapGetters('cosmos/blocks', [
      'highlightedBlock',
      'isTableSheetActive'
    ]),    
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