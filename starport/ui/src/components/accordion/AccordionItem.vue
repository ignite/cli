<template>
  <div
    :id="groupId + '-' + itemData.id"
    class="item" 
    :class="{'-is-active': itemData.isActive}"
  >
    <div
      class="item__trigger"
      @click="toggle"
      role="button"
    >
      <slot name="trigger"></slot>
    </div>
    <transition
      name="transition"
      @enter="startTransition"
      @after-enter="endTransition"
      @before-leave="startTransition"
      @after-leave="endTransition"
    >
      <div v-if="itemData.isActive" class="item__contents">
        <slot name="contents"></slot>
      </div>
    </transition>
  </div>  
</template>

<script>
export default {
  props: ['itemData', 'multiple', 'groupId'],
  methods: {
    getAccordionWrapper(node) {
      if (node.$el.id !== this.groupId) {
        this.getAccordionWrapper(node.$parent)
      }
      return node.$parent
    },
    getAccordionItem(parentNode) {  
      return parentNode.$children.map(childNode => {
        if (childNode.groupId !== this.groupId) {
          getAccordionItem(childNode.$el)
        }
        return childNode
      })
    },
    toggle(event) {
      if (this.multiple) {
        this.itemData.isActive = !this.itemData.isActive
        return
      }
      
      const $accordionWrapper = this.getAccordionWrapper(this.$parent)
      
      $accordionWrapper.$children.forEach((item) => {
        const $accordItem = this.getAccordionItem(item)
        $accordItem.forEach((accordItem) => {
          const isClickedItem = accordItem.$el.id === event.currentTarget.parentElement.id
          
          if (isClickedItem) {
            accordItem.itemData.isActive = !accordItem.itemData.isActive
          } else {
            accordItem.itemData.isActive = false
          }
        })
      }) 
      
    },
    startTransition(el) {
      el.style.height = el.scrollHeight + 'px'
    },
    endTransition(el) {
      el.style.height = ''
    }
  }  
}
</script>

<style scoped>

.item__trigger:hover {
  cursor: pointer;
}
.item__contents {
  overflow: hidden;
}

.transition-enter-active, .transition-leave-active {
  will-change: height;
  transition: height 0.2s ease;
}
.transition-enter, .transition-leave-to {
  height: 0 !important;
}

</style>