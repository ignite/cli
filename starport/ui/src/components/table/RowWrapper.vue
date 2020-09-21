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
  }
}
</script>

<style scoped>

.table__row {
  border-radius: var(--bd-radius-primary);
  font-size: 0.9rem;
}
.table__row.-is-active {
  background-color: var(--c-bg-primary);
}
.table__row.-is-active >>> .table__cells {
  color: var(--c-txt-primary);
  font-weight: 500;
}
.table__row.-is-button {
  cursor: pointer;
}
.table__row:hover {
  background-color: var(--c-bg-primary);
  transition: background-color .3s;
}
.table__row { transition: background-color .3s; }

.table__row >>> .accord-item__contents .side-tab-list {
  margin-top: 1rem;
  padding-bottom: 1rem;
}

</style>