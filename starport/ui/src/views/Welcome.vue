<template>

  <div class="main">

    <div class="hero">
      <h4>Welcome to Starport!</h4>
      <h2>Your blockchain <br> is ready.</h2>
      <p>Starport has scaffolded and launched a Cosmos blockchain for you. Your blockchain has its own tokens, accounts, governance, custom data types and more.</p>
    </div>

    <div class="dashboard grid-col-3">
      <div class="left-top dashboard__headline">BUILD LOG</div>
      <div class="center-top dashboard__headline">STACK</div>

      <div class="left dashboard__card dashboard__log">
        <IconItem :iconType="'check'"  :itemText="'Depencies installed'" />        
        <IconItem :iconType="'check'"  :itemText="'Source code scaffolded'" />        
        <IconItem :iconType="'check'"  :itemText="'Build complete'" />        
        <IconItem :iconType="'check'"  :itemText="'Blockchain initialized'" />        
        <IconItem :iconType="'check'"  :itemText="'Accounts created'" />        
        <IconItem :iconType="'check'"  :itemText="'Blockchain node started'" />        
      </div>

      <div 
        v-for="(card, index) in stack"
        :key="card.id+index"
        :class="['dashboard__card', `-${card.id}`]"
      >
        <div class="dashboard__card-logo">
          <LogoCosmosSdk v-if="card.id === 'rpc'" />
          <LogoTendermint v-else-if="card.id === 'api'" />
          <LogoStarport v-else />
        </div>
        <div class="dashboard__card-main">
          <span class="dashboard__card-heading">{{card.name}}</span>
          <p class="dashboard__card-blurb">{{card.blurb}}</p>
          <IconItem :isActive="backendRunningStates[card.id]" :itemText="`localhost: ${card.port}`" />        
        </div>
        <div 
          v-if="card.id === 'api' && blockStack.length>0"
          class="dashboard__card-blocks"
        >

          <transition-group name="list" tag="ul">
            <div 
              v-for="(block, index) in blockStack"
              :key="block.hash"
              class="card-counter"
            >
              <div class="card-counter__top">
                <div class="card-counter__top-left">
                  <p>BLOCK</p>
                  <p>{{block.height}}</p>
                </div>
                <div class="card-counter__top-right">
                  <span>{{block.time}}</span>
                </div>
              </div>
              <div class="card-counter__btm">
                <p ref="blockHash" class="card-counter__hash">{{block.hash}}</p>
              </div>
              <div class="card-counter__bg">
                <Box/>
              </div>
            </div>
          </transition-group>

        </div>        
      </div>      

    </div>

    <div class="grid-col-3 intro">
      <div class="intro__side">
        <span>Architecture</span>
        <h3>Brief intro</h3>
      </div>
      <div class="intro__main">
        <p>Your blockchain is built with Cosmos SDK, a modular framework for building blockchains. It includes modules such as auth, bank, staking, governance, and more. Every feature is packaged as a separate module that can interact with other modules. Starport has actually generated a module for you, which you can use to start developing your own application and features.</p>
      </div>
    </div>
    
    <div class="tutorials">
      <div class="tutorials__top">
        <h3>Build your app</h3>
      </div>

      <div class="tutorials__articles">
        <div class="grid-col-3 cards">
          <a 
            v-for="card in articles"
            :key="card.title"
            :class="['text-card', { '-is-dark': card.tagline === 'tutorial' }]"
            :href="card.link"
          >
            <div class="text-card__top">
              <span class="text-card__tagline">{{card.tagline}}</span>
              <p class="text-card__title">{{card.title}}</p>
            </div>
            <div class="text-card__btm">
              <p>{{card.blurb}}</p>
            </div>
          </a>
        </div>
      </div>

      <div class="tutorials__videos">
        <a 
          v-for="card in videos"
          :key="card.title"
          class="image-card"
          :href="card.link"
        >
          <img class="image-card__img" src="https://i.ytimg.com/vi/h6Ur_40LB9k/hq720.jpg" alt="Image">
          <div class="image-card__text">
            <span class="image-card__text-h1">{{card.title}}</span>
            <span class="image-card__text-p">{{card.length}}</span>
          </div>          
        </a>
      </div>
    </div>

    <div class="footer">
      <div 
        v-for="block in footerBlocks"
        :key="block.title"
        class="footer__block"
      >
        <span>{{block.title}}</span>
        <a :href="block.link.url">{{block.link.text}}</a>
      </div>
    </div>

  </div>
  
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import moment from 'moment'

