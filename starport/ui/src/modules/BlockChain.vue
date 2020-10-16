<template>
  <div
    :class="[
      'chain-container', {
        '-is-loading': isLoading,
        '-has-higher': hasHigherBlocks,
        '-has-lower': hasLowerBlocks
      }
    ]"  
  >
    <div
      :class="['chain']"
      ref="chain"
      @scroll="[handleTableScroll($event), updateScrollValue()]"
    >  
      <div v-if="blocks" class="chain__blocks">
        <button 
          :class="['chain__block', {
            '-has-txs': block.txs.length>0,
            '-is-active': block.blockMsg.blockHash === highlightedBlock.id
          }]"
          v-for="block in blocks"
          :key="block.blockMsg.blockHash"          
          :id="block.blockMsg.blockHash"
          @click="handleCardClicked"
        >
          <BlockCard
            :title="block.blockMsg.height"
            :note="getFmtTime(block.blockMsg.time)"
            :isActive="block.blockMsg.blockHash === highlightedBlock.id"   
          >
            <div v-if="block.txs.length>0" class="block-info">
              <span v-if="getFailedTxsCount(block.txs)>0" class="block-info__indicator"></span>
              <span class="block-info__text">{{block.txs.length}} transactions</span>
              ·
              <span class="block-info__text">{{getMsgsAmount(block.txs)}} messages</span>
            </div>
          </BlockCard>   
        </button>   
      </div>     
    </div>  

    <button class="util-btn -top" @click="handleNavClick('top')"><IconArrow/></button>      
    <!-- <button class="util-btn -btm" @click="handleNavClick('btm')"><IconArrow/></button>         -->
  </div>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import moment from 'moment'

import blockHelpers from '@/mixins/blocks/helpers'

import BlockCard from '@/components/cards/BlockCard'
import IconArrow from '@/assets/icons/Arrow'

