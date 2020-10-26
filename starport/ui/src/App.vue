<template>
  <div>
    <layout />
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import initCosmos from '@/mixins/initCosmos'
import Layout from '@/layouts/Layout'

export default {
  mixins: [initCosmos],
  components: {
    Layout,
  },
  methods: {
    ...mapActions('cosmos/blocks', [ 'initBlockConnection' ]),
  },
  created() {
    this.initBlockConnection()
  }  
};
</script>

<style>
body {
  margin: 0;
  font-family: var(--f-primary);
  color: var(--c-txt-primary);
  background-color: var(--c-bg-primary);
}

button:hover {
  cursor: pointer;
}

a {
  color: var(--c-txt-primary);
}
a:visited {
  color: inherit;
}
</style>