import IconItem from '@/components/list/IconItem'
import LogoTendermint from "@/assets/LogoTendermint.vue";
import LogoCosmosSdk from "@/assets/LogoCosmosSdk.vue";
import LogoStarport from "@/assets/LogoStarport.vue";
import Box from "@/assets/icons/Box.vue";

const stack = [
  {
    id: 'api',
    name: 'Blockchain',
    blurb: `The consensus engine, powered by Tendermint Core.`,
    port: '26657'
  },
  {
    id: 'rpc',
    name: 'API Server',
    blurb: `The back-end of your app, powered by Cosmos SDK.`,
    port: '1317'
  },
  {
    id: 'frontend',
    name: 'User Interface',
    blurb: `The front-end of your app, built with Vue.js and generated by Starport.`,
    port: '8080'
  }
]

const articles = [
  {
    tagline: 'tutorial',
    title: 'Build a Polling App',
    blurb: 'Build a voting application with a web-based UI.',
    link: '#'
  },
  {
    tagline: 'documentation',
    title: 'Starport Handbook',
    blurb: 'Create your own blockchain - from star to ecosystem',
    link: '#'
  },
  {
    tagline: 'tutorial',
    title: 'Build a Blog',
    blurb: 'Learn how Starport works by building a blog.',
    link: '#'
  },
]

const videos = [
  {
    title: 'Cosmos Code With Us - Building your first Cosmos app',
    length: '1:39:07',
    link: '#'
  },
  {
    title: 'Getting started with Starport, the easiest way to build a Cosmos SDK blockchain',
    length: '3:31',
    link: '#'
  },
]

const footerBlocks = [
  { title: 'Chat with developers', link: { text: 'Cosmos Discord →', url: '#' } },
  { title: 'Join the community', link: { text: 'Cosmos SDK Forum →', url: '#' } },
  { title: 'Found an issue?', link: { text: 'Suggest improvements →', url: '#' } },
]

export default {
  components: {
    IconItem,
    LogoTendermint,
    LogoCosmosSdk, 
    LogoStarport,
    Box,
  },
  data() {
    return {
      stack,
      articles,
      videos,
      footerBlocks,
      blockStack: []
    }
  },
  computed: {
    ...mapGetters('cosmos', [ 'backendRunningStates', 'backendEnv' ]),   
    ...mapGetters('cosmos/blocks', [ 'latestBlock' ]), 
  },    
  methods: {
    getFmtBlockTime(block) {
      if (!block) return '_'

      const time = block.blockMeta.block.header.time
      return moment(time).format('H:mm:ss')
    }
  },
  watch: {
    latestBlock() {
      if (this.blockStack.length>2) this.blockStack.splice(0, 1)

      this.blockStack.splice(2, 0, {
        height: this.latestBlock.height,
        hash: this.latestBlock.blockMeta.block_id.hash,
        time: this.getFmtBlockTime(this.latestBlock),
      })
    }
  }
}
</script>

<style scoped>

.main {
  position: relative;
  margin-top: -1rem;
  margin-bottom: 6rem;
}

.hero {
  margin-bottom: 4rem;
  margin-right: 45%;
}
.hero h2 {
  font-size: 5.625rem;
  font-weight: var(--f-w-extra-bold);
  line-height: 108%;
  letter-spacing: -0.055em;
  margin-bottom: 2rem;
}

