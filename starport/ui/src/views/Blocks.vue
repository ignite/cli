<template>
  <div>
    <!-- <BlocksExplorer/> -->
    <SpBlockExplorer/>
  </div>
</template>

<script>
import { SpBlockExplorer } from '@tendermint/vue'
// import BlocksExplorer from '@/modules/BlocksExplorer'

export default {
  components: {
    // BlocksExplorer,
    SpBlockExplorer
  }
}
</script>

<style scoped>
</style>