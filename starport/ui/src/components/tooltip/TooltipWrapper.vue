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
      class="tooltip"
      ref="tooltip"
      :style="{ bottom: `calc(50% - ${tooltipOffsetValue}px / 2)` }"
    >
      <span class="tooltip__content" v-html="content"></span>
    </div>
  </div>
</template>

<script>
export default {
  // TODO: implement direction options
  props: {
    content: { type: String, required: true },
    direction: {
      type: String,
      default: 'right',
      validator: function(value) {
        return ['top', 'right'].indexOf(value) !== -1
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
    }
  },
  methods: {
    handleMouseEnter() {
      this.isActive = true
    },
    handleMouseLeave() {
      this.isActive = false
    },
    setTooltipHeight() {
      this.tooltipHeight = this.$refs.tooltip.clientHeight
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
  position: absolute;
  bottom: 0;
  left: calc(100% + 1rem);
  padding: 0.5rem 0.8rem 0.6rem 0.6rem;
  border-radius: 4px;
  font-size: 0.875rem;
  line-height: 1.35;
  color: var(--c-txt-contrast-secondary);
  background-color: var(--c-bg-contrast-secondary);
  box-shadow: 0px 0px 8px rgba(0,0,0,.1);
  /* width: 15vw; */
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
  top: calc(50% - 6px);
  left: -6px;
  width: 0; 
  height: 0; 
  border-top: 6px solid transparent;
  border-bottom: 6px solid transparent; 
  border-right: 8px solid var(--c-bg-contrast-secondary);
}

.tooltip-wrapper.-is-active .tooltip {
  opacity: 1;
  pointer-events: initial;
}

.tooltip__content span {
  font-weight: 600;
  text-decoration: underline;
}

</style>