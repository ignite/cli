<template>
  <div class="blocks-page">
    <SpBlockExplorer/>
  </div>
</template>

<script>
import { SpBlockExplorer } from '@tendermint/vue'

export default {
  components: { SpBlockExplorer }
}
</script>

<style scoped>

.blocks-page >>> .explorer {
  height: calc(100vh - var(--header-height) - 1px - 2.25rem) !important;
}

</style>