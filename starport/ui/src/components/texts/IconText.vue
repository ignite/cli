<script>
import IconCopy from '@/assets/icons/Copy'

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
      validator(value) {
        return ['copy'].indexOf(value)>=0
      }
    }
  },
  methods: {
    iconComp() {
      switch (this.iconType) {
        case 'copy':
          return <IconCopy/>
        default:
          return <IconCopy/>
      }
    }
  },
  render(h) {
    const { text, link, isIconClickable, iconType } = this.$props

    const textContent = !link ? (
      <span>{text}</span>
    ) : (
      <a href={link}>{text}</a>
    )

    const iconContent = !isIconClickable ? (
      <span>{this.iconComp()}</span>
    ) : (
      <button onClick={() => this.$emit('clicked')}>{this.iconComp()}</button>
    )

    return (
      <p class="main">
        {textContent}
        <span class={`icon -is-${this.iconType}`}>{iconContent}</span>
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
  transform: translate3d(0, 1.5px, 0);
}
.icon * >>> svg path {
  fill: var(--c-txt-third);
}

</style>