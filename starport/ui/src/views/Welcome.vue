<template>

  <div class="main">

    <div class="hero">
      <h2>Your blockchain <br> is ready.</h2>
      <p>Starport has scaffolded and launched a Cosmos blockchain for you. Your blockchain has its own tokens, accounts, custom data types and more.</p>
    </div>

    <div class="dashboard -grid-col-3">
      <div class="-left-top -f-cosmos-overline-0">BUILD LOG</div>
      <div class="-center-top -f-cosmos-overline-0">STACK</div>

      <div class="-left dashboard__card -log">
        <IconItem :iconType="'check'" :itemText="'Dependencies installed'" />        
        <IconItem :iconType="'check'" :itemText="'Source code scaffolded'" />        
        <IconItem :iconType="'check'" :itemText="'Build complete'" />        
        <IconItem :iconType="'check'" :itemText="'Blockchain initialized'" />        
        <IconItem :iconType="'check'" :itemText="'Accounts created'" />        
        <IconItem 
          :iconType="'check'"
          :itemText="'Blockchain node started'"
          :isActive="backendRunningStates.api"          
        />        
      </div>

      <div 
        v-for="(card, index) in stack"
        :key="card.id+index"
        :class="['dashboard__card', `-${card.id}`, {'-is-active': backendRunningStates[card.id]}]"
      >
        <div class="dashboard__card-logo">
          <LogoCosmosSdk v-if="card.id === 'rpc'" />
          <LogoTendermint v-else-if="card.id === 'api'" />
          <LogoStarport v-else />
        </div>
        <div class="dashboard__card-main">
          <span class="dashboard__card-heading">{{card.name}}</span>
          <p class="dashboard__card-blurb">{{card.blurb}}</p>
          <IconItem 
            :isActive="backendRunningStates[card.id]"
            :itemText="`localhost: ${card.port}`"
            :toInjectSlot="card.id === 'frontend'"
          >
            <p v-if="card.id === 'frontend'" class="item__main">
              <a class="-with-arrow" :href="appEnv.FRONTEND" target="_blank">localhost: {{card.port}}</a>
            </p>
          </IconItem>
        </div>
        <BlockInfoCard 
          v-if="card.id === 'api'"
          :blockCards="blockCards"
        />
      </div>

    </div>

    <div class="-grid-col-3 intro">
      <div class="intro__side">
        <p class="-f-cosmos-overline-0">Architecture</p>
        <h3>Brief intro</h3>
      </div>
      <div class="intro__main">
        <p>Your blockchain is built with 
          <a href="https://github.com/cosmos/cosmos-sdk" target="_blank">Cosmos SDK</a>
          , a modular framework for building blockchains. Every feature in the Cosmos SDK is packaged as a separate module that can interact with other modules. We've installed the 
          <span>auth</span>, <span>bank</span>, and <span>staking</span> modules for you. 
          We've also generated an empty module, which you can use to start developing your own application features.          
        </p>
      </div>
    </div>
    
    <div class="tutorials">
      <div class="tutorials__top">
        <h3>Build your app</h3>
      </div>

      <div class="tutorials__articles">
        <div class="-grid-col-3 cards">
          <a 
            v-for="card in articles"
            :key="card.title"
            :class="['card-wrapper text-card', { '-is-dark': card.tagline === 'tutorial' }]"
            :href="card.link"
            target="_blank"
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

      <div class="tutorials__videos -grid-col-3">
        <a 
          v-for="card in videos"
          :key="card.title"
          class="image-card"
          :href="card.link"
          target="_blank"
        >
          <img class="image-card__img card-wrapper" :src="card.imgUrl" :alt="card.alt">
          <div class="image-card__text">
            <span class="image-card__text-h1">{{card.title}}</span>
            <span class="image-card__text-p">{{card.length}}</span>
          </div>          
        </a>
      </div>
    </div>

    <div class="footer">
      <div class="footer__main -grid-col-3">
        <div 
          v-for="block in footerBlocks"
          :key="block.title"
          class="footer__main-item"
        >
          <span>{{block.title}}</span>
          <a :href="block.link.url">{{block.link.text}}</a>
        </div>
      </div>
      <div class="footer__sub -grid-col-3">
        <div class="footer__sub-item -logo -first">
          <span><LogoTendermint/></span> Built by Tendermint Inc.
        </div>
        <div class="footer__sub-item -second">
          © Starport 2020
        </div>
        <div class="footer__sub-item -third">
          <a href="https://github.com/tendermint/starport"><LogoGithub/></a>
        </div>        
      </div>
    </div>

  </div>
  
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import moment from 'moment'

