<template>
  <div 
    :class="['sheet']"
    ref="tableInnerSheet"
  >
    <div class="sheet__top -container -border-btm">
      <h3 class="sheet__heading">Block #3,290,279</h3>
    </div>
    <div class="sheet__sub -container -border-btm">
      <div class="sheet__entry">
        <span class="sheet__entry-label">Hash</span>
        <p class="sheet__entry-content">BBE79CF80CD48AC8308D5611653E1FFDD48F819C251F4FF795FDD70A9F9F436A</p>
      </div>
      <div class="sheet__entry">
        <span class="sheet__entry-label">Time</span>
        <p class="sheet__entry-content">9 Sep 2020, 2:45:24am UTC (50 seconds ago)</p>
      </div>
    </div>
    <div class="sheet__main -container">
      <div class="cards-container">
        <div class="cards-container__top">
          <h4 class="cards-container__label">Transactions</h4>
        </div>

        <div class="cards-container__card">
          <div class="card">
            <div class="card__container">
              <div class="sheet__entry">
                <span class="sheet__entry-label">TxHash</span>
                <p class="sheet__entry-content">37549643B101E3D1FBD3955032A4FCEFBC720950A6A1E592C11E555F079FFF15</p>
              </div>
              <div class="sheet__entry">
                <span class="sheet__entry-label">Status</span>
                <p class="sheet__entry-content">Success</p>
              </div>
              <div class="sheet__entry">
                <span class="sheet__entry-label">Fee</span>
                <p class="sheet__entry-content">0.000750ATOM</p>
              </div>
              <div class="sheet__entry">
                <span class="sheet__entry-label">Gas</span>
                <p class="sheet__entry-content">41432</p>
              </div>
              <div class="sheet__entry">
                <span class="sheet__entry-label">Memo</span>
                <p class="sheet__entry-content">Testing 123</p>
              </div>
            </div>
            <div class="card__container">
              <SideTabList
                :list="[
                  {
                    title: 'Send',
                    subItems: [
                      { title: 'Delegator', content: 'cosmos1hvsdf03tl6w5pnfvfv5g8uphjd4wfw2hsucxnd' },
                      { title: 'Validator', content: 'cosmosvaloper1hvsdf03tl6w5pnfvfv5g8uphjd4wfw2h' },
                      { title: 'Reward Amount', content: '0.000012 ATOM' },
                    ]
                  }
                ]"
              />                   
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script>
import SideTabList from '@/components/table/SideTabList'

export default {
  components: {
    SideTabList
  },
  data() {
    return {
      isActive: false
    }
  }
}
</script>

<style scoped>

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

.card__container:not(:last-child) {
  border-bottom: 1px solid var(--c-theme-secondary);
  margin-bottom: 1.25rem;
  padding-bottom: 1.25rem;
}

</style>