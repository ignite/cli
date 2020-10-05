<template>
  <div class="container">
    <transition-group name="list" tag="ul">
      <div 
        v-for="(block) in blockCards"
        :key="block.hash"
        class="card"
        @click="handleRowClick(block)"
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
          <p ref="blockHash" class="card__hash">{{block.hash}}</p>
        </div>
        <div class="card__bg">
          <Box/>
        </div>
      </div>
    </transition-group>
  </div>        
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
        return value
          .filter(block => block.height && block.time && block.hash)
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
    ...mapGetters('cosmos/ui', [ 'blocksExplorerTableId' ]),
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
    handleRowClick({ height, hash }) {
      const blockData = this.blockByHeight(height)
      const fmtBlockData = this.blockFormatter.blockForTable(blockData)[0]
      this.setHighlightedBlock({
        block: { id: hash, data: fmtBlockData }
      })
      this.setTableSheetState({
        tableId: this.blocksExplorerTableId,
        sheetState: true
      })      

      this.$router.push('blocks')
    }
  }
}
</script>

<style scoped>

.card {
  position: relative;
  padding: 1.5rem;
  background-color: var(--c-bg-primary);
  border-radius: 12px;
}
.card__top {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5rem;
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
  right: 0;
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
  height: 100%;
  position: relative;
  transform: translate3d(0, 4rem, 0);
  /* box-shadow: 0px 8px 40px rgba(0, 3, 66, 0.08);   */
  perspective: 1000px;
}
.card {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 100%;
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
  transform: translate3d(0,-20px,-50px);
  z-index: -1;
  transition: transform .5s;
}
.card:nth-last-child(3) {
  transform: translate3d(0,-40px,-100px);
  z-index: -2;
  transition: transform .5s;
}
.card:nth-last-child(4) {
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

</style>