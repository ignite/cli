<template>
  <IconText
    class="copy-icon-text"
    :text="text"
    :link="link"
    :iconType="'copy'"
    :isIconClickable="true"
    @clicked="handleClicked"
  />
</template>

<script>
import IconText from '@/components/texts/IconText'

export default {
  components: {
    IconText
  },
  props: {
    text: { type: String, required: true },
    link: { type: String, default: null },    
  },
  methods: {
    handleClicked() {
      console.log('test')
    }
  }
}
</script>

<style scoped>

/* .copy-icon-text >>> a {
  
} */

</style>