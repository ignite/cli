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
    }
  },
  methods: {
    handleSheetClose() {
      this.$emit('sheet-closed')
    }
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
  width: calc(100% - var(--container-collapsed-width) - 4rem);
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


@media only screen and (max-width: 1200px) {
  .container {
    --container-collapsed-width: 25vw;
  }
}
@media only screen and (max-width: 992px) {
  .container {
    --container-collapsed-width: 25vw;
  }  
  .container__sheet {
    width: calc(100% - var(--container-collapsed-width) - 2rem);
  }
}
@media screen and (max-width: 768px) {
  .container {
    --container-collapsed-width: 240px;
  }
  .container__wrapper {
    overflow-x: scroll;
  }    
  .container__main {
    min-width: 240px;
  }
}


</style>