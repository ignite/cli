<template>
  <div v-if="!block.data" :class="['sheet -is-empty']">Waiting for block data</div>

  <div v-else class="sheet">
    <div class="sheet__header">
      <div class="sheet__header-main">
        <p>{{block.data.blockMsg.height}}</p>
      </div>
      <div class="sheet__header-side">
        <div class="sheet__header-side-top">
          <CopyIconText 
            :text="block.data.blockMsg.blockHash_sliced" 
            :link="`${appEnv.RPC}/block?hash=${block.data.blockMsg.blockHash}`"
          />
        </div>
        <div class="sheet__header-side-btm">
          <span>{{getFmtTime(block.data.blockMsg.time)}}</span>
        </div>
      </div>
    </div>

    <div class="sheet__main">
      <div 
        v-if="block.data.blockMsg.txs>0 && block.data.txs.length>0"
        class="txs"
      >
        <div class="txs__header">
          <p>Transactions</p>
        </div>

        <div class="txs__tx tx">
          <div class="tx__main">
            <p class="tx__title">Messages</p>
          </div>
          <div class="tx__side">
            <p class="tx__title">Tx Info</p>
          </div>
        </div>
      </div>
      <div v-else class="txs -is-empty">
        <p>0 Transactions</p>
      </div>      
    </div>
  </div>
</template>

<script>
import moment from 'moment'
import { mapGetters } from 'vuex'

import CopyIconText from '@/components/texts/CopyIconText'

export default {
  components: {
    CopyIconText
  },
  props: {
    block: { type: Object }
  },
  computed: {
    ...mapGetters('cosmos', [ 'appEnv' ]),    
  },    
  methods: {
    getFmtTime(time) {
      const momentTime = moment(time)
      return momentTime.format('MMM D YYYY, HH:mm:ss')
    } 
  }
}
</script>

<style scoped>

.sheet {
  height: 100%;
}

.sheet.-is-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  animation: tempLoadingEffect 1.5s ease-in-out infinite;
}
@keyframes tempLoadingEffect {
  0% { color: var(--c-txt-grey); }
  50% { color: var(--c-txt-secondary); }
  100% { color: var(--c-txt-grey); }
}

.sheet__header {
  display: flex;
  align-items: center;
}

.sheet__header-main {
  margin-right: 1.5rem;
}
.sheet__header-main p {
  font-size: 3.1875rem;
  font-weight: var(--f-w-bold);  
}

.sheet__header-side-top {
  margin-bottom: 4px;
}

.sheet__header-side-btm {
  margin-bottom: 4px;
}
.sheet__header-side-btm span {
  font-size: 0.8125rem;
  color: var(--c-txt-secondary);
}

.sheet__main {
  height: 100%;
}

.txs {
  height: 100%;
}
.txs.-is-empty {
  display: flex;
  align-items: center;
  justify-content: center;
}
.txs.-is-empty p {
  font-size: 3rem;
  font-weight: var(--f-w-light);
  color: var(--c-txt-gray);
  margin-bottom: 5rem;  
}

</style>