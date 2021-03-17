<template>
  <div>
    <layout />
  </div>
</template>

<script>
import { mapActions, mapMutations } from 'vuex'
import Layout from '@/layouts/Layout'

export default {
  components: {
    Layout
  },
  async created() {
    this.$store.dispatch("cosmos/init", { account: false })    
  }
};
</script>

<style>
body {
  margin: 0;
  font-family: var(--sp-f-primary);
  color: var(--sp-c-txt-primary);
  background-color: var(--sp-c-bg-primary);
}

button:hover {
  cursor: pointer;
}

a {
  color: var(--sp-c-txt-primary);
}
/* a:visited {
  color: inherit;
} */
</style>