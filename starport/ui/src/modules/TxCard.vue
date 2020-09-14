<template>
  <div class="tx-card">
    <div class="tx-card__container">
      <ListWrapper :listItems="[
        { headText: 'TxHash', contentText: txData.txMsg.hash },
        { headText: 'Status', contentText: txData.txMsg.status },
        { headText: 'Fee', contentText: txData.txMsg.fee },
        { headText: 'Gas', contentText: txData.txMsg.gas },
        { headText: 'Memo', contentText: txData.txMsg.memo }
      ]" />
    </div>
    <div class="tx-card__container">
      <SideTabList
        :list="getFmtMsgForInnerTable(txData.msgs)"
      />                   
    </div>
  </div>  
</template>

<script>
import SideTabList from '@/components/table/SideTabList'
import ListWrapper from '@/components/list/ListWrapper'

export default {
  props: {
    txData: { type: Object, required: true }
    // TODO: add validator
  },  
  components: {
    SideTabList,
    ListWrapper
  },  
  methods: {
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
  }
}
</script>

<style scoped>

.tx-card__container:not(:last-child) {
  border-bottom: 1px solid var(--c-theme-secondary);
  margin-bottom: 1.25rem;
  padding-bottom: 1.25rem;
}

</style>