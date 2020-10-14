<template>
  <div class="chain">
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
            <span class="block-info__text">{{block.txs.length}} transactions</span>
            ·
            <span class="block-info__text">{{getMsgsAmount(block.txs)}} messages</span>
          </div>
        </BlockCard>   
      </button>
    </div>
  </div>  
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex'
import moment from 'moment'

import blockHelpers from '@/mixins/blocks/helpers'

import BlockCard from '@/components/cards/BlockCard'

export default {
  props: {
    blocks: { type: Array }
  },
  components: {
    BlockCard,    
  },
  data() {
    return {
      blockFormatter: blockHelpers.blockFormatter(),
      states: {
        isScrolledInTopHalf: true,
        isLoading: false
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
    handleCardClicked(event) {
      const blockHash = event.currentTarget.id
      const blockPayload = {
        id: blockHash,
        data: this.blockByHash(blockHash)[0]
      }
      
      this.setHighlightedBlock({ block: blockPayload })           
      this.localHighlightedBlock=blockPayload
    }  
  },
  watch: {
    latestBlock() {
      if (!this.localHighlightedBlock) {
        this.setHighlightedBlock({ block: {
          id: this.latestBlock.blockMeta.block_id.hash,
          data: this.latestBlock
        }})
      }
    }
  },  
}
</script>

<style scoped>

.chain {
  --bg-offset: 0.5rem;
}

.chain {
  min-height: 100%;
  max-height: 100%;
  overflow-y: scroll;
  overflow-x: visible;
}
.chain::-webkit-scrollbar {
  opacity: 0;
}

.chain__block {
  position: relative;
  width: 100%;
  padding-left: calc(var(--g-offset-side) - var(--bg-offset));
  padding-right: var(--bg-offset);
}
.chain__block:after {
  content: '';
  position: absolute;
  bottom: 0;
  left: calc(var(--g-offset-side) - var(--bg-offset)*2);
  width: calc(100% - var(--g-offset-side));
  height: 1px;
  background-color: var(--c-border-primary);  
}
.chain__block:before {
  content: '';
  position: absolute;
  z-index: -1;
  top: 0;
  left: calc(var(--g-offset-side) - var(--bg-offset) * 4);
  width: calc(100% - var(--g-offset-side)/2 + 8px);
  height: 100%;
  border-radius: 16px;
  background-color: var(--c-bg-secondary);
  opacity: 0;
  transition: opacity .6s ease-in;
}
.chain__block.-is-active:before {
  opacity: 1;
  transition: opacity .5s ease-in-out;
}

.chain__block >>> .card {
  margin-right: 1rem;
}
.chain__block.-has-txs >>> .card__title {
  color: var(--c-txt-highlight);
}

.block-info__text:first-child {
  font-weight: var(--f-w-medium);
}
.block-info__text:last-child {
  color: var(--c-txt-secondary);
}

</style>