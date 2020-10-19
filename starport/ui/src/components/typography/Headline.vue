<template>
  <p class="headline"><slot/></p>
</template>

<script>
export default {
  props: {
    withBorder: {
      type: Boolean,
      default: true
    }
  }
}
</script>

<style scoped>

.headline {
  display: block;
  position: relative;
  font-size: 0.8125rem;
  color: var(--c-txt-grey);
  padding-bottom: 6px;
  margin-bottom: 1rem;
}
.headline:after {
  content: '';
  position: absolute;
  left: -6px;
  bottom: 0;
  width: 100%;
  height: 1px;
  background-color: var(--c-txt-contrast-secondary);
}

</style>