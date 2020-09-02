<template>
  <div>
    <div class="container">
      <div class="navbar">
        <navbar>
          <div class="navbar__logo">
            <h1>Starport</h1>    
            <p>Dev Tool</p>
          </div>          
        </navbar>
      </div>
      <div class="content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  --g-margin: 3rem;
  --g-offset-top: 10rem;
}
.container {
  display: flex;
}

.navbar {
  padding-left: var(--g-margin);
  max-width: 30vw;
}
.navbar__logo {
  height: var(--g-offset-top);
  margin-left: 7px;

  display: flex;
  align-items: center;
}
.navbar__logo h1 {
  font-size: 1.5rem;
  font-weight: 800;
  margin-right: 0.5rem;
}
.navbar__logo p {
  font-size: 1.4rem;
}

.content {
  -webkit-overflow-scrolling: touch;
  flex-grow: 1;
  overflow: scroll;
}
.content {
  margin-top: var(--g-offset-top);
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
