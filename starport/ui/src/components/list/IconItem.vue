<template>
  <div 
    :class="['item', {'-is-active': isActive}]"
  >
    <div class="item__head">
      <div v-if="!isActive" class="item__head-icon -is-loading"><Spinner/></div>
      <span v-else :class="['item__head-icon', `-is-${iconType}`]"></span>
    </div>
    <p v-if="!toInjectSlot" class="item__main">{{ itemText }}</p>    
    <div v-else><slot/></div>
    <!-- <TooltipWrapper :content="backendRunningStates[item.id] ? item.noteActive : item.noteInactive">
      <p class="item__main"><a :href="getBackendUrl(item.port)">{{item.name}}</a></p>
    </TooltipWrapper>             -->
  </div>  
</template>

<script>
import Spinner from '@/components/loaders/Spinner'

export default {
  components: {
    Spinner
  },
  props: {
    isActive: { type: Boolean, default: true },
    itemText: { type: String },
    toInjectSlot: { type: Boolean, default: false },
    iconType: {
      type: String,
      default: 'dot',
      validator(value) {
        const options = ['dot', 'check', 'slot']
        return options.filter(opt => opt === value).length>0
      }
    }
  }
}
</script>

<style scoped>

.item {
  --c-active: #4ACF4A;
  --c-active-sub: #7fe87f;
}
.item {
  display: flex;
  align-items: center;
  opacity: .6;
}

.item__head-icon {
  display: block;
  margin-right: 1rem;  
} 
.item__head-icon.-is-loading {
  width: 8px;
  height: 8px;
} 
.item__head-icon.-is-dot {
  width: 6px;
  height: 6px;  
  border-radius: 100%;  
  margin-top: 0.25px;
  background-color: var(--c-active);  
}
.item__head-icon.-is-check {
  width: 4px;  
  height: 8px;
  border-bottom: 2px solid var(--c-active);
  border-right: 2px solid var(--c-active);  
  transform: rotate(45deg);  
  margin-top: -2px;
}

.item__main {
  font-size: 1rem;
  color: var(--c-txt-grey);
}
.item__main a {
  position: relative;
  text-decoration: none;
}
.item__main a:after {
  content: '→';  
  position: absolute;
  top: 1px;
  right: -20px;
}
.item__main a:hover:after {
  right: -24px;
  transition: right .3s;
}

.item.-is-active {
  opacity: 1;
}
.item.-is-active .item__head-icon.-is-dot {
  animation: tempActiveEffect 2s ease-in-out infinite;
}
.item.-is-active .item__main {
  color: var(--c-txt-primary);
}
.item.-is-active .item__main a {
  color: var(--c-txt-highlight);
}


@keyframes tempLoadingEffect {
  0% { background-color: var(--c-txt-grey); }
  50% { background-color: var(--c-txt-secondary); }
  100% { background-color: var(--c-txt-grey); }
}
@keyframes tempActiveEffect {
  0% { background-color: var(--c-active); }
  50% { background-color: var(--c-active-sub); }
  100% { background-color: var(--c-active); }
}

</style>