<template>
  <div>
    <router-view />
  </div>
</template>

<style>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700;800;900&display=swap");

body {
  margin: 0;
  font-family: "Inter", "Helvetica", sans-serif;
}
</style>

<script>
export default {
  created() {
    this.$store.dispatch("init");
  }
};
</script>
