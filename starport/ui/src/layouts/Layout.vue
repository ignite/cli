<template>

  <div :class="['layout', `-route-${$route.name}`]">
    <div class="container">
      <div class="container__left">
        <div class="navbar">
          <Navbar @ham-clicked="handleHamClick">
            <div class="navbar__logo">
              <h1>Starport</h1>    
            </div>          
          </Navbar>
        </div>
      </div>
      <div class="container__right">
        <div class="content">
          <transition name="fade" mode="out-in">
            <keep-alive include="Welcome">
              <router-view />
            </keep-alive>
          </transition>
        </div>
      </div>

      <div v-if="this.$route.path === '/'" class="bg-image"><Planet/></div>      
    </div>

    <Modal :visible="visible" v-if="visible" @visible="visible = $event">
      <div v-if="modalTrigger === 'ham'" class="sheet -nav">
        <div class="sheet__main">
          <div class="sheet__links">
            <div class="sheet__links-headline">
              <Headline>Pages</Headline> 
            </div>
            <router-link
              class="tab"
              to="/"
            >Welcome</router-link>
            <router-link
              class="tab -flex"
              to="/blocks"
            >Blocks</router-link>          
          </div>
        </div>
        <div class="sheet__sub">
          <BackendIndicators/>
        </div>
      </div>
    </Modal>    
  </div>
</template>

<script>
import Navbar from '@/components/Navbar'
import Modal from '@/components/Modal'
import BackendIndicators from '@/modules/BackendIndicators'
import Headline from '@/components/typography/Headline'
import Planet from "@/assets/images/Planet.vue";

export default {
  components: {
    Navbar,
    Modal,
    BackendIndicators,
    Headline,
    Planet
  },  
  data() {
    return {
      visible: false,
      modalTrigger: 'ham'
    }
  },
  methods: {
    handleHamClick() {
      this.visible = !this.visible
      this.modalTrigger = 'ham'
    }
  },
}
</script>

<style scoped>
.layout {
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
.layout.-route-welcome .container__right {
  /* margin-left: 0; */
  max-width: 1120px;
}
@media only screen and (max-width: 1200px) {
  .layout {
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
  font-size: 28px;
  font-weight: var(--f-w-bold);
  letter-spacing: -0.016em;
}
@media only screen and (max-width: 1200px) {
  .navbar__logo {
    margin-left: 0;
    justify-content: center;
  }
}


.sheet {
  height: 100%;
  padding: 2rem var(--g-margin);
  box-sizing: border-box;
  overflow-x: hidden;
  /* display: flex;
  flex-direction: column;
  justify-content: space-between; */
}
.sheet__main {
  margin-bottom: 6rem;
}
.sheet__links-headline {
  margin-bottom: 1.5rem;
}

.tab {
  position: relative;
  display: block;
  text-decoration: none;

  font-size: 1.5rem;
  font-weight: 400;
  color: var(--c-txt-grey);

  transition: color .3s;
}
.tab:hover {
  color: var(--c-txt-primary);
  transition: color .3s;
}
.tab:not(:last-child) {
  margin-bottom: 1rem;
}
.tab.router-link-exact-active {
  pointer-events: none;  
  font-weight: 800;
  color: var(--c-txt-primary);
}

.bg-image {
  position: absolute;
  top: 0;
  right: 0;
  z-index: -1;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity .3s ease-in-out;
}
.fade-enter,
.fade-leave-active {
  opacity: 0;
  transition: opacity .3s ease-in-out;
}
.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: opacity .3s ease-in-out;
}
.fade-fast-enter,
.fade-fast-leave-active {
  opacity: 0;
  transition: opacity .3s ease-in-out;
}

</style>