.hero h4 {
  font-size: 1rem;
  font-weight: var(--f-w-medium);
  text-transform: uppercase;
  color: rgba(0,5,66,.621);
  margin-bottom: 1rem;
}
.hero p {
  font-size: 1.3125rem;
  font-weight: var(--f-w-light);
  line-height: 145%;
  letter-spacing: -0.007em;  
}

.grid-col-3 {
  display: grid;  
  grid-column-gap: 2rem;      
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.dashboard__card.left-top { grid-area: left-top; }
.dashboard__card.center-top { grid-area: center-top; }
.dashboard__card.left { grid-area: left; }
.dashboard__card.-api { grid-area: api; }
.dashboard__card.-rpc { grid-area: rpc; }
.dashboard__card.-frontend { grid-area: frontend; }

.dashboard {
  grid-template-areas: 
    'left-top center-top center-top'
    'left api rpc'
    'left api frontend';
  grid-row-gap: 1rem;

  margin-bottom: 8rem;
}
.dashboard__card {
  position: relative;
  background-color: #F8F8FD;
  border-radius: 12px;
  padding: 1.75rem 1.75rem;
}
.dashboard__card.-api {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.dashboard__card-heading {
  display: block;
  font-size: 1rem;
  font-weight: var(--f-w-bold);
  margin-bottom: 0.45rem;
}
.dashboard__card-blurb {
  font-size: 0.75rem;
  line-height: 130.9%;
  color: rgba(0,4,56,73.8%);
  margin-bottom: 2.5rem;
  width: 80%;
}
.dashboard__card-logo {
  position: absolute;
  top: 0.8rem;
  right: 0.8rem;
}
.dashboard__headline {
  font-size: 1rem;
  font-weight: var(--f-w-medium);
  color: rgba(0,5,66,62.1%);
}
.dashboard__log {
  background-color: transparent;
  border: 1px solid rgba(0,13,158, 7%);
}
.dashboard__log > div:not(:last-child) {
  margin-bottom: 1rem;
}

.dashboard__card-blocks {
  /* transform: translate3d(0, 4rem, 0);
  box-shadow: 0px 8px 40px rgba(0, 3, 66, 0.08); */
  /* margin-top: 2rem; */
}
.card-counter {
  position: relative;
  padding: 1.5rem;
  background-color: var(--c-bg-primary);
  border-radius: 12px;
}
.card-counter__top {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5rem;
}
.card-counter__top-left {
  color: #4251fa;
}
.card-counter__top-left p:first-child {
  font-size: 0.75rem;
  font-weight: var(--f-w-medium);
  margin-bottom: 4px;
}
.card-counter__top-left p:last-child {
  font-size: 1.3125rem;
  font-weight: var(--f-w-bold);
}
.card-counter__top-right {
  font-size: 0.75rem;
  color: rgba(0, 5, 66, 0.621);
}
.card-counter__btm {
  color: rgba(0, 5, 66, 0.621);
}
.card-counter__bg {
  position: absolute;
  right: 0;
  bottom: 30%;
}
.card-counter__hash {
  display: block;
  white-space: nowrap; /* forces text to single line */
  overflow: hidden;
  text-overflow: ellipsis;  
}

.dashboard__card-blocks {
  height: 100%;
  position: relative;
  transform: translate3d(0, 4rem, 0);
  /* box-shadow: 0px 8px 40px rgba(0, 3, 66, 0.08);   */
  perspective: 1000px;
}
.card-counter {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 100%;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  box-shadow: 0px 8px 40px rgba(0, 3, 66, 0.08);  
  transform-origin: center;
}
.card-counter:nth-last-child(1) {
  z-index: 0;
}
.card-counter:nth-last-child(2) {
  transform: translate3d(0,-20px,-50px);
  z-index: -1;
  transition: transform .5s;
}
.card-counter:nth-last-child(3) {
  transform: translate3d(0,-40px,-100px);
  z-index: -2;
  transition: transform .5s;
}
.card-counter:nth-last-child(4) {
  transform: translate3d(0,-60px,-150px);
  z-index: -3;
  transition: transform .5s;
  opacity: 0;
}

.list-enter-active {
  animation: slideIn 1s;
}
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translate3d(0, 24px, 50px);
  }
  to {
    opacity: 1;
    transform: translate3d(0, 0, 0);
  }
}



