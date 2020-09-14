<template>
  <div>
    <layout />
  </div>
</template>

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

<script>
import Layout from "@/layouts/Layout";

export default {
  components: {
    Layout,
  },
};
</script>