export default {
  props: {
    blocks: { type: Array }
  },
  components: {
    BlockCard,    
    IconArrow
  },
  data() {
    return {
      blockFormatter: blockHelpers.blockFormatter(),
      states: {
        lastScrolledHeight: 0,
        lastScrolledTop: 0,
        isScrolledTop: false,
        isScrolledBottom: false,
        isScrolledAwayFromTop: false,
        isLoading: false,
        hasHigherBlocks: false,
        hasLowerBlocks: false
      },
      localHighlightedBlock: null
    }
  },  
  computed: {
    /*
     *
     * Vuex 
     *
     */
    ...mapGetters('cosmos', [ 'appEnv' ]),
    ...mapGetters('cosmos/blocks', [ 'highlightedBlock', 'blocksStack', 'lastBlock', 'stackChainRange', 'latestBlock', 'blockByHash' ]),
    /*
     *
     * Local 
     *
     */   
    isLoading() { return this.states.isLoading },
    isScrolledTop() { return this.states.isScrolledTop },
    isScrolledBottom() { return this.states.isScrolledBottom },
    hasHigherBlocks() { return this.states.hasHigherBlocks },
    hasLowerBlocks() { return this.states.hasLowerBlocks },
  },
  methods: {
    /*
     *
     * Vuex 
     *
     */    
    ...mapActions('cosmos/blocks', [ 'addBlockEntry', 'getBlockchain', 'setHighlightedBlock', 'popOverloadBlocks', 'sortBlocksStack' ]),    
    /*
     *
     * Local 
     *
     */       
    getFmtTime(time) {
      const momentTime = moment(time)
      const duration = moment.duration(moment().diff(momentTime))

      if (duration.as('years')>=1) {
        return momentTime.format('MMM D YYYY, HH:mm:ss')
      } else if (duration.as('days')>=1) {
        return momentTime.format('MMM D, HH:mm:ss')
      }
      return momentTime.format('HH:mm:ss')
    },
    getMsgsAmount(txs) {
      return txs
        .map(tx => tx.tx.value.msg.length)
        .reduce((accu, curr) => accu+curr)
    },
    getFailedTxsCount(txs) {
      return txs.filter(tx => tx.code).length
    },    
    handleCardClicked(event) {
      const blockHash = event.currentTarget.id
      const blockPayload = {
        id: blockHash,
        data: this.blockFormatter.blockForTable(this.blockByHash(blockHash))[0]
      }
      
      this.setHighlightedBlock({ block: blockPayload })           
      this.localHighlightedBlock=blockPayload
    },
    handleTableScroll: _.debounce(function(event) {
      const $table = event.target
      const scrolledHeight = $table.scrollTop + $table.offsetHeight
      const tableScrollHeight = $table.scrollHeight

      const isScrolledToTop = scrolledHeight <= $table.offsetHeight
      const isScrolledToBottom = scrolledHeight + 100 >= tableScrollHeight
      const isOnTopHalf = $table.scrollTop < (tableScrollHeight-$table.offsetHeight) / 2
      const isScrollAwayFromTop = scrolledHeight / tableScrollHeight > 0.2

      const isCallableScrolledDistance = 
        $table.offsetHeight / Math.abs(scrolledHeight-this.states.lastScrolledHeight) > 25
      
      if (isCallableScrolledDistance) {
        this.isScrolledAwayFromTop = isScrollAwayFromTop

        if (isScrolledToBottom) {
          this.states.isScrolledTop=false
          this.states.isScrolledBottom=true
          this.handleScrollBottom()
        }
        if (isScrolledToTop) {
          this.states.isScrolledTop=true
          this.states.isScrolledBottom=false          
          this.handleScrollTop()
        }
      }
    }, 200),
    updateScrollValue() {
      const $table = event.target
      const scrolledHeight = $table.scrollTop + $table.offsetHeight 
      this.states.lastScrolledHeight = scrolledHeight     
    },
    /*
     *
     // Pop overloaded blocks (over maxStackCount)
     // only when scrolling to upperhalf of the table
     *
     */         
    async getHigherBlocks() {
      if (!this.latestBlock) return 
      
      const isShowingLatestBlock = (parseInt(this.latestBlock.height) === this.stackChainRange.highestHeight)

      if (!isShowingLatestBlock && !this.states.isLoading) {
        this.states.isLoading=true

        await this.getBlockchain({ 
          blockHeight: this.stackChainRange.highestHeight,
          toGetLowerBlocks: false
        }).then(() => {
          this.states.isLoading=false
          this.setHasHigherBlocksState()      
          setTimeout((() => {
            this.$refs.chain.scrollBy({
              top: 24,
              left: 0,
              behavior: 'smooth'
            })            
          }).bind(this), 100)
        })        
      }
    },    
    /*
     *
     // Load extra 20 blocks
     // only when scrolling to bottom of the table
     *
     */          
    async getLowerBlocks() {
      this.states.isLoading=true

      await this.getBlockchain({ 
        blockHeight: this.lastBlock.height,
        toGetLowerBlocks: true
      }).then(() => {
        this.states.isLoading=false
        this.setHasHigherBlocksState()  
        this.$refs.chain.scrollBy({
          bottom: 24,
          left: 0,
          behavior: 'smooth'
        })            
      })
    },    
    handleScrollTop() {
      this.getHigherBlocks()
    },
    handleScrollBottom() {
      this.getLowerBlocks()
    },
    handleNavClick(dir) {
      if (dir==='top' && this.states.hasHigherBlocks) {
        this.getHigherBlocks()
        this.$refs.chain.scrollTo(0,0)        
      }
    },
    setHasHigherBlocksState() {
      if (
        (this.stackChainRange.highestHeight !== this.latestBlock.height) ||
        this.isScrolledAwayFromTop &&
        !this.isScrolledTop
      ) {
        this.states.hasHigherBlocks=true
      } else {
        this.states.hasHigherBlocks=false
      }    
    },
    setHasLowerBlocksState() {
      if (this.stackChainRange.highestHeight !== 1) {
        this.states.hasLowerBlocks=true
      } else {
        this.states.hasLowerBlocks=false
      }      
    }
  },
  watch: {
    latestBlock() {
      /**
       * 
       // If no block is clicked (selected),
       // set highlighted block to be latest block.
       * 
       */
      if (!this.localHighlightedBlock) {
        this.setHighlightedBlock({ block: {
          id: this.latestBlock.blockMeta.block_id.hash,
          data: this.blockFormatter.blockForTable([this.latestBlock])[0]
        }})
      }

      this.setHasHigherBlocksState()      
      this.setHasLowerBlocksState()
    },
    isScrolledTop() {
      this.setHasHigherBlocksState()
    },
    isScrolledBottom() {
      this.setHasLowerBlocksState()
    }
  },  
  mounted() {
    if (this.latestBlock) {
      this.setHighlightedBlock({ block: {
        id: this.latestBlock.blockMeta.block_id.hash,
        data: this.blockFormatter.blockForTable([this.latestBlock])[0]
      }})      
    }
  }
}
</script>

