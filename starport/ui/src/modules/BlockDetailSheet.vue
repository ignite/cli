<template>
  <div v-if="!block.data" :class="['sheet -is-empty']">Waiting for block data</div>

  <div v-else class="sheet">
    <div class="sheet__header">
      <div class="sheet__header-main">
        <p>{{block.data.height}}</p>
      </div>
      <div class="sheet__header-side">
        <div class="sheet__header-side-top">
          <CopyIconText :text="'hihihi'" :link="'#'" />
        </div>
        <div class="sheet__header-side-btm">
          <span>13:48:39, Oct 1 2019</span>
        </div>
      </div>
    </div>

    <div class="sheet__main">
      <div class="txs">
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
    </div>
  </div>
</template>

<script>
import CopyIconText from '@/components/texts/CopyIconText'

export default {
  components: {
    CopyIconText
  },
  props: {
    block: { type: Object }
  } 
}
</script>

<style scoped>

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

</style>