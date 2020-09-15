<template>
  <div
    :id="groupId + '-' + itemData.id"
    class="accord-item" 
    :class="[
      itemData.isActive ? '-is-active' : '',
      isDisabled ? '-is-disabled' : ''
    ]"
  >
    <div
      class="accord-item__trigger"
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
      <div v-if="localItemData.isActive" class="accord-item__contents">
        <slot name="contents"></slot>
      </div>
    </transition>
  </div>  
</template>

<script>
export default {
  props: {
    itemData: { type: Object, required: true },
    isDisabled: { type: Boolean, default: false },
    groupId: { type: String, required: true },
    multiple: { type: Boolean, default: false }
  },
  data() {
    return {
      localItemData: this.itemData
    }
  },
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
    setParentState(node, parentTag, toState=undefined)  {
      const isTargetParent = node.$options._componentTag === parentTag
      const isTableWrapper = node.$options._componentTag === 'TableWrapper'

      if (!isTargetParent && !isTableWrapper) {
        this.setParentState(node.$parent, parentTag, toState)
      }
      
      if (node.isActive !== undefined) {
        if (toState === undefined) {
          node.isActive = !node.isActive
        } else {
          node.isActive = toState
        }
      }
    },
    toggle(event) {
      if (this.multiple) {
        this.localItemData.isActive = !this.localItemData.isActive
        return
      }
      if (this.isDisabled) return
      
      const $accordionWrapper = this.getAccordionWrapper(this.$parent)
      
      $accordionWrapper.$children.forEach((item) => {
        const $accordItem = this.getAccordionItem(item)
        $accordItem.forEach((accordItem) => {
          const isClickedItem = accordItem.$el.id === event.currentTarget.parentElement.id
          
          if (isClickedItem) {
            accordItem.localItemData.isActive = !accordItem.localItemData.isActive
            this.setParentState(accordItem, 'TableRowWrapper')
          } else {
            accordItem.localItemData.isActive = false
            this.setParentState(accordItem, 'TableRowWrapper', false)
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

.accord-item:not(.-is-disabled) .accord-item__trigger:hover {
  cursor: pointer;
}
.accord-item__contents {
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