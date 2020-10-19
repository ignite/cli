<template>
  <transition name="fade" mode="out-in" key="default">
    <div v-if="!isBlocksStackEmpty && isBackendAlive" class="explorer">
      <FullWidthContainer>
        <div slot="sideSheet" class="explorer__block">
          <transition name="fadeMild" mode="out-in">
            <BlockDetailSheet v-if="highlightedBlock" :block="highlightedBlock" :key="blockSheetKey" />
          </transition>
        </div>
        <div slot="mainContent" class="explorer__chain">
          <div class="explorer__chain-header">Blocks</div>
          <div class="explorer__chain-main">
            <BlockChain :blocks="fmtBlockData" />
          </div>
        </div>      
      </FullWidthContainer>
    </div>

    <div v-else class="explorer -is-empty" key="empty">
      <div class="explorer__container">
        <IconBox/>
        <p>Generating blocks</p>
      </div>    
    </div>
  </transition>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import blockHelpers from '@/mixins/blocks/helpers'

import _ from 'lodash'

import FullWidthContainer from '@/components/containers/FullWidthContainer'
import BlockDetailSheet from '@/modules/BlockDetailSheet'
import BlockChain from '@/modules/BlockChain'
import IconBox from "@/assets/icons/Box.vue"

export default {
  components: {
    FullWidthContainer,
    BlockChain,
    BlockDetailSheet,
    IconBox
  },
  data() {
    return {
      blockFormatter: blockHelpers.blockFormatter(),
      states: {
        isHidingBlocksWithoutTxs: false,
        isScrolledInTopHalf: true,
        isLoading: false
      },
    }
  },
  computed: {
    /*
     *
     * Vuex 
     *
     */
    ...mapGetters('cosmos', [ 'appEnv', 'backendRunningStates' ]),
    ...mapGetters('cosmos/blocks', [ 'highlightedBlock', 'blocksStack', 'lastBlock', 'stackChainRange', 'latestBlock' ]),
    /*
     *
     * Local
     * 
     */    
    blockSheetKey() {
      if (this.highlightedBlock?.data) {
        return this.highlightedBlock.data.blockMsg.blockHash
      }
      return ''
    },
    fmtBlockData() {
      const fmtBlockForTable = this.blockFormatter.blockForTable(this.blocksStack)

      if (!fmtBlockForTable) return null

      if (this.states.isHidingBlocksWithoutTxs) {
        return this.blockFormatter.filterBlock(fmtBlockForTable).hideBlocksWithoutTxs()
      }

      return fmtBlockForTable
    },
    isBlocksStackEmpty() {
      return this.blocksStack.length<=0 ||
        !this.fmtBlockData || 
        this.fmtBlockData?.length<=0
    },
    isBackendAlive() {
      return this.backendRunningStates.api
    }
  },  
  methods: {
    /*
     *
     * Vuex 
     *
     */    
    ...mapMutations('cosmos/blocks', [ 'popOverloadBlocks', 'sortBlocksStack' ]),
    ...mapActions('cosmos/blocks', [ 'addBlockEntry', 'getBlockchain', 'setHighlightedBlock' ]),
    /*
     *
     * Local 
     *
     */      
    /*
     *
     // Pop overloaded blocks (over maxStackCount)
     // only when scrolling to upperhalf of the table
     *
     */         
    async handleScrollTop() {
      this.states.isScrolledInTopHalf=true      
      if (!this.latestBlock) return 
      
      const isShowingLatestBlock = (parseInt(this.latestBlock.height) === this.stackChainRange.highestHeight)

      if (!isShowingLatestBlock && !this.states.isLoading) {
        this.states.isLoading=true

        await this.getBlockchain({ 
          blockHeight: this.stackChainRange.highestHeight,
          toGetLowerBlocks: false
        }).then(() => this.states.isLoading=false)        
      }
    },
    /*
     *
     // Load extra 20 blocks
     // only when scrolling to bottom of the table
     *
     */          
    async handleScrollBottom() {
      this.states.isScrolledInTopHalf=false
      this.states.isLoading=true

      await this.getBlockchain({ 
        blockHeight: this.lastBlock.height,
        toGetLowerBlocks: true
      }).then(() => this.states.isLoading=false)
    }   
  },
  beforeDestroy() {
    if (this.latestBlock) {
      this.getBlockchain({ 
        blockHeight: this.latestBlock.height,
        toReset: true,
        toGetLowerBlocks: true
      })
    }
  }
}
</script>

<style scoped>

.explorer {
  --top-offset: 2.25rem;
}
.explorer {
  height: calc(100vh - var(--header-height) - 1px - 2.25rem);
  padding-top: 2.25rem;
}
@media only screen and (max-width: 992px) {
  .explorer {
    --top-offset: 1.5rem;
  }
}

.explorer__chain {
  height: inherit;
}
.explorer__chain-main {
  height: calc(100% - 40px);
}
.explorer__chain-header {
  font-size: 3.1875rem;
  font-weight: var(--f-w-bold);
  margin-bottom: 2rem;
  padding-left: calc(var(--g-offset-side) - 4px);
}
@media only screen and (max-width: 992px) {
  .explorer__chain-header {
    margin-bottom: 1rem;
  }  
}

.explorer__block {
  height: 100%;
  min-width: 400px;
}

.explorer.-is-empty {
  display: flex;
  justify-content: center;
  align-items: center;
  color: var(--c-txt-light);
  animation: tempLoading 5s ease-in-out infinite;
}
.explorer.-is-empty .explorer__container svg {
  display: block;
  margin: 0 auto 0.5rem auto;
}
.explorer.-is-empty .explorer__container svg >>> path {
  fill: var(--c-txt-light);
  fill-opacity: 1;
}
@keyframes tempLoading {
	0%, 100% { opacity: 0.3;}
	50% { opacity: .8; }
	75% { opacity: .8; }
}


</style>