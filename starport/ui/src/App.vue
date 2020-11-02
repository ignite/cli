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
  // methods: {
  //   ...mapMutations('cosmos/env', [ 'setTimer', 'clearTimer' ]),
  //   ...mapActions('cosmos/env', [ 'setStatusState' ]),
  //   ...mapActions('cosmos/blocks', [ 'initBlockConnection' ]),
  // },
  async created() {
    this.$store.dispatch("cosmos/init", { account: false })    
  },
  // beforeDestroy() {
  //   this.clearTimer()
  // }  
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