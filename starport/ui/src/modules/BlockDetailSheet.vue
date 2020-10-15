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
          <p class="txs__header-title">Transactions</p>
          <p class="txs__header-note">
            <span>{{block.data.txs.length}}</span>
            <span v-if="failedTxsCount"> / </span>
            <span v-if="failedTxsCount" class="txs__header-note-warn">{{ failedTxsCount }} error</span>
          </p>
        </div>

        <div 
          v-for="(tx, txIndex) in block.data.txs"
          :key="txIndex+tx.txhash"
          class="txs__tx tx"
        >
          <div class="tx__main">
            <div v-if="tx.code" class="tx__error">
              <span class="tx__error-title">Error</span>
              <p class="tx__error-msg">{{ tx.raw_log }}</p>
            </div>

            <TxMsgCards :msgs="tx.tx.value.msg" />
          </div>
          <div class="tx__side">
            <div class="tx__info">
              <p class="tx__title">Tx Info</p>

              <div class="tx__info-container">
                <div class="tx__info-content tx-info">
                  <span class="tx-info__title">Hash</span>
                  <CopyIconText 
                    :text="tx.txhash" 
                    :link="`${appEnv.RPC}/block?hash=${tx.txhash}`"
                  />                  
                </div>
                <div class="tx__info-content tx-info">
                  <span class="tx-info__title">Gas Used / Wanted</span>
                  <p class="tx-info__content">{{ `${tx.gas_used} / ${tx.gas_wanted}` }}</p>
                </div>
                <div class="tx__info-content tx-info">
                  <span class="tx-info__title">Fee</span>
                  <p class="tx-info__content">{{ getTxFee(tx) }}</p>
                </div>                
              </div>
            </div>
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
import { uuid } from 'vue-uuid'
import { mapGetters } from 'vuex'
import { getters } from '@/mixins/helpers'

import CopyIconText from '@/components/texts/CopyIconText'
import TxMsgCards from '@/modules/TxMsgCards'

export default {
  components: {
    CopyIconText,
    TxMsgCards
  },
  props: {
    block: { type: Object }
  },
  computed: {
    ...mapGetters('cosmos', [ 'appEnv' ]),
    failedTxsCount() {
      return this.block.data.txs.filter(tx => tx.code).length
    }    
  },    
  methods: {
    getFmtTime(time) {
      const momentTime = moment(time)
      return momentTime.format('MMM D YYYY, HH:mm:ss')
    },
    getTxFee(tx) {
      const fee = tx.tx.value.fee.amount[0]
      return `${fee.amount} ${fee.denom}`
    }
  }
}
</script>

<style scoped>

.sheet {
  height: 100%;
  padding-right: var(--g-offset-side);
}
.sheet {
  overflow-y: scroll;
}
.sheet::-webkit-scrollbar {
  width: 0px;
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
  margin-bottom: 2.5rem;
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
  font-size: 2rem;
  font-weight: var(--f-w-light);
  color: var(--c-txt-gray);
  margin-bottom: 5rem;  
}

.txs__header {
  display: flex;
  align-items: flex-end;
  margin-left: 2px;
  margin-bottom: 1.5rem;
}
.txs__header-title {
  font-size: 1.3125rem;
  font-weight: var(--f-w-medium);
  margin-right: 0.85rem;
}
.txs__header-note {
  font-size: 1rem;
  color: var(--c-txt-third);
  margin-bottom: 1.8px;
}
.txs__header-note-warn {
  color: var(--c-txt-danger);
}

.tx {
  display: flex;
  margin-bottom: 3rem;  
}
.tx:not(:last-child) {
  padding-bottom: 3rem;
  border-bottom: 1px solid var(--c-border-primary);
}
.tx__main {
  flex-grow: 1;
  margin-right: 3rem;
}
.tx__side {
  width: 15vw;
  max-width: 180px;
}

.tx__error {
  color: var(--c-txt-danger);  
  padding: 1.25rem 1.5rem;
  border-radius: 12px;
  background-color: var(--c-danger-light);
  margin-bottom: 1.5rem;
}
.tx__error-title {
  display: block;
  font-size: 0.75rem;    
  font-weight: var(--f-w-bold);
  text-transform: uppercase;  
  margin-bottom: 0.5rem;
}
.tx__error-msg {
  font-size: 0.875rem;  
}

.tx__title {
  font-weight: var(--f-w-medium);
  font-size: 0.75rem;
  line-height: 130.9%;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--c-txt-third);
  margin-bottom: 0.85rem;
}

.tx-info:not(:last-child) {
  margin-bottom: 1.5rem;
}
.tx-info__title {
  display: inline-block;
  /* font-weight: var(--f-w-medium); */
  font-size: 0.75rem;
  line-height: 130.9%;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--c-txt-third);
  margin-bottom: 4px;
}
.tx-info__content {
  color: var(--c-txt-secondary);
}

</style>

