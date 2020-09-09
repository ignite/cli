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
.table__cells.-header {
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid var(--c-theme-secondary);
}

</style>