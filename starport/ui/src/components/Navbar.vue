<template>
  <nav>
    <div class="nav__main">
      <button class="nav__ham"><HamIcon/></button>
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
      <div class="nav__chips">

        <div 
          v-for="(chip, index) in localHosts"
          :key="chip.id+index"
          :class="['chip', {'-is-active': backendRunningStates[chip.id]}]"
        >
          <div class="chip__head">
            <span class="chip__head-icon"></span>
          </div>
          <TooltipWrapper :content="backendRunningStates[chip.id] ? chip.noteActive : chip.noteInactive">
            <p class="chip__main"><a :href="`localhost:${chip.url}`">{{chip.name}}</a></p>
          </TooltipWrapper>            
        </div>

      </div>
    </div>
  </nav>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'

import TooltipWrapper from '@/components/tooltip/TooltipWrapper'
import HamIcon from '@/assets/Ham'

const localHosts = [
  {
    id: 'frontend',
    name: 'User interface',
    noteActive: `The front-end of your app. A Vue application is running on port <span>8080</span>`,
    noteInactive: `Can't start UI, because Node.js is not found.`,
    url: '8080'
  },
  {
    id: 'rpc',
    name: 'Cosmos',
    noteActive: `The back-end of your app. Cosmos API is running locally on port <span>1317</span>`,
    noteInactive: ``,
    url: '1317'
  },
  {
    id: 'api',
    name: 'Tendermint',
    noteActive: `The consensus engine. Tendermint API is running on port <span>26657</span>`,
    noteInactive: ``,
    url: '26657'
  }
]

export default {
  components: {
    TooltipWrapper,
    HamIcon
  },
  data() {
    return {
      isBlinking: false,
      localHosts
    }
  },
  computed: {
    ...mapGetters('cosmos', [ 'backendRunningStates' ]),    
    ...mapGetters('cosmos/blocks', [ 'blockEntries' ]),
  },
  watch: {
    blockEntries() {
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

/* chip */
.nav__chips {
  margin-left: 0.85rem;
}
.nav__chips > div:not(:last-child) {
  margin-bottom: 0.65rem;
}
.chip {
  --c-active: #4ACF4A;
  --c-active-sub: #7fe87f;
}
.chip {
  display: flex;
  align-items: center;
}
/* .chip:not(:last-child) {
  margin-bottom: 0.65rem;
} */
.chip__head-icon {
  display: block;
  width: 6px;
  height: 6px;
  border-radius: 100%;  
  margin-right: 0.65rem;
  margin-top: 0.25px;
  background-color: var(--c-txt-grey);  
  animation: tempLoadingEffect 1.5s ease-in-out infinite;
} 
.chip__main {
  font-size: 0.9375rem;
  color: var(--c-txt-grey);
}
.chip__main a {
  text-decoration: none;
}
.chip.-is-active .chip__head-icon {
  animation: tempActiveEffect 2s ease-in-out infinite;
}
.chip.-is-active .chip__main {
  color: var(--c-txt-secondary);
}
@keyframes tempLoadingEffect {
  0% { background-color: var(--c-txt-grey); }
  50% { background-color: var(--c-txt-secondary); }
  100% { background-color: var(--c-txt-grey); }
}
@keyframes tempActiveEffect {
  0% { background-color: var(--c-active); }
  50% { background-color: var(--c-active-sub); }
  100% { background-color: var(--c-active); }
}

</style>
