<template>
  <div>
    <div class="container">
      <div class="container__left">
        <div class="navbar">
          <navbar>
            <div class="navbar__logo">
              <h1>Starport</h1>    
              <p>Dev Tool</p>
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
  padding-left: var(--g-margin);
  width: 15vw;  
  max-width: 320px;  
}
.container__right {
  padding-right: var(--g-margin);
  max-width: 1024px;
  margin-top: var(--g-offset-top);  
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
}
.navbar__logo h1 {
  font-size: 1.25rem;
  font-weight: 800;
  margin-right: 0.5rem;
}
.navbar__logo p {
  font-size: 1.25rem;
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
