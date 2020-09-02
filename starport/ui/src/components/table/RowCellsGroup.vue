<template>
  <div>
    <div class="-row-group">
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
    tableCells: { type: Array, required: true }
  }
}
</script>

<style scoped>

</style>