.intro__main { grid-area: intro-main; }
.intro__side { grid-area: intro-side; }
.intro {
  grid-template-areas: 
    'intro-side intro-main intro-main';  
}
.intro__main {
  width: 80%;
}
.intro__side span {
  display: block;
  font-size: 1rem;
  font-weight: var(--f-w-medium);
  text-transform: uppercase;
  color: rgba(0,5,66,62.1%);  
  margin-bottom: 4px;
}
.intro__side h3 {
  font-size: 2.375rem;
  font-weight: var(--f-w-bold);
  margin-left: -2px;
}
.intro__main p {
  line-height: 162.5%;
}
.intro {
  margin-bottom: 5rem;
}

.tutorials__top {
  margin-bottom: 2rem;
}
.tutorials__top h3 {
  font-size: 2.375rem;
  font-weight: var(--f-w-bold);
  margin-left: -2px;  
}
.tutorials__articles {
  margin-bottom: 2rem;
}

.text-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0px 0px 1px rgba(0, 0, 0, 0.07), 0px 8px 16px rgba(0, 0, 0, 0.05), 0px 20px 44px rgba(0, 3, 66, 0.12);  
}
.text-card.-is-dark {
  color: var(--c-txt-contrast-primary);
  background: linear-gradient(124.57deg, #1E1741 0%, #222262 100%);
}
a.text-card {
  text-decoration: none;
}
.text-card__top {
  margin-bottom: 8rem;
}
.text-card__tagline {
  display: block;
  font-weight: var(--f-w-bold);
  font-size: 0.75rem;
  line-height: 130.9%;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #CFD1E7;
  margin-bottom: 4px;
}
.text-card__title {
  font-weight: var(--f-w-bold);
  font-size: 1.3125rem;
  line-height: 128.7%;
  letter-spacing: -0.007em;
}
.text-card__btm p {
  line-height: 130%;
  color: #CFD1E7;
  width: 90%;
}

.tutorials__videos {
  display: flex;
  justify-content: space-between;
}
.tutorials__videos .image-card {
  width: calc((100% - 2rem) / 2);
}

.image-card {
  text-decoration: none;
}
.image-card img {
  width: 100%;
  object-fit: cover;
  border-radius: 12px;
  box-shadow: 0px 0px 1px rgba(0, 0, 0, 0.07), 0px 8px 16px rgba(0, 0, 0, 0.05), 0px 20px 44px rgba(0, 3, 66, 0.12);  
  margin-bottom: 0.5rem;  
}
.image-card__text-h1 {
  display: block;
  font-weight: var(--f-w-bold);
  line-height: 130%;
}
.image-card__text-p {
  font-size: 0.75rem;
  line-height: 130.9%;
  letter-spacing: 0.005em;
  color: #616489;
}

.tutorials__videos {
  margin-bottom: 8rem;
}

.footer {
  display: flex;
  padding-top: 3rem;
  border-top: 1px solid rgba(0, 11, 119, 0.185);
}
.footer__block:not(:last-child) {
  margin-right: 10%;
}
.footer__block span {
  display: block;
  font-weight: var(--f-w-bold);
  line-height: 130%;
  margin-bottom: 0.5rem;
}
.footer__block a {
  text-decoration: none;
  font-size: 16px;
  letter-spacing: -0.007em;
  font-weight: var(--f-w-medium);
  color: #4251FA;
}

</style>