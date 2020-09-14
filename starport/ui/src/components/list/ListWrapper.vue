<template>
  <div class="list-wrapper">
    <div
      v-for="(item, index) in listItems"
      :key="item.headText+index"        
      class="list-wrapper__item"
    >
      <ListItem
        :headText="item.headText" 
        :contentText="item.contentText"
      />
    </div>
  </div>
</template>

<script>
import ListItem from './ListItem'

export default {
  props: {
    listItems: {
      type: Array,
      default: [],
      validator: function(value) {
        return value
          .filter(item => item?.headText && item?.contentText)
          .length === value.length
      }      
    }
  },
  components: {
    ListItem
  },  
}
</script>

<style scoped>

.list-wrapper__item:not(:last-child) {
  margin-bottom: 0.8rem;
}

</style>