import BlockInfoCard from '@/modules/BlockInfoCard'
import IconItem from '@/components/list/IconItem'
import LogoTendermint from "@/assets/logos/LogoTendermint"
import LogoCosmosSdk from "@/assets/logos/LogoCosmosSdk"
import LogoStarport from "@/assets/logos/LogoStarport"
import LogoGithub from "@/assets/logos/LogoGithub"

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
    link: 'https://tutorials.cosmos.network/voter/'
  },
  {
    tagline: 'documentation',
    title: 'Starport Handbook',
    blurb: 'Create your own blockchain - from star to ecosystem',
    link: 'https://github.com/tendermint/starport/tree/develop/docs'
  },
  {
    tagline: 'tutorial',
    title: 'Build a Blog',
    blurb: 'Learn how Starport works by building a blog.',
    link: 'https://tutorials.cosmos.network/blog/tutorial/01-index.html'
  },
]

const videos = [
  {
    title: 'Cosmos Code With Us - Building your first Cosmos app',
    length: '1:39:07',
    imgUrl: 'https://i.ytimg.com/vi/h6Ur_40LB9k/hq720.jpg',
    alt: 'Cosmos Code With Us - Building your first Cosmos app',
    link: 'https://www.youtube.com/watch?v=h6Ur_40LB9k'
  },
  {
    title: 'Getting started with Starport, the easiest way to build a Cosmos SDK blockchain',
    length: '3:31',
    imgUrl: 'https://i.ytimg.com/vi/rmbPjCGDXek/hq720.jpg',
    alt: 'Getting started with Starport, the easiest way to build a Cosmos SDK blockchain',
    link: 'https://www.youtube.com/watch?v=rmbPjCGDXek'
  },
  {
    title: 'Tendermint Workshop: A 5 minute Blockchain Using Starport',
    length: '56:28',
    imgUrl: '/images/brian-workshop.png',
    alt: 'Tendermint Workshop: A 5 minute Blockchain Using Starport',
    link: 'https://www.youtube.com/watch?v=PGLAW-HrzWg&t=10s'
  },
]

const footerBlocks = [
  { title: 'Chat with developers', link: { text: 'Cosmos Discord', url: 'https://discord.gg/W8trcGV' } },
  { title: 'Join the community', link: { text: 'Cosmos SDK Forum', url: 'https://forum.cosmos.network' } },
  { title: 'Found an issue?', link: { text: 'Suggest improvements', url: 'https://github.com/tendermint/starport/issues' } },
]

export default {
  name: 'Welcome',
  components: {
    IconItem,
    LogoTendermint,
    LogoCosmosSdk, 
    LogoStarport,
    LogoGithub,
    BlockInfoCard
  },
  data() {
    return {
      stack,
      articles,
      videos,
      footerBlocks,
      blockCards: []
    }
  },
  computed: {
    ...mapGetters('cosmos', [ 'backendRunningStates', 'backendEnv', 'appEnv' ]),   
    ...mapGetters('cosmos/blocks', [ 'latestBlock', 'blockByHeight' ])
  },    
  methods: {
    insertBlockToStack(index, block) {
      this.blockCards.splice(index, 0, this.getFmtBlock(block))
    },
    setInitBlockCards() {
      const latestBlock = this.latestBlock

      if (latestBlock) {
        this.insertBlockToStack(0, latestBlock)

        for (let i=1; i<=2; i++) {
          if (parseInt(latestBlock.height)-i>0) {
            this.insertBlockToStack(i, this.blockByHeight(parseInt(latestBlock.height)-i)[0])
          } else {
            break
          }
        }
      }      
    },
    getFmtBlock(block) {
      return {
        height: block.height,
        hash: block.blockMeta.block_id.hash,
        time: this.getFmtBlockTime(block.blockMeta.block.header.time),        
      }
    },
    getFmtBlockTime(time) {
      return !time ? '_' : moment(time).format('H:mm:ss')
    },
  },
  watch: {
    latestBlock() {
      if (this.blockCards.length===0) {
        this.setInitBlockCards()
        return 
      }

      if (this.blockCards.length>2) this.blockCards.splice(0, 1)
      this.insertBlockToStack(2, this.latestBlock)
    },
    backendRunningStates() {
      if (!this.backendRunningStates.api) {
        this.blockCards=[]
      }
    }
  },
  created() {
    this.setInitBlockCards()
  }
}
</script>

