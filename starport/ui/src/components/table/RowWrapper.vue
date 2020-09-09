<template>
  <div 
    :class="[
      'table__row', 
      isActive ? '-is-active' : '',
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
    isWithInnerSheet: { type: Boolean, default: false }
  },
  data() {
    return {
      isActive: false
    }
  },
  methods: {
    getParentTableNode(parentNode) {
      if (parentNode.$refs.table === undefined) this.getParentTableNode(parentNode.$parent)
      return parentNode
    },
    handleClick() {
      const $table = this.getParentTableNode(this.$parent)
      $table.isSheetActive = !$table.isSheetActive
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

.table__row >>> .accord-item__contents .side-tab-list {
  margin-top: 1rem;
  padding-bottom: 1rem;
}

</style>