<template>
  <div>
    <layout />
  </div>
</template>

<script>
import initBlockConnection from '@/mixins/blocks/initConnection'

import Layout from "@/layouts/Layout"

export default {
  mixins: [initBlockConnection],
  components: {
    Layout,
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

a:visited {
  color: inherit;
}
</style>