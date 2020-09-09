<template>
  <div v-if="!blockData">waiting for block data...</div>
  <div 
    v-else
    :class="['sheet']"
    ref="tableInnerSheet"
  >
    <div class="sheet__top -container -border-btm">
      <h3 class="sheet__heading">Block #{{blockData.blockMsg.height}}</h3>
    </div>
    <div class="sheet__sub -container -border-btm">
      <div class="sheet__entry">
        <span class="sheet__entry-label">Hash</span>
        <p class="sheet__entry-content">{{blockData.blockMsg.blockHash}}</p>
      </div>
      <div class="sheet__entry">
        <span class="sheet__entry-label">Time</span>
        <p class="sheet__entry-content">{{blockData.blockMsg.time}}</p>
      </div>
    </div>
    <div class="sheet__main -container">
      <div class="cards-container">
        <div class="cards-container__top">
          <h4 class="cards-container__label">Transactions</h4>
        </div>

        <!-- transactions -->
        <div 
          v-if="blockData.blockMsg.txs>0 && blockData.txs.length>0"
          class="cards-container__card"
        >
          <div 
            v-for="tx in messagesForTable"
            :key="tx.tableData.id"
            class="card"
          >
            <div class="card__container">
              <div class="sheet__entry">
                <span class="sheet__entry-label">TxHash</span>
                <p class="sheet__entry-content">{{tx.txMsg.hash}}</p>
              </div>
              <div class="sheet__entry">
                <span class="sheet__entry-label">Status</span>
                <p class="sheet__entry-content">{{tx.txMsg.status}}</p>
              </div>
              <div class="sheet__entry">
                <span class="sheet__entry-label">Fee</span>
                <p class="sheet__entry-content">{{tx.txMsg.fee}}</p>
              </div>
              <div class="sheet__entry">
                <span class="sheet__entry-label">Gas</span>
                <p class="sheet__entry-content">{{tx.txMsg.gas}}</p>
              </div>
              <div class="sheet__entry">
                <span class="sheet__entry-label">Memo</span>
                <p class="sheet__entry-content">{{tx.txMsg.memo}}</p>
              </div>
            </div>
            <div class="card__container">
              <SideTabList
                :list="getFmtMsgForInnerTable(tx.msgs)"
              />                   
            </div>
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

export default {
  props: {
    blockData: { type: Object }
  },
  components: {
    SideTabList
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
    getFmtMsgForInnerTable(msgs) {
      return msgs.map(({
        type,
        amount,
        delegator,
        validator,
        from,
        to
      }) => {
        const fromType = type === 'MsgSend' ? {
          title: 'From', content: from
        } : {
          title: 'Delegator', content: delegator
        }
        const toType = type === 'MsgSend' ? {
          title: 'To', content: to
        } : {
          title: 'Validator', content: validator
        }

        return {
          title: type,
          subItems: [
            fromType,
            toType,
            { title: 'Amount', content: amount },
          ]
        }
      })
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
  background-color: var(--c-bg-primary);
  
  border-left: 1px solid var(--c-theme-secondary);
  box-shadow: -2px 0 6px rgba(0,0,0,.05);
}

.sheet__top,
.sheet__sub {
  position: relative;
}
.sheet .-container {
  padding-left: 1rem;
  padding-right: 1rem;
}
.sheet .-border-btm:after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 1rem;
  width: calc(100% - 2rem);
  height: 1px;
  background-color: var(--c-theme-secondary);
}

.sheet .sheet__top {
  padding-top: 0.8rem;
  padding-bottom: 0.8rem;
}
.sheet .sheet__sub {
  padding-top: 1rem;
  padding-bottom: 1rem;
}
.sheet .sheet__main {
  padding-top: 1.5rem;
}

.sheet__heading {
  font-size: 1rem;
}

.sheet__entry {
  display: flex;
}
.sheet__entry:not(:last-child) {
  margin-bottom: 0.8rem;
}
.sheet__entry span,
.sheet__entry p  {
  font-size: 0.875rem;
}
.sheet__entry *:first-child {
  width: 15%;
  color: var(--c-txt-grey);
}
.sheet__entry *:last-child {
  flex-grow: 1;
}

/* cards__container */
.cards-container__top {
  margin-bottom: 1rem;
}
.cards-container__label {
  font-size: 0.9375rem;
  color: var(--c-txt-grey);
}
.cards-container__card {
  padding: 1rem 1.25rem 1.5rem 1.25rem;
  border: 1px solid var(--c-theme-secondary);
  box-shadow: -2px 0 6px rgba(0,0,0,.05);
}

/* card */
.card__container:not(:last-child) {
  border-bottom: 1px solid var(--c-theme-secondary);
  margin-bottom: 1.25rem;
  padding-bottom: 1.25rem;
}

</style>