<template>
  <transition name="fade">
    <div v-if="blockCards.length>0" class="container" key="default">
      <transition-group name="list" tag="ul">
        <div 
          v-for="(block) in blockCards"
          :key="block.hash"
          class="card"
          @click="handleCardClick(block)"
        >
          <div class="card__top">
            <div class="card__top-left">
              <p>BLOCK</p>
              <p>{{block.height}}</p>
            </div>
            <div class="card__top-right">
              <span>{{block.time}}</span>
            </div>
          </div>
          <div class="card__btm">
            <p class="card__hash">{{block.hash}}</p>
          </div>
          <div class="card__bg">
            <Box/>
          </div>
        </div>
      </transition-group>
    </div>        

    <div v-else class="empty-card" key="empty">
      <div class="empty-card__container">
        <Box/>
        <p>Generating blocks</p>
      </div>
    </div>  
  </transition>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import blockHelpers from '@/mixins/blocks/helpers'

import Box from "@/assets/icons/Box.vue"

export default {
  components: {
    Box,
  },
  props: {
    blockCards: {
      type: Array,
      required: true,
      validator(value) {
        return value.filter(block => block.height && block.time && block.hash)
          .length===value.length
      }
    }
  },
  data() {
    return {
      blockFormatter: blockHelpers.blockFormatter(),
    }
  },
  computed: {
    ...mapGetters('cosmos/blocks', [ 'blockByHeight' ]),
  },
  methods: {
    /*
     *
     * Vuex 
     *
     */        
    ...mapMutations('cosmos/ui', [ 'setTableSheetState' ]),
    ...mapActions('cosmos/blocks', [ 'setHighlightedBlock' ]),
    /*
     *
     * Local 
     *
     */          
    handleCardClick({ height, hash }) {
      const blockData = this.blockByHeight(height)
      const fmtBlockData = this.blockFormatter.blockForTable(blockData)[0]
      this.setHighlightedBlock({
        block: { id: hash, data: fmtBlockData }
      })   

      this.$router.push('blocks')
    }
  }
}
</script>

<style scoped>

.card {
  overflow: hidden;
  padding: 1.5rem;
  background-color: var(--c-bg-primary);
  border-radius: 12px;
}
.card__top {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4.5rem;
}
.card__top-left {
  color: var(--c-txt-highlight);
}
.card__top-left p:first-child {
  font-size: 0.75rem;
  font-weight: var(--f-w-medium);
  margin-bottom: 4px;
}
.card__top-left p:last-child {
  font-size: 1.3125rem;
  font-weight: var(--f-w-bold);
}
.card__top-right {
  font-size: 0.75rem;
  color: rgba(0, 5, 66, 0.621);
}
.card__btm {
  color: rgba(0, 5, 66, 0.621);
}
.card__bg {
  position: absolute;
  right: -24px;
  bottom: 30%;
}
.card__hash {
  display: block;
  white-space: nowrap; /* forces text to single line */
  overflow: hidden;
  text-overflow: ellipsis;  
  font-family: var(--f-secondary);
}

.container {
  /* height: 100%; */
  position: relative;
  transform: translate3d(0, 3rem, 0);
  /* box-shadow: 0px 8px 40px rgba(0, 3, 66, 0.08);   */
  perspective: 1000px;
}
.card {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  box-shadow: 0px 8px 40px rgba(0, 3, 66, 0.08);  
  transform-origin: center;
  pointer-events: none;
  transition: transform .3s ease-in-out;
  will-change: transform;
}
.card:hover {
  transform: translate3d(0, 2px, 8px);
  transition: transform .3s ease-in-out;
  will-change: transform;
}
.card:nth-last-child(1) {
  z-index: 0;
  pointer-events: initial;
}
.card:nth-last-child(1):hover {
  cursor: pointer;
}
.card:nth-last-child(2) {
  transform: translate3d(0,-24px,-50px);
  z-index: -1;
  color: transparent;
  transition: transform .5s;
}
.card:nth-last-child(3) {
  transform: translate3d(0,-48px,-100px);
  z-index: -2;
  color: transparent;
  transition: transform .5s;
}
.card:nth-last-child(2) *,
.card:nth-last-child(3) * {
  color: transparent;
  transition: color .5s;
}
.card:nth-last-child(4) {
  transform: translate3d(0,-60px,-150px);
  z-index: -3;
  transition: transform .5s;
  opacity: 0;
}
@media screen and (max-width: 576px) {
  .card__top {
    margin-bottom: 3rem;
  }
}

.empty-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;  
  padding: 1.5rem 0;  
  border: 1px solid rgba(0, 13, 158, 0.07);
  border-radius: 12px;  

  animation: tempLoading 5s ease-in-out infinite;
}
.empty-card__container svg {
  display: block;
  margin: 0 auto 1rem auto;
}
.empty-card__container p {
  font-size: 0.75rem;
  line-height: 130.9%;
  letter-spacing: 0.005em;
  color: var(--c-txt-grey);
  opacity: .8;
}

@keyframes tempLoading {
	0%, 100% { opacity: 0.3;}
	50% { opacity: 1; }
	75% { opacity: 1; }
}
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translate3d(0, 48px, 0);
  }
  to {
    opacity: 1;
    transform: translate3d(0, 0, 0);
  }
}
@keyframes cardSlideIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
@keyframes containerSlideIn {
  from {
    opacity: 0;
    transform: translate3d(0, 4rem, 0);
  }
  to {
    opacity: 1;
    transform: translate3d(0, 3rem, 0);
  }
}
@keyframes firstCardSlideIn {
  from { 
    opacity: 0;
    transform: translate3d(0,16px,0); 
  }
  to { 
    opacity: 1;
    transform: translate3d(0,0,0); 
  }
}
@keyframes secondCardSlideIn {
  from { 
    opacity: 0;
    transform: translate3d(0,0px,-50px); 
  }
  to { 
    opacity: 1;
    transform: translate3d(0,-24px,-50px); 
  }
}
@keyframes thirdCardSlideIn {
  from { 
    opacity: 0;
    transform: translate3d(0,-24px,-100px); 
  }
  to { 
    opacity: 1;
    transform: translate3d(0,-48px,-100px); 
  }
}

.list-enter-active {
  animation: slideIn 1s;
}

.fade-leave-active {
  animation: none;  
  opacity: 0;
  transition: opacity .5s;
}
.container.fade-enter-active {
  transition: 2s;
  animation: cardSlideIn 1.5s ease-in-out;
}
.container.fade-enter-active .card:nth-last-child(1) {
  animation: firstCardSlideIn .5s ease-out;
}
.container.fade-enter-active .card:nth-last-child(2) {
  animation: secondCardSlideIn .5s ease-out;
  animation-delay: .25s;
}
.container.fade-enter-active .card:nth-last-child(3) {
  animation: thirdCardSlideIn .5s ease-out;
  animation-delay: .5s;
}

</style>
