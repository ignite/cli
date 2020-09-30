<template>
  <div class="side-tab-list">
    <ListItem :headText="'Messages'" :isSlotContentType="true">
      <div
        slot="content"
        class="side-tab-list__contents"
      >
        <div 
          v-for="(item, index) in list"
          :key="item.title+index"
          class="side-tab-list__content"
        >
          <p class="side-tab-list__content-title">{{item.title}}</p>
          <ListWrapper 
            :listItems="item.subItems" />   
        </div>   
      </div>    
    </ListItem>
  </div>
</template>

<script>
import ListItem from '@/components/list/ListItem'
import ListWrapper from '@/components/list/ListWrapper'

export default {
  props: {
    tab: { type: String, default: 'Messages' },
    list: { type: Array, required: true }
  },
  components: {
    ListItem,
    ListWrapper,
  }  
}
</script>

<style scoped>

p.side-tab-list__content-title {
  border-bottom: 1px solid var(--c-theme-secondary);
  padding-bottom: 0.5rem;
  margin-bottom: 0.625rem;
  
  /* font-size: 1rem; */
  font-weight: 600;
  color: var(--c-txt-secondary);
}


/* @media only screen and (max-width: 1400px) {
  .side-tab-list >>> .list-item__head {
    margin-bottom: 0.8rem;
  }
} */

</style>