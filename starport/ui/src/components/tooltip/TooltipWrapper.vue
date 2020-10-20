<template>
  <div 
    :class="['tooltip-wrapper', { '-is-active': isTooltipActive }]" 
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div class="tooltip-wrapper__trigger">
      <slot/>
    </div>
    <div 
      :class="['tooltip', `-${direction}`]"
      ref="tooltip"
      :style="{ bottom: bottomOffset }"
    >
      <span class="tooltip__content" v-html="content"></span>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    content: { type: String, required: true },
    isEventTriggerType: {
      type: Object,
      default: null,
      validator: function(value) {
        if (value == null || !value) return true
        
        return typeof value?.triggerActiveState === 'boolean'
      }
    },
    direction: {
      type: String,
      default: 'right',
      validator: function(value) {
        return ['top', 'right', 'left'].indexOf(value) !== -1
      }     
    }
  },
  data() {
    return {
      isActive: false,
      tooltipHeight: 0
    }
  },
  computed: {
    isTooltipActive() {
      return this.isActive
    },
    tooltipOffsetValue() {
      return this.tooltipHeight
    },
    bottomOffset() {
      if (this.direction === 'right' || this.direction === 'left') {
        return `calc(50% - ${this.tooltipOffsetValue}px / 2)`
      } else if (this.direction === 'top') {
        return `calc(100% + 6px)`
      }
    }
  },
  methods: {
    handleMouseEnter() {
      if (!this.isEventTriggerType) {
        this.isActive = true
      }
    },
    handleMouseLeave() {
      if (!this.isEventTriggerType) {
        this.isActive = false
      }
    },
    setTooltipHeight() {
      this.tooltipHeight = this.$refs.tooltip.clientHeight
    }
  },
  watch: {
    isEventTriggerType() {
      this.isActive = this.isEventTriggerType.triggerActiveState
    },
    content() {
      this.setTooltipHeight()
    }
  },
  created() {
    window.addEventListener('resize', this.setTooltipHeight)
  },
  mounted() {
    this.setTooltipHeight()
  },
  destroyed() {
    window.removeEventListener('resize', this.setTooltipHeight)
  }
}
</script>

<style>

.tooltip-wrapper {
  position: relative;
}

.tooltip {
  --tooltip-size: 6px;
}

.tooltip {
  position: absolute;
  bottom: 0;
  left: calc(100% + 0.85rem);
  padding: 0.35rem 0.6rem 0.4rem 0.6rem;
  border-radius: 6px;
  font-size: 0.875rem;
  line-height: 1.35;
  color: var(--c-txt-contrast-secondary);
  background-color: var(--c-bg-contrast-secondary);
  box-shadow: 0px 0px 8px rgba(0,0,0,.1);
  white-space: nowrap;
  overflow-wrap: normal;
  transition: opacity .3s;
}
.tooltip {
  opacity: 0;
  pointer-events: none;
  transition: opacity .3s;
}
.tooltip:before {
  content: '';
  position: absolute;
  top: calc(50% - var(--tooltip-size));
  left: calc(var(--tooltip-size) * -1);
  width: 0; 
  height: 0; 
  border-top: var(--tooltip-size) solid transparent;
  border-bottom: var(--tooltip-size) solid transparent; 
  border-right: calc(var(--tooltip-size) + 2px) solid var(--c-bg-contrast-secondary);
}
.tooltip-wrapper.-is-active .tooltip {
  opacity: 1;
  pointer-events: initial;
}

.tooltip.-top {
  left: -50%;
  width: 12vw;
  padding: 0.5rem 0.8rem 0.6rem 0.8rem;
  overflow-wrap: break-word;
  white-space: break-spaces;
}
.tooltip.-top:before {
  top: auto;
  bottom: calc(var(--tooltip-size) * -1 - 2px);
  left: calc(50% - var(--tooltip-size));
  transform: rotate(-90deg);
}
.tooltip.-left {
  left: auto;
  right: calc(100% + 0.8rem);
}
.tooltip.-left:before {
  left: auto;
  right: calc(var(--tooltip-size) * -1);
  transform: rotate(180deg);
}

.tooltip__content span {
  font-weight: 600;
  text-decoration: underline;
}

</style>