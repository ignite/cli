<template>
  <div class="side-tab-list container">
    <div class="container__left">
      <p>{{tab}}</p>
    </div>
    <div class="container__right">
      <div 
        v-for="(item, index) in list"
        :key="item.title+index"
        class="list"
      >
        <p class="list__title">{{item.title}}</p>
        <div class="list__items">
          <div
            v-for="(subItem, subIndex) in item.subItems"
            :key="item.title+index+subIndex"
            class="list__item"
          >
            <p class="list__item-content -title">{{subItem.title}}</p>
            <p class="list__item-content -content">{{subItem.content}}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    tab: { type: String, default: 'Messages' },
    list: { type: Array, required: true }
  }
}
</script>

<style scoped>

p {
  margin: 0;
  font-size: 0.875rem;
}

.container {
  display: flex;
}
.container__left {
  width: 15%;
}
.container__right {
  flex-grow: 1;
}
.container__left p {
  color: var(--c-txt-secondary);
  margin-top: 1px;
}

.list:not(:last-child) {
  margin-bottom: 1.8rem;
}

.list__title {
  border-bottom: 1px solid var(--c-theme-secondary);
  padding-bottom: 0.5rem;
  margin-bottom: 0.625rem;
  
  font-size: 1rem;
  font-weight: 600;
}
.list__item {
  display: flex;
}
.list__item:not(:last-child) {
  margin-bottom: 0.625rem;
}
.list__item-content:first-child {
  width: 25%;
}
.list__item-content:last-child {
  flex-grow: 1;
}
.list__item-content.-title {
  color: var(--c-txt-grey);
}

</style>