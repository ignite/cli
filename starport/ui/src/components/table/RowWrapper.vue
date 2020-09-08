<template>
  <div :class="['table__row', isActive ? '-is-active' : '']">
    <slot/>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isActive: false
    }
  }
}
</script>

<style scoped>

.table__row.-is-active {
  background-color: snow;
}

</style>