<template>
  <div>
    <div :class="['table__cells', isTableHead ? '-header' : '']">
      <div 
        v-for="(n, index) in tableCells"
        :key="n.id"
        class="table__col"
        :style="[getCellWidth(cellWidths[index])]"
      >{{n.content}}</div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    isTableHead: { type: Boolean, default: false },
    tableCells: { type: Array, required: true },
    cellWidths: {
      type: Array,
      default() {
        return ['8%', '10%', '1', '15%']
      },
      validator(value) {
        return value.filter(val => typeof val === 'string').length === value.length
      }      
    }
  },
  methods: {
    getCellWidth(val) {
      if (val.includes('%')) {
        return {
          maxWidth: val,
          minWidth: val
        }
      }

      return { flexGrow: val }
    }
  }
}
</script>

<style scoped>

.table__cells {
  display: flex;
}
div.table__cells.-header {
  padding-top: 1.25rem;
  padding-bottom: 1rem;
  /* border-bottom: 1px solid var(--c-theme-secondary); */
  color: var(--c-txt-secondary);
  margin-bottom: 0.5rem;
  font-weight: 500;
  font-size: 0.9375rem;
}
div.table__cells:not(.-header) {
  color: var(--c-txt-grey);
}

</style>