<template>
  <div>
    <layout />
  </div>
</template>

<style>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700;800;900&display=swap");
@import './styles/_normalize.css';
@import './styles/_colors.css';
@import './styles/_typography.css';

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
import Layout from "@/components/Layout";

export default {
  components: {
    Layout,
  },
};
</script>
