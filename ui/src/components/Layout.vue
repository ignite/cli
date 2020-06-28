<template>
  <div>
    <div class="container">
      <!-- <div class="navbar">
        <navbar />
      </div> -->
      <div class="content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  --navbar-height: 3rem;
}
.navbar {
  height: var(--navbar-height);
  width: 100%;
  position: fixed;
}
.content {
  padding-top: var(--navbar-height);
  overflow: scroll;
  -webkit-overflow-scrolling: touch;
}
</style>

<script>
import Navbar from "@/components/Navbar";

export default {
  components: {
    Navbar,
  },
};
</script>
