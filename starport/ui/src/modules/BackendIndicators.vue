<template>
  <div class="container">
    <div class="container__title">
      <Headline>Backend Status</Headline>
    </div>
    <div class="indicators">
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
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import TooltipWrapper from '@/components/tooltip/TooltipWrapper'
import Headline from '@/components/typography/Headline'

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
    Headline
  },  
  data() {
    return {
      localHosts
    }
  },
  computed: {
    ...mapGetters('cosmos', [ 'backendRunningStates' ]),    
  },  
}
</script>

<style scoped>

.indicators > div:not(:last-child) {
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
  font-size: 1rem;
  color: var(--c-txt-grey);
  opacity: .8;
}
.chip__main a {
  text-decoration: none;
}
.chip.-is-active .chip__head-icon {
  animation: tempActiveEffect 2s ease-in-out infinite;
}
.chip.-is-active .chip__main a {
  color: var(--c-txt-grey);
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

@media only screen and (max-width: 1200px) {
  /*
   *
   * ⚠️TODO:
   * Temporarily hiding tooltip in sidesheet because
   * it extends over sidesheet width (hence hidden).
   * 
   * The feature to solve this problem 
   * will be implemented in TooltilWrapper component
   * with direction options.
   *
   */
  .indicators >>> .tooltip-wrapper .tooltip {
    display: none;
  }
}

</style>