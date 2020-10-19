<script>
import IconCopy from '@/assets/icons/Copy'
import TooltipWrapper from '@/components/tooltip/TooltipWrapper'

export default {
  components: {
    IconCopy
  },
  props: {
    text: { type: String, required: true },
    link: { type: String, default: null },
    isIconClickable: { type: Boolean, default: false },    
    iconType: {
      type: String,
      required: true,
      default: 'copy',
      validator(value) {
        return ['copy'].indexOf(value)>=0
      }
    },
    tooltipOption: {
      type: String,
      default: 'none',
      validator(value) {
        return ['none', 'iconWrapper', 'textWrapper', 'compWrapper'].indexOf(value)>=0
      }
    },
    tooltipStates: { type: Object, default: null },
    tooltipDirection: { type: String, default: 'right' }
  },
  methods: {
    getIconType() {
      switch (this.iconType) {
        case 'copy':
          return <IconCopy/>
        default:
          return <IconCopy/>
      }
    },
    getIconContent() {
      return !this.isIconClickable ? (
        <span>{this.getIconType()}</span>
      ) : (
        <button onClick={() => this.$emit('iconClicked')}>{this.getIconType()}</button>
      )
    },
    getIconComp() {
      const IconContent = this.getIconContent()

      switch (this.tooltipOption) {
        case 'iconWrapper':
          return (
            <TooltipWrapper 
              content={this.tooltipStates.text} 
              isEventTriggerType={{ triggerActiveState: this.tooltipStates.state }}
              direction={this.tooltipDirection}
            >
              {IconContent}      
            </TooltipWrapper>              
          )
        default:
          return IconContent
      }
    }
  },
  render(h) {
    const { text, link, isIconClickable, iconType } = this.$props

    const textContent = !link ? (
      <span>{text}</span>
    ) : (
      <a href={link} target="_blank">{text}</a>
    )

    return (
      <p class="main">
        {textContent}
        <span class={`icon -is-${this.iconType}`}>{this.getIconComp()}</span>
      </p>      
    )
  }
}
</script>

<style scoped>

.main {
  font-size: 1rem;
}
.main > *:first-child {
  margin-right: 4px;
}

.main a {
  color: var(--c-txt-highlight)
}

.icon.-is-copy {
  display: inline-block;
  transform: translate3d(0, 1px, 0);
}
.icon * >>> svg path {
  fill: var(--c-txt-third);
}

</style>