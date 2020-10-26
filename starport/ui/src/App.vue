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
    Layout,
  },
  methods: {
    ...mapMutations('cosmos/env', [ 'setTimer', 'clearTimer' ]),
    ...mapActions('cosmos/env', [ 'setStatusState' ]),
    ...mapActions('cosmos/blocks', [ 'initBlockConnection' ]),
  },
  async created() {
    /*
     *
     // 1. Fetch backend status regularly
     *
     */
    this.timer = setInterval(this.setStatusState.bind(this), 5000)
    
    try {
      await this.setStatusState()
    } catch {
      console.log(`Can't fetch /env`)
    }

    /*
     *
     // 2. Start block fetching 
     *
     */
    this.initBlockConnection()
  },
  beforeDestroy() {
    this.clearTimer()
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