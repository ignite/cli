<template>

  <div :class="['layout', `-route-${routeName}`]">
    <div class="navbar">
      <Navbar @ham-clicked="handleHamClick" />
    </div>

    <div class="container">
      <div class="container__main">
        <div class="content">
          <transition name="fade" mode="out-in" @enter="enter">
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
      modalTrigger: 'ham',
      routeName: 'welcome'
    }
  },
  methods: {
    handleHamClick() {
      this.visible = !this.visible
      this.modalTrigger = 'ham'
    },
    enter() {
      this.routeName = this.$route.name
    }
  },
  mounted() {
    this.routeName = this.$route.name
  }
}
</script>

<style scoped>
.layout {
  --g-offset-side: 3rem;
  --g-offset-top: 10rem;
  --header-height: 72px;
}

.container {
  display: flex;
  justify-content: space-between;
}
.container__main {
  z-index: 0;
  padding: 0 var(--g-offset-side);
  max-width: 1120px;
  margin: 0 auto;
  flex-grow: 1;
  -webkit-overflow-scrolling: touch;  
}
.layout.-route-blocks .container__main {
  padding: 0;
  max-width: 100%;
}
@media only screen and (max-width: 1200px) {
  .layout {
    --g-offset-side: 2rem;
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
    padding: 0 var(--g-offset-side);
    margin-bottom: 2rem;
    border-bottom: 1px solid var(--c-txt-contrast-secondary);
  }
  .container__main {
    padding: 0 var(--g-offset-side);
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
  }
}


.navbar {
  /* position: sticky;
  z-index: 2;
  top: 0; */
  width: 100%;
  max-width: 100vw;
  padding: 0 var(--g-offset-side);
  box-sizing: border-box;
  border-bottom: 1px solid var(--c-border-primary);
}


.sheet {
  height: 100%;
  padding: 2rem var(--g-offset-side);
  box-sizing: border-box;
  overflow-x: hidden;
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

</style>