<style scoped>

.chain-container {
  position: relative;
  height: inherit;
}

.chain {
  --bg-offset: 0.5rem;
}

.chain {
  position: relative;
  min-height: calc(100% - 0rem);
  max-height: calc(100% - 0rem);
  overflow-y: scroll;
  overflow-x: visible;

	background:
		linear-gradient(white 30%, rgba(255,255,255,0)),
		linear-gradient(rgba(255,255,255,0), white 70%) 0 100%,
		radial-gradient(farthest-side at 50% 0, rgba(0,0,0,.05), rgba(0,0,0,0)),
		radial-gradient(farthest-side at 50% 100%, rgba(0,0,0,.05), rgba(0,0,0,0)) 0 100%;
  background-repeat: no-repeat;
  background-size: 100% 24px, 100% 24px, 100% 24px, 100% 24px;
  background-attachment: local, local, scroll, scroll;  
}
.chain::-webkit-scrollbar {
  width: 0px;
}

.chain__block {
  position: relative;
  width: 100%;
  padding-left: var(--g-offset-side);
  padding-right: var(--bg-offset);
}
.chain__block:after {
  content: '';
  position: absolute;
  bottom: 0;
  left: calc(var(--g-offset-side) - var(--bg-offset));
  width: calc(100% - var(--g-offset-side));
  height: 1px;
  background-color: var(--c-border-primary);  
}
/* .chain__block:before {
  content: '';
  position: absolute;
  z-index: -1;
  top: 0;
  left: calc(var(--g-offset-side) - var(--bg-offset) * 3);
  width: calc(100% - var(--g-offset-side)/2);
  height: 100%;
  border-radius: 16px;
  background-color: var(--c-bg-secondary);
  opacity: 0;
  transition: opacity .6s ease-in;
} */
.chain__block:before {
  content: '';
  position: absolute;
  z-index: 2;
  top: -3px;
  left: calc(var(--g-offset-side) - 1.85rem);
  width: 4px;
  height: calc(100% + 6px);
  /* border-radius: 16px; */
  background-color: var(--c-txt-highlight);
  opacity: 0;
  transition: opacity .3s ease-in;
}
.chain__block.-is-active:before {
  opacity: 1;
  transition: opacity .3s ease-in-out;
}
.chain__block >>> .card {
  margin-right: 1rem;
}
.chain__block.-has-txs >>> .card__title {
  color: var(--c-txt-highlight);
}
@media only screen and (max-width: 992px) {
  .chain__block:before {
    left: calc(var(--g-offset-side) - 1.5rem);
  }
}

.block-info__text:first-child {
  font-weight: var(--f-w-medium);
  color: var(--c-txt-secondary);
}
.block-info__text:last-child {
  color: var(--c-txt-third);
}
.block-info__indicator {
  display: inline-block;
  width: 5px;
  height: 5px;
  border-radius: 100%;
  background-color: var(--c-danger-primary);
  margin-right: 6px;
  transform: translate3d(0, -2px, 0);
}

.util-btn {
  position: absolute;
  top: -0.8rem;
  left: calc((100% - var(--g-offset-side)) / 2 + 22px);
  width: 22px;
  height: 22px;
  background-color: var(--c-bg-primary);
  border-radius: 100%;
  box-shadow: 0px 0px 1px rgba(0, 0, 0, 0.07), 0px 0.5px 3px rgba(0, 0, 0, 0.1), 0px 1.25px 6px rgba(0, 3, 66, 0.08);  
}
.util-btn.-btm {
  top: auto;
  bottom: 1rem;
  transform: rotate(180deg);
  box-shadow: 0px 0px 1px rgba(0, 0, 0, 0.07), 0px 0.5px 3px rgba(0, 0, 0, 0.1), 0px 1.25px 6px rgba(0, 3, 66, 0.08);  
}
.util-btn {
  opacity: 0;
  pointer-events: none;
  transition: opacity .3s;
}
.chain-container.-has-higher .util-btn.-top,
.chain-container.-has-lower .util-btn.-btm {
  opacity: 1;
  pointer-events: initial;
  transition: opacity .3s;
}

</style>