<template>
  <div class="accordion box">
    <slot/>
  </div>  
</template>

<script>
import AccordionItem from './AccordionItem'

export default {
  components: {
    AccordionItem,
  },
  props: {
    // content: {
    //   type: Array,
    //   required: true
    // },
    multiple: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      groupId: null
    }
  },  
  mounted() {
    this.groupId = this.$el.id
  }
}

export function accordionDataTemplate({
  id,
  active
}) {
  return {
    id,
    active
  }
}
</script>

<style>

</style>