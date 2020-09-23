<template>
  <nav>
    <div class="nav__main">
      <button class="nav__ham" @click="handleHamClick"><HamIcon/></button>
      <div class="nav__logo">
        <slot />
      </div>
      <div class="nav__tabs">
        <router-link
          class="tab"
          to="/"
        >
          Welcome
        </router-link>
        <router-link
          class="tab -flex"
          to="/blocks"
        >
          <span :class="['circle', {'-is-active': isBlinking}]" ref="circle"></span>
          Blocks
        </router-link>
      </div>
    </div>
    <div class="nav__sub">
      <div class="nav__sub-chips">
        <BackendIndicators/>
      </div>
    </div>
  </nav>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'

import BackendIndicators from '@/modules/BackendIndicators'
import HamIcon from '@/assets/icons/Ham'

export default {
  components: {
    HamIcon,
    BackendIndicators
  },
  data() {
    return {
      isBlinking: false
    }
  },
  computed: {
    ...mapGetters('cosmos/blocks', [ 'blocksStack' ]),
  },
  methods: {
    handleHamClick() {
      this.$emit('ham-clicked')
    }
  },
  watch: {
    blocksStack() {
      /* TODO: refactor this temp code */
      this.isBlinking = true
      setTimeout(function() {
        this.isBlinking = false
      }.bind(this), 1500)
    }
  }
}

</script>

<style scoped>

nav {
  position: sticky;
  top: 0;
  z-index: 2;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding-bottom: 4rem;
  box-sizing: border-box;
}
.nav__tabs {
  margin-top: -1rem;
}
.nav__ham {
  display: none;
  height: 1.5rem;
  width: 1.5rem;
}
.nav__sub-chips {
  margin-left: 0.85rem;
}
@media only screen and (max-width: 1200px) {
  nav {
    height: auto;
    padding-bottom: 0;
    
  }
  .nav__ham {
    display: block;
  }
  .nav__main {
    display: grid;
    grid-template-columns: 30% 1fr 30%;
    align-items: center;    
    padding: 0 0.25rem;
  }
  .nav__sub {
    display: none;
  }
  .nav__tabs {
    display: flex;
    justify-content: flex-end;
    margin-top: 0;
  }
}

.tab {
  position: relative;
  display: block;
  text-decoration: none;

  font-size: 1.25rem;
  font-weight: 300;
  color: var(--c-txt-grey);

  padding: 0.85rem 0 0.85rem 1.85rem;
  border-radius: var(--bd-radius-primary);
  transition: background-color .3s;
}
.tab:hover {
  background-color: var(--c-bg-secondary);
  transition: background-color .3s;
}
.tab:not(:last-child) {
  margin-bottom: 0.5rem;
}
.tab.router-link-exact-active {
  pointer-events: none;  
  font-weight: 600;
  color: var(--c-txt-secondary);
  background-color: var(--c-bg-secondary);
}
@media only screen and (max-width: 1200px) {
  .tab {
    padding: 0;
    font-size: 1.125rem;
  }
  .tab:hover {
    background-color: transparent;
    color: var(--c-txt-primary);
    transition: color .3s;
  }  
  .tab:not(:last-child) {
    margin-right: 2.25rem;
    margin-bottom: 0;
  }  
  .tab.router-link-exact-active {
    background-color: transparent;
  }
}

/* temp loading effect */
.circle {
  --c-active: #4ACF4A;
  --circle-size: 8px;
}
.circle {
  border-radius: 100%;
  box-shadow: none;
  display: block;
  width: var(--circle-size);
  height: var(--circle-size);
  padding: 0;
  position: absolute;
  top: calc(50% - calc(var(--circle-size) / 1.8));
  left: 12px;
}
.circle::before,
.circle::after {
  border-radius: 100%;
}
.circle.-is-active {
  animation: tempActiveEffect 1s;
  animation-iteration-count: 1;
}
@keyframes tempActiveEffect {
  0% { background-color: var(--c-active); }
  50% { background-color: transparent; }
  75% { background-color: var(--c-active); }
  100% { background-color: transparent; }
}
@media only screen and (max-width: 1200px) {
  .circle {
    --circle-size: 6px;
  }  
  .circle {
    left: -0.8rem;
  }
}

</style>
