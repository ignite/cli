<template>
  <dl class="accordion box" role="presentation">
    <accordion-item
      v-for="item in content"
      :multiple="multiple"
      :item="item"
      :groupId="groupId"
      :key="item.id"
    >
      <template slot="trigger">
        <slot name="trigger"></slot>
      </template>
    </accordion-item>
  </dl>  
</template>

<script>
import AccordionItem from './AccordionItem'

export default {
  components: {
    AccordionItem,
  },
  props: {
    content: {
      type: Array,
      required: true
    },
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
  active,
  title,
  content
}) {
  return {
    id,
    active,
    title,
    content
  }
}
</script>

<style>

</style>