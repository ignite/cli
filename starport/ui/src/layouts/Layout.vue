<template>
  <div>
    <div class="container">
      <div class="container__left">
        <div class="navbar">
          <Navbar @ham-clicked="handleHamClick">
            <div class="navbar__logo">
              <h1>Starport</h1>    
              <p>DevTools</p>
            </div>          
          </Navbar>
        </div>
      </div>
      <div class="container__right">
        <div class="content">
          <router-view />
        </div>
      </div>
    </div>
    <Modal :width="'40vw'" :visible="visible" v-if="visible" @visible="visible = $event">
      <div v-if="modalTrigger === 'ham'" class="sheet -nav">
        <div class="sheet__main">
          <router-link
            class="tab"
            to="/"
          >Welcome</router-link>
          <router-link
            class="tab -flex"
            to="/blocks"
          >Blocks</router-link>          
        </div>
        <div class="sheet__sub">
          <p>test</p>
        </div>
      </div>
    </Modal>    
  </div>
</template>

<script>
import Navbar from '@/components/Navbar'
import Modal from '@/components/Modal'

export default {
  data() {
    return {
      visible: false,
      modalTrigger: 'ham'
    }
  },
  components: {
    Navbar,
    Modal
  },
  methods: {
    handleHamClick() {
      this.visible = !this.visible
      this.modalTrigger = 'ham'
    }
  }
};
</script>

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

@media only screen and (max-width: 1200px) {
  .container {
    --g-margin: 2rem;
    --g-offset-top: 5rem;
  }  
  .container {
    flex-direction: column;
  }
  .container__left {
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
    margin-right: 0;
    padding: 0 var(--g-margin);
    margin-bottom: 2rem;
    border-bottom: 1px solid var(--c-txt-contrast-secondary);
  }
  .container__right {
    padding: 0 var(--g-margin);
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
  }
}

.sidesheet {
  
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
@media only screen and (max-width: 1200px) {
  .navbar__logo {
    margin-left: 0;
    justify-content: center;
  }
}
</style>