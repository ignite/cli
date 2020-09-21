<template>
  <div>
    <div class="container">
      <div class="container__left">
        <div class="navbar">
          <navbar>
            <div class="navbar__logo">
              <h1>Starport</h1>    
              <p>DevTools</p>
            </div>          
          </navbar>
        </div>
      </div>
      <div class="container__right">
        <div class="content">
          <router-view />
        </div>
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
  justify-content: space-between;
}
.container__left {
  z-index: 1;
  padding-left: calc(var(--g-margin) - 1.5rem);
  width: 12vw;  
  min-width: 200px;  
  max-width: 320px;  
  margin-right: 2.5rem;  
}
.container__right {
  z-index: 0;
  padding-right: var(--g-margin);
  max-width: 1420px;
  margin-top: calc(var(--g-offset-top) / 1.85);  
  margin-right: auto;
  margin-left: auto;
  flex-grow: 1;
  -webkit-overflow-scrolling: touch;  
}

.navbar {
  position: sticky;
  top: 0;
}
.navbar__logo {
  height: var(--g-offset-top);

  display: flex;
  align-items: center;

  margin-left: 1.5rem;
}
.navbar__logo h1 {
  font-size: 1.125rem;
  font-weight: 800;
  margin-right: 0.5rem;
}
.navbar__logo p {
  font-size: 1.125rem;
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
