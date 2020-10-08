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
          <div class="tab__icon"><CompassIcon/></div>
          Welcome
        </router-link>
        <router-link
          class="tab -flex"
          to="/blocks"
        >
          <span :class="['circle', {'-is-active': isBlinking}]" ref="circle"></span>
          <div class="tab__icon"><BlocksIcon/></div>
          Blocks
        </router-link>
      </div>
    </div>
    <!-- <div class="nav__sub">
      <div class="nav__sub-chips">
        <BackendIndicators/>
      </div>
    </div> -->
  </nav>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'

import BackendIndicators from '@/modules/BackendIndicators'
import HamIcon from '@/assets/icons/Ham'
import CompassIcon from '@/assets/icons/Compass'
import BlocksIcon from '@/assets/icons/Blocks'

export default {
  components: {
    HamIcon,
    BackendIndicators,
    CompassIcon,
    BlocksIcon
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
@media screen and (max-width: 768px) {
  .nav__tabs {
    display: none;
  }
}

.tab {
  position: relative;
  text-decoration: none;

  display: flex;
  align-items: center;  

  font-weight: var(--f-w-medium);
  line-height: 130%;
  color: rgba(0, 5, 66, 0.621);

  padding: 0 0 0 1.4rem;
  border-radius: var(--bd-radius-primary);
  transition: background-color .3s;
}
.tab:hover {
  color: var(--c-txt-primary);
  transition: color .3s;
}
.tab:hover .tab__icon svg >>> path {
  fill: var(--c-txt-primary);
  transition: fill .3s;
}
.tab:not(:last-child) {
  margin-bottom: 1rem;
}
.tab__icon {
  display: flex;
  align-items: center;
  margin-right: 0.8rem;
}
.tab.router-link-exact-active {
  pointer-events: none;  
  color: var(--c-txt-primary);
}
.tab.router-link-exact-active .tab__icon svg >>> path {
  fill: #4251FA;
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
  --circle-size: 6px;
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
  left: 6px;
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
