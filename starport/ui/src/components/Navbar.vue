<template>
  <nav>
    <div class="nav__main">
      <router-link to="/">      
        <div class="nav__logo">
          <span><LogoStarport/></span>
          <h1>Starport</h1>
        </div>
      </router-link>      
    </div>
    <div class="nav__center">
      <div v-if="this.$route.name === 'welcome'" class="nav__center-msg -f-cosmos-overline-0">Welcome to Starport</div>
    </div>
    <div class="nav__sub">
      <button class="nav__ham" @click="handleHamClick"><HamIcon/></button>      
      <div class="nav__tabs">
        <router-link
          class="tab"
          to="/"
        >
          <!-- <div class="tab__icon"><CompassIcon/></div> -->
          Welcome
        </router-link>
        <router-link
          class="tab -flex"
          to="/blocks"
        >
          <span :class="['circle', {'-is-active': isBlinking}]" ref="circle"></span>
          <!-- <div class="tab__icon"><BlocksIcon/></div> -->
          Blocks
        </router-link>
      </div>      
    </div>
  </nav>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'

import BackendIndicators from '@/modules/BackendIndicators'
import LogoStarport from '@/assets/logos/LogoStarportSmall'
import HamIcon from '@/assets/icons/Ham'
import CompassIcon from '@/assets/icons/Compass'
import BlocksIcon from '@/assets/icons/Blocks'

export default {
  components: {
    LogoStarport,
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
    ...mapGetters('cosmos/blocks', [ 'latestBlock' ]),
  },
  methods: {
    handleHamClick() {
      this.$emit('ham-clicked')
    }
  },
  watch: {
    latestBlock() {
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
  width: 100%;
  height: var(--header-height);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav__tabs {
  display: flex;
}
.nav__tabs .tab:not(:last-child) {
  margin-right: 1.8rem;
}

.nav__ham {
  display: none;
  height: 1.5rem;
  width: 1.5rem;
}

.nav__logo {
  display: flex;
  align-items: center;
}
.nav__logo span {
  display: block;
  height: 26px;
  margin-right: 8px;
}
.nav__logo svg {
  width: auto;
  height: 100%;
}
.nav__logo h1 {
  font-size: 20px;
  font-weight: var(--f-w-bold);
  letter-spacing: -0.016em;
}

.nav__center-msg {
  color: var(--c-txt-third);
}

@media only screen and (max-width: 768px) {
  .nav__center-msg {
    display: none;
  }
}
@media only screen and (max-width: 576px) {
  nav {
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
  color: var(--c-txt-third);
}
.tab:hover {
  color: var(--c-txt-primary);
  transition: color .3s;
}
.tab:hover .tab__icon svg >>> path {
  fill: var(--c-txt-primary);
  transition: fill .3s;
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
  left: -0.8rem;
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
