<template>
  <div class="tx-card">
    <div class="tx-card__container">
      <ListWrapper :listItems="getFmtTxMeta(txData.meta)" />
    </div>
    <div class="tx-card__container">
      <SideTabList :list="getFmtTxMsg(txData.msgs)"/>                   
    </div>
  </div>  
</template>

<script>
import SideTabList from '@/components/table/SideTabList'
import ListWrapper from '@/components/list/ListWrapper'

export default {
  props: {
    txData: { type: Object, required: true } // TODO: add validator
  },  
  components: {
    SideTabList,
    ListWrapper
  },  
  methods: {
    getFmtTxMeta(txMeta) {
      const fmtTxMeta = []
      for (const [key, val] of Object.entries(txMeta)) {
        fmtTxMeta.push({ headText: key, contentText: val })
      }      
      return fmtTxMeta
    },
    getFmtTxMsg(msgs) {
      return msgs.map((msg) => {
        const fmtSubItems = []
        for (const [key, val] of Object.entries(msg)) {
          if (key !== 'type') {
            fmtSubItems.push({ headText: key, contentText: val })
          }
        }

        return {
          title: msg.type,
          subItems: fmtSubItems
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