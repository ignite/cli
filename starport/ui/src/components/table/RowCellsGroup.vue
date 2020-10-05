<template>
  <div>
    <div :class="['table__cells', isTableHead ? '-header' : '']">
      <div 
        v-for="(n, index) in tableCells"
        :key="n+index"
        class="table__col"
      >{{n}}</div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    isTableHead: { type: Boolean, default: false },
    tableCells: { type: Array, required: true }
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