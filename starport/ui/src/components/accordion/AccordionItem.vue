<template>
  <div
    :id="groupId + '-' + item.id"
    class="accordion-item" 
    :class="{'is-active': item.active}"
  >
    <dt class="accordion-item-title">
      <!-- <button @click="toggle" class="accordion-item-trigger"> -->
        <div @click="toggle">
        <slot name="trigger"></slot>
        </div>
      <!-- </button> -->
    </dt>
    <transition
      name="accordion-item"
      @enter="startTransition"
      @after-enter="endTransition"
      @before-leave="startTransition"
      @after-leave="endTransition">
      <dd v-if="item.active" class="accordion-item-details">
        <div v-html="item.content" class="accordion-item-details-inner"></div>
      </dd>
    </transition>
  </div>  
</template>

<script>
export default {
  props: ['item', 'multiple', 'groupId'],
  methods: {
    toggle(event) {
      if (this.multiple) this.item.active = !this.item.active
      else {
        this.$parent.$children.forEach((item, index) => {
          if (item.$el.id === event.currentTarget.parentElement.parentElement.id) item.item.active = !item.item.active
          else item.item.active = false
        }) 
      }
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


.accordion-item-details {
  overflow: hidden;
  background-color: whitesmoke;
}

.accordion-item-enter-active, .accordion-item-leave-active {
  will-change: height;
  transition: height 0.2s ease;
}
.accordion-item-enter, .accordion-item-leave-to {
  height: 0 !important;
}

</style>