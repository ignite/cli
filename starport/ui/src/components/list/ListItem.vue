<template>
  <div :class="['list-item', getItemClass]">
    <span v-if="!isSlotHeadType" class="list-item__head">{{headText}}</span>
    <div v-else class="list-item__head"><slot name="head"/></div>

    <p v-if="!isSlotContentType" class="list-item__content">{{contentText}}</p>    
    <div v-else class="list-item__content"><slot name="content"/></div>
  </div>
</template>

<script>
export default {
  props: {
    headText: { type: String, default: '' },
    contentText: { type: String, default: '' },
    isSlotHeadType: { type: Boolean, default: false },
    isSlotContentType: { type: Boolean, default: false }
  },
  computed: {
    getItemClass() {
      if (this.headText === 'Status') {
        return ['-prefix-dot', this.contentText === 'Fail' ? '-is-warn' : '-is-success']
      }
    }
  }
}
</script>

<style scoped>

.list-item {
  display: flex;
}

.list-item * {
  font-size: 0.875rem;
}
.list-item__head {
  width: 20%;
  min-width: 20%;
  margin-right: 0.5rem;  
  margin-top: 1px;
  font-size: 0.8125rem;  
  color: var(--c-txt-grey);
}
.list-item__content {
  position: relative;
  flex: 1 0;
  color: var(--c-txt-secondary);
  overflow-wrap: break-word;
  word-break: break-word;
}
.list-item.-prefix-dot .list-item__content {
  text-indent: 12px;
}
.list-item.-prefix-dot .list-item__content:before {
  content: '';
  position: absolute;
  top: calc(51% - 6px / 2);
  left: 1px;
  width: 6px;
  height: 6px;
  border-radius: 100%;
}
.list-item.-prefix-dot.-is-warn .list-item__content:before {
  background-color: #FF1A1A;
}
.list-item.-prefix-dot.-is-success .list-item__content:before {
  background-color: #4ACF4A;
}

@media only screen and (max-width: 1400px) {
  .list-item {
    flex-direction: column;
  }
  .list-item__head {
    width: 100%;
    font-size: 0.8125rem;
    margin-bottom: 0.35rem;
  }
}

</style>