<style scoped>

.-grid-col-3 {
  display: grid;  
  grid-column-gap: 2rem;      
  grid-template-columns: repeat(3, minmax(0, 1fr));
}
@media screen and (max-width: 992px) {
  .-grid-col-3 {
    grid-column-gap: 1rem;     
  }
}

.main {
  position: relative;
  margin-top: 4rem;
  margin-bottom: 6rem;
}
@media screen and (max-width: 768px) {
  .main {
    margin-top: 2rem;
  }
}

.hero {
  margin-bottom: 4rem;
  margin-right: 35%;
}
.hero h2 {
  font-size: 5.625rem;
  font-weight: var(--f-w-extra-bold);
  line-height: 112%;  
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
@media screen and (max-width: 1400px) {
  .hero {
    margin-right: 20%;
  }
}
@media screen and (max-width: 768px) {
  .hero {
    margin-right: 0;
  }
}
@media screen and (max-width: 576px) {
  .hero h2 {
    font-size: 16vw;
  }
}

.dashboard .-left-top { grid-area: left-top; }
.dashboard .-center-top { grid-area: center-top; }
.dashboard .-left { grid-area: left; }
.dashboard .-api { grid-area: api; }
.dashboard .-rpc { grid-area: rpc; }
.dashboard .-frontend { grid-area: frontend; }

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
.dashboard__card:not(.-log) {
  opacity: .6;
  transition: all .3s ease-in-out;  
}
.dashboard__card.-is-active {
  opacity: 1;
  transition: all .3s ease-in-out;
}
.dashboard__card.-api {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.dashboard__card.-frontend.-is-active {
  background-color: #FDFDFD;
  box-shadow: 0px 8px 40px rgba(0, 3, 66, 0.08);
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

.dashboard__card.-log {
  background-color: transparent;
  border: 1px solid rgba(0,13,158, 7%);
}
.dashboard__card.-log > div:not(:last-child) {
  margin-bottom: 1rem;
}

@media screen and (max-width: 768px) {
  .dashboard {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin-bottom: 5rem;        
  }
  .dashboard {
    grid-template-areas: 
      'center-top center-top'
      'api rpc'
      'api frontend'
      'left-top left-top'
      'left left';      
  }
  .dashboard .-left-top {
    margin-top: 4rem;
  }
}
@media screen and (max-width: 576px) {
  .dashboard {
    grid-template-areas: 
      'center-top center-top'
      'api api'
      'rpc rpc'
      'frontend frontend'
      'left-top left-top'
      'left left';  
  }
  .dashboard__card.-api {
    height: 272px;
    margin-bottom: 3rem;
  }  
  .dashboard__card-main {
    /* margin-bottom: 5rem; */
  }
}


.intro__main { grid-area: intro-main; }
.intro__side { grid-area: intro-side; }
.intro {
  grid-template-areas: 
    'intro-side intro-main intro-main';  
}
.intro {
  margin-bottom: 10rem;
}
.intro__main {
  width: 80%;
}
.intro__side p {
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
.intro__main p a {
  color: var(--c-txt-highlight);
}
.intro__main span {
  font-family: var(--f-secondary);
}
@media screen and (max-width: 768px) {
  .intro {
    grid-template-areas: 
      'intro-side'
      'intro-main';
    grid-template-columns: repeat(1, minmax(0, 1fr));
    margin-bottom: 5rem;        
  }
  .intro__side {
    margin-bottom: 1rem;
  }
  .intro__main {
    width: 100%;
  }
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
.tutorials__videos {
  margin-bottom: 8rem;
}

.card-wrapper {
  border-radius: 12px;
  box-shadow: 0px 0px 1px rgba(0, 0, 0, 0.07), 0px 8px 16px rgba(0, 0, 0, 0.05), 0px 20px 44px rgba(0, 3, 66, 0.12);    
  transition: box-shadow .25s ease-out,transform .25s ease-out,opacity .4s ease-out;  
}
.card-wrapper:hover {
  box-shadow: 0px 0px 1px rgba(0, 0, 0, 0.07), 0px 12px 24px rgba(0, 0, 0, 0.02), 0px 30px 66px rgba(0, 3, 66, 0.14);
  transform: translateY(-2px);
  transition: box-shadow .25s ease-out,transform .25s ease-out,opacity .4s ease-out;  
  transition-duration: .1s; 
}

.text-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;  
  padding: 1.5rem;
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
  color: #616489;
  margin-bottom: 4px;
}
.text-card.-is-dark .text-card__tagline {
  color: #CFD1E7;
}
.text-card__title {
  font-weight: var(--f-w-bold);
  font-size: 1.3125rem;
  line-height: 128.7%;
  letter-spacing: -0.007em;
}
.text-card__btm p {
  line-height: 130%;
  color: #616489;
  width: 90%;
}
.text-card.-is-dark .text-card__btm p {
  color: #CFD1E7;
}
@media screen and (max-width: 768px) {
  .text-card__top {
    margin-bottom: 3.5rem;
  }
  .text-card__btm p {
    font-size: 0.8125rem;
    width: 100%;
  }
}
@media screen and (max-width: 576px) {
  .cards {
    grid-template-columns: repeat(1, minmax(0, 1fr));
  }
  .cards .card-wrapper:not(:last-child) {
    margin-bottom: 1.5rem;
  }
}

/* .tutorials__videos {
  display: flex;
  justify-content: space-between;
}
.tutorials__videos .image-card {
  width: calc((100% - 2rem) / 2);
}
@media screen and (max-width: 768px) {
  .tutorials__videos .image-card {
    width: calc((100% - 1rem) / 2);
  }
} */
@media screen and (max-width: 576px) {
  .tutorials__videos {
    grid-template-columns: repeat(1, minmax(0, 1fr));
  }
  .image-card__text {
    margin-left: 2px;
  }
  .tutorials__videos {
    margin-bottom: 3rem;
  }
  .tutorials__videos .image-card:not(:last-child) {
    margin-bottom: 1.5rem;
  }  
}


.image-card {
  text-decoration: none;
}
.image-card img {
  width: 100%;
  object-fit: cover;
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
.image-card:hover .card-wrapper {
  box-shadow: 0px 0px 1px rgba(0, 0, 0, 0.07), 0px 12px 24px rgba(0, 0, 0, 0.02), 0px 30px 66px rgba(0, 3, 66, 0.14);
  transform: translateY(-2px);
  transition-duration: .1s;   
}


.footer {
  padding-top: 3rem;
  border-top: 1px solid rgba(0, 11, 119, 0.185);
  width: 100%;
}
.footer__main,
.footer__sub {
  width: 100%;
}
.footer__main {
  margin-bottom: 5rem;
}
.footer__main-item span {
  display: block;
  font-weight: var(--f-w-bold);
  line-height: 130%;
  margin-bottom: 0.5rem;
}
.footer__main-item a {
  position: relative;
  text-decoration: none;
  font-size: 16px;
  letter-spacing: -0.007em;
  font-weight: var(--f-w-medium);
  color: var(--c-txt-highlight);
}
.footer__main-item a:after {
  content: '→';  
  position: absolute;
  top: 1px;
  right: -20px;
}
.footer__main-item a:hover:after {
  right: -24px;
  transition: right .3s;
}
.footer__sub-item {
  display: flex;
  align-items: center;  
  font-weight: var(--f-w-medium);
  font-size: 0.7rem;
  letter-spacing: -0.007em;
  color: #989BB9;
}
.footer__sub-item span {
  display: inline-block;
  margin-right: 0.5rem;
}
.footer__sub-item.-logo span {
  transform: translate3d(0, 4px, 0);
}
.footer__sub-item:last-child {
  justify-content: flex-end;
}
.footer__sub-item a:hover svg >>> path {
  fill: #616489;
  transition: fill .3s;
}
@media screen and (max-width: 768px) {
  .footer__main {
    grid-template-columns: repeat(1, minmax(0, 1fr));
    grid-row-gap: 2rem;
  }
  .footer__sub .-first { grid-area: first; }
  .footer__sub .-second { grid-area: second; }
  .footer__sub .-third { grid-area: third; }
  .footer__sub {
    grid-template-columns: repeat(5, minmax(0, 1fr));
    grid-row-gap: 2rem;
    grid-template-areas: 'first first second second third';
  }  
}
@media screen and (max-width: 576px) {
  .footer__sub {
    grid-template-columns: repeat(6, minmax(0, 1fr));
    grid-template-areas: 'first first first second second third';
  }
}
@media screen and (max-width: 376px) {
  .footer__main {
    margin-bottom: 3rem;
  }
  .footer__sub {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    grid-template-areas: 'first first' 'second third';
    grid-row-gap: 0.5rem;
  }
  .footer__sub-item.-logo >>> svg {
    width: 24px;
    height: 24px;
  }
  .footer__sub-item:last-child >>> svg {
    width: 24px;
    height: 24px;
  }
}

</style>