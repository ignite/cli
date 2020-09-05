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

.table__row >>> .accord-item__contents .side-tab-list {
  margin-top: 1rem;
  padding-bottom: 1rem;
}

</style>