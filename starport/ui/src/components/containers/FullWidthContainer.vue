<template>
  <div 
    :class="['container', {'-is-collapsed': isSheetActive, '-is-loading': isTableLoading}]"
    ref="container"
  >

    <div class="container__wrapper">
      <div 
        v-if="hasSideSheet"
        :class="['container__sheet', {'-is-active': isSheetActive}]"
      >
        <slot name="sideSheet"/>
      </div>    

      <div class="container__main">
        <slot name="mainContent"/>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapMutations } from 'vuex'
import _ from 'lodash'

export default {
  props: {
    hasSideSheet: { type: Boolean, default: true },
    isTableLoading: { type: Boolean, default: false },
    tableEmptyMsg: { type: String, default: 'Waiting for blocks' },
  },
  data() {
    return {
      isSheetActive: true,
      lastScrolledHeight: 0,
      lastScrolledTop: 0,
      lastTimestamp: undefined
    }
  },
  methods: {
    handleSheetClose() {
      this.$emit('sheet-closed')
    },
    handleTableScroll: _.debounce(function(event) {
      const $table = event.target
      const scrolledHeight = $table.scrollTop + $table.offsetHeight
      const tableScrollHeight = $table.scrollHeight

      const isScrolledToTop = scrolledHeight <= $table.offsetHeight
      const isScrolledToBottom = scrolledHeight + 100 >= tableScrollHeight
      const isOnTopHalf = $table.scrollTop < (tableScrollHeight-$table.offsetHeight) / 2

      const isCallableScrolledDistance = 
        $table.offsetHeight / Math.abs(scrolledHeight-this.lastScrolledHeight) > 25
      
      if (isCallableScrolledDistance) {
        if (isScrolledToBottom) this.$emit('scrolled-bottom')
        if (isScrolledToTop) {
          this.$emit('scrolled-top')
          /*
           *
           // Scroll down the table a bit to prevent staying on top
           *
           */
          $table.scrollBy({
            top: 5,
            left: 0,
            behavior: 'smooth'
          })
        }
      }
    }, 250),
    updateScrollValue() {
      const $table = event.target
      const scrolledHeight = $table.scrollTop + $table.offsetHeight 
      this.lastScrolledHeight = scrolledHeight     
      this.lastScrolledTop = $table.scrollTop     
    }
  },
  created() {  
    // function checkScrollPosition(timestamp) {
    //   if (this.lastTimestamp === undefined) {
    //     this.lastTimestamp = timestamp
    //   }
    //   if (timestamp - this.lastTimestamp > 500) {
    //     if (this.lastScrolledTop === 0) this.$emit('scrolled-top')
    //     this.lastTimestamp = timestamp            
    //   }
    //   window.requestAnimationFrame(checkScrollPosition.bind(this))
    // }
    // window.requestAnimationFrame(checkScrollPosition.bind(this))    
  }
}
</script>

<style scoped>

.container {
  --container-collapsed-width: 20vw;
}
.container {
  position: relative;
  height: inherit;
}

.container__wrapper {
  position: relative;
  height: inherit;
  overflow: hidden;

  transition: opacity .3s ease-in-out;
}
.container__wrapper:before,
.container__wrapper:after {
  position: absolute;
  z-index: 1;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  opacity: 0;
  transition: all .3s ease-in-out;  
}
.container__wrapper:before {
  content: '';
  background-color: var(--c-bg-secondary);
}
.container__wrapper:after {
  content: 'Fetching blocks';
  display: flex;  
  justify-content: center;
  align-items: center;
  color: var(--c-txt-grey);
  animation: tempLoadingEffect 2s ease-in-out infinite;  
}
.container.-is-loading .container__wrapper:before {
  opacity: .8;
}
.container.-is-loading .container__wrapper:after {
  opacity: 1;
}
.container.-is-loading .container__wrapper:before,
.container.-is-loading .container__wrapper:after {
  pointer-events: initial;
  transition: all .3s ease-in-out;
}

.container__main {
  box-sizing: border-box;
  height: inherit;
  max-height: inherit;
}
.container .container__main {
  width: 100%;
  transition: width .3s ease-in-out;
  will-change: width;
}
.container.-is-collapsed .container__main {
  width: var(--container-collapsed-width);
  transition: width .3s ease-in-out;
  will-change: width;
}

.container__sheet {
  position: absolute;
  top: 0;
  right: 0;
  width: calc(100% - var(--container-collapsed-width) - 2.5rem);
  height: 100%;  
}
.container__sheet {  
  transform: translate3d(100%, 0, 0);
  opacity: 0;
  transition: transform ease-out .25s;  
  will-change: transform;
}
.container__sheet.-is-active {  
  transform: translate3d(0,0,0);
  opacity: 1;  
  transition: transform ease-out .3s;
  will-change: transform;
}

@keyframes tempLoadingEffect {
  0% { color: var(--c-txt-grey); }
  50% { color: var(--c-txt-secondary); }
  100% { color: var(--c-txt-grey); }
}


@media only screen and (max-width: 1400px) {
  .container {
    --container-collapsed-width: 20vw;
  }
}
@media only screen and (max-width: 992px) {
  .container {
    min-width: 850px;
  }
  .container.-is-collapsed .container__main {
    width: 320px;
  }
  .container__sheet {
    width: calc(100% - 320px - 1rem);  
  }  
}


</style>