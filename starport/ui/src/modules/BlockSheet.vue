<template>
  <div v-if="!blockData">waiting for block data...</div>
  
  <div 
    v-else
    :class="['sheet']"
  >
    <div class="sheet__top -container -border-btm">
      <h3 class="sheet__heading">Block #{{blockData.blockMsg.height}}</h3>
    </div>
    <div class="sheet__sub -container -border-btm">
      <ListWrapper :listItems="[
        { headText: 'Hash', contentText: blockData.blockMsg.blockHash },
        { headText: 'Time', contentText: blockData.blockMsg.time },
      ]" />      
    </div>
    <div class="sheet__main -container">
      <div class="cards-container">
        <div class="cards-container__top">
          <h4 class="cards-container__label">Transactions</h4>
        </div>

        <!-- transactions -->
        <div v-if="blockData.blockMsg.txs>0 && blockData.txs.length>0">
          <div 
            v-for="tx in messagesForTable"
            :key="tx.tableData.id"          
            class="cards-container__card"
          >
            <TxCard :txData="tx" />
          </div>
        </div>
        <div v-else-if="blockData.blockMsg.txs>0 && blockData.txs.length<=0">🚨 Error fetching transaction data</div>
        <div v-else>No transactions are included</div>

      </div>
    </div>
  </div>
</template>

<script>
import SideTabList from '@/components/table/SideTabList'
import ListWrapper from '@/components/list/ListWrapper'
import TxCard from '@/modules/TxCard'

export default {
  props: {
    blockData: { type: Object }
  },
  components: {
    SideTabList,
    ListWrapper,
    TxCard
  },
  data() {
    return {
      isActive: false
    }
  },
  computed: {
    messagesForTable() {
      return this.blockData.txs.map(item => {
        const {
          fee,
          msg,
          memo
        } = item

        return {
          txMsg: {
            hash: 'faketransactionhashfornow', // temp
            status: 'Fakestatus', // temp
            fee: fee.amount[0].amount, // temp
            gas: fee.gas, // temp
            memo: memo && memo.length>0 ? memo : 'N/A'
          },
          msgs: msg.map(({
            type,
            value
          }) => ({
            type: this.getMsgType(type),
            amount: this.getAmount(value.amount),
            delegator: value.delegator_address,
            validator: value.validator_address,
            from: value.from_address,
            to: value.to_address
          })),
          tableData: {
            id: item.signatures[0].signature, // temp
            isActive: false
          },
        }
      })
    }    
  },
  methods: {
    getAmount(amountObj) {
      return amountObj.amount 
        ? amountObj.amount+amountObj.denom
        : amountObj[0].amount+amountObj[0].denom
    },
    getMsgType(type) {
      return type.replace('cosmos-sdk/', '')
    }
  }  
}
</script>

<style scoped>

/* sheet */
.sheet {
  width: 100%;
  height: 100%;
  background-color: var(--c-bg-third);
  overflow-y: scroll;
  overflow-x: hidden;
  padding-bottom: 2.5rem;
  box-sizing: border-box;
  
  /* border-left: 1px solid var(--c-theme-secondary); */
  /* box-shadow: -2px 0 6px rgba(0,0,0,.05); */
  border-radius: 8px;
  color: var(--c-txt-grey);
}
.sheet::-webkit-scrollbar { /* width */
  width: 6px;
}
.sheet::-webkit-scrollbar-track { /* Track */
  /* box-shadow: inset 0 0 1px var(--c-bg-grey);  */
  background: var(--c-bg-third); 
}
.sheet::-webkit-scrollbar-thumb { /* Handle */
  background: var(--c-bg-secondary); 
  border-radius: 10px;
}
.sheet::-webkit-scrollbar-thumb:hover { /* Handle on hover */
  background: var(--c-contrast-secondary); 
}

.sheet__top,
.sheet__sub {
  position: relative;
}
.sheet__sub {
  margin-left: 4px;
}
.sheet .-container {
  padding-left: 2rem;
  padding-right: 2rem;
}
.sheet .-border-btm:after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 1.5rem;
  width: calc(100% - 3.5rem);
  height: 1px;
  background-color: var(--c-theme-secondary);
}

.sheet .sheet__top {
  padding-top: 1.25rem;
  padding-bottom: 1rem;
}
.sheet .sheet__sub {
  padding-top: 1rem;
  padding-bottom: 1rem;
}
.sheet .sheet__main {
  padding-top: 1.5rem;
}

.sheet__heading {
  font-size: 1.25rem;
  color: var(--c-txt-primary);
}

/* cards__container */
.cards-container__top {
  margin-bottom: 1rem;
}
.cards-container__label {
  font-size: 0.875rem;
  color: var(--c-txt-primary);
  font-weight: 500;
  margin-left: 4px;
}
.cards-container__card {
  padding: 1.5rem 1.5rem 2rem 1.5rem;
  /* border: 1px solid var(--c-theme-secondary); */
  /* box-shadow: -2px 0 6px rgba(0,0,0,.05); */
  background-color: var(--c-bg-primary);
  border-radius: 8px;
}
.cards-container__card:not(:last-child) {
  margin-bottom: 1rem;
}

/* card */
.card__container:not(:last-child) {
  border-bottom: 1px solid var(--c-theme-secondary);
  margin-bottom: 1.25rem;
  padding-bottom: 1.25rem;
}

</style>