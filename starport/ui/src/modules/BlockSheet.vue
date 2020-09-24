<template>
  <div v-if="!fmtBlockData">waiting for block data...</div>
  
  <div 
    v-else
    :class="['sheet']"
  >
    <div class="sheet__top -container -border-btm">
      <h3 class="sheet__heading">Block #{{fmtBlockData.blockMsg.height}}</h3>
      <button class="sheet__btn" @click="handleJsonCopy">Copy block JSON</button>
    </div>
    <div class="sheet__sub -container -border-btm">
      <ListWrapper :listItems="[
        { headText: 'Hash', contentText: fmtBlockData.blockMsg.blockHash },
        { headText: 'Time', contentText: fmtBlockData.blockMsg.time },
      ]" />      
    </div>
    <div class="sheet__main -container">
      <div class="cards-container">

        <div class="cards-container__top">
          <h4 class="cards-container__label">Transactions</h4>
        </div>

        <!-- transactions -->
        <div v-if="fmtBlockData.blockMsg.txs>0 && blockData.txs.length>0">
          <div 
            v-for="tx in messagesForTable"
            :key="tx.tableData.id"          
            class="cards-container__card"
          >
            <TxCard :txData="tx" />
          </div>
        </div>
        <div 
          v-else-if="fmtBlockData.blockMsg.txs>0 && blockData.txs.length<=0"
          class="cards-container__card -is-empty"
        >
          <p>🚨 Error fetching transaction data</p>
        </div>
        <div v-else class="cards-container__card -is-empty">
          <p>No transactions</p>
        </div>

      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

import blockHelpers from '@/mixins/blocks/helpers'

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
    /*
     *
     * Vuex 
     *
     */    
    ...mapGetters('cosmos/blocks', [ 'chainId' ]),
    /*
     *
     * Local 
     *
     */    
    messagesForTable() {
      return blockHelpers.blockFormatter()
        .txForCard(this.fmtBlockData.txs, this.chainId)
    },
    fmtBlockData() {
      return this.blockData.data
    }
  },
  methods: {
    handleJsonCopy() {
      function fallbackCopyTextToClipboard(text) {
        var textArea = document.createElement("textarea")
        textArea.value = text
        
        // Avoid scrolling to bottom
        textArea.style.top = "0"
        textArea.style.left = "0"
        textArea.style.position = "fixed"

        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()

        try {
          var successful = document.execCommand('copy')
          var msg = successful ? 'successful' : 'unsuccessful'
          console.log('Fallback: Copying text command was ' + msg)
        } catch (err) {
          console.error('Fallback: Oops, unable to copy', err)
        }

        document.body.removeChild(textArea)
      }
      function copyTextToClipboard(text) {
        if (!navigator.clipboard) {
          fallbackCopyTextToClipboard(text)
          return
        }
        navigator.clipboard.writeText(text).then(function() {
          console.log('Async: Copying to clipboard was successful!')
        }, function(err) {
          console.error('Async: Could not copy text: ', err)
        })
      }

      copyTextToClipboard(JSON.stringify(this.blockData.rawJson))
    }
  }
}
</script>

<style scoped>

/* sheet */
.sheet {
  width: 100%;
  height: 100%;
  background-color: var(--c-bg-secondary);
  overflow-y: scroll;
  overflow-x: hidden;
  padding-bottom: 2.5rem;
  box-sizing: border-box;
  
  /* border-left: 1px solid var(--c-theme-secondary); */
  /* box-shadow: -2px 0 6px rgba(0,0,0,.05); */
  border-radius: var(--bd-radius-primary);
  color: var(--c-txt-grey);
}
.sheet::-webkit-scrollbar { /* width */
  width: 6px;
}
.sheet::-webkit-scrollbar-track { /* Track */
  /* box-shadow: inset 0 0 1px var(--c-bg-grey);  */
  background: var(--c-bg-secondary); 
}
.sheet::-webkit-scrollbar-thumb { /* Handle */
  background-color: var(--c-bg-third); 
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

.sheet__top {
  padding-top: 1.25rem;
  padding-bottom: 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.sheet__sub {
  padding-top: 1rem;
  padding-bottom: 1rem;
}
.sheet__main {
  padding-top: 1.5rem;
}
.sheet__btn {
  font-size: 0.8125rem;
  color: var(--c-txt-grey);
  transition: color .3s;
}
.sheet__btn:hover {
  color: var(--c-txt-secondary);
  transition: color .3s;
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
  border-radius: var(--bd-radius-primary);
}
.cards-container__card:not(:last-child) {
  margin-bottom: 1rem;
}
.cards-container__card.-is-empty {
  padding-top: 4rem;
  padding-bottom: 4rem;
  text-align: center;
  height: 100%;
}
.cards-container__card.-is-empty p {
  font-weight: 300;
}

/* card */
.card__container:not(:last-child) {
  border-bottom: 1px solid var(--c-theme-secondary);
  margin-bottom: 1.25rem;
  padding-bottom: 1.25rem;
}

</style>