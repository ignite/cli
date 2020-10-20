<template>
  <div>
    <transition name="overlay" appear>
      <div
        class="overlay"
        ref="overlay"
        :style="{
          'background-color': backgroundColor || 'rgba(0, 0, 0, 0.35)',
        }"
        v-if="visible && visibleLocal"
        @click="close"
        @touchstart="touchstart"
        @touchmove="touchmove"
        @touchend="touchend"
      ></div>
    </transition>
    <transition :name="`sidebar__${side}`" @after-leave="emitVisible()" appear>
      <div
        :class="['sidebar', `sidebar__side__${side}`]"
        ref="sidebar"
        v-if="visible && visibleLocal"
        :style="style"
        @click.self="sidebarClick"
        @touchstart="touchstart"
        @touchmove="touchmove"
        @touchend="touchend"
      >
        <!-- @slot Contents of the sidebar. -->
        <div
          @scroll="setScrolling(true)"
          ref="content"
          :class="[
            `sidebar__content`,
            `sidebar__content__side__${side}`,
            `sidebar__fullscreen__${!!fullscreenComputed}`,
          ]"
        >
          <slot />
        </div>
        <div
          class="close"
          v-if="side === 'center' && buttonClose"
          @click="close"
        >
          <svg
            class="close__icon"
            width="100%"
            height="100%"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M4 4l16 16m0-16L4 20"
              stroke-width="1.5"
              stroke-linecap="round"
            />
          </svg>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 100000;
  height: 100vh;
  width: 100vw;
}
.sidebar {
  position: fixed;
  overflow-y: hidden;
  z-index: 100000;
  transform: translateX(var(--sidebar-translate-x))
    translateY(var(--sidebar-translate-y));
  -webkit-overflow-scrolling: touch;
}
.close {
  border-radius: 50%;
  stroke: rgba(255, 255, 255, 0.75);
  box-sizing: border-box;
  padding: 8px;
  width: 48px;
  height: 48px;
  position: absolute;
  top: 0;
  right: 0;
  margin: 1rem;
  cursor: pointer;
}
.close__icon {
  display: block;
}
.sidebar.sidebar__side__left {
  top: 0;
  left: 0;
  right: initial;
  width: var(--sidebar-width, 300px);
  max-width: var(--sidebar-max-width, 75%);
  height: var(--sidebar-height, 100%);
  max-height: var(--sidebar-max-height, 100%);
  box-shadow: var(--sidebar-box-shadow);
}
.sidebar.sidebar__side__right {
  top: 0;
  left: initial;
  right: 0;
  width: var(--sidebar-width, 300px);
  max-width: var(--sidebar-max-width, 75%);
  height: var(--sidebar-height, 100%);
  max-height: var(--sidebar-max-height, 100%);
  box-shadow: var(--sidebar-box-shadow);
}
.sidebar.sidebar__side__bottom {
  top: 0;
  left: 0;
  right: 0;
  max-width: initial;
  overflow-y: scroll;
  -webkit-overflow-scrolling: touch;
  width: var(--sidebar-width, 100%);
  max-width: var(--sidebar-width, 100%);
  height: 100%;
  margin-left: auto;
  margin-right: auto;
}
.sidebar.sidebar__side__center {
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  max-width: 100%;
  overflow-y: scroll;
  -webkit-overflow-scrolling: touch;
  pointer-events: all;
  height: 100vh;
}
.sidebar__content {
  background: var(--c-bg-primary);
  position: absolute;
  pointer-events: all;
  /* overflow-y: scroll; */
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  width: 100%;
  height: 100%;
}
.sidebar__content.sidebar__content__side__bottom {
  margin-top: var(--sidebar-margin-top);
  overflow-y: hidden;
  height: auto;
  box-shadow: var(--sidebar-box-shadow);
}

.sidebar__content.sidebar__content__side__center {
  pointer-events: all;
  position: absolute;
  width: var(--sidebar-width, 600px);
  max-width: var(--sidebar-max-width, 90%);
  height: var(--sidebar-height, auto);
  max-height: var(--sidebar-max-height, none);
  top: var(--sidebar-margin-top);
  transform: translateX(-50%);
  left: 50%;
  border-radius: var(--sidebar-border-radius);
  box-shadow: var(--sidebar-box-shadow);
  margin-bottom: 20px;
  border-radius: 0.5rem;
}

.sidebar__content.sidebar__content__side__center.sidebar__fullscreen__true {
  top: 0;
  width: 100%;
  height: 100%;
  max-width: 100%;
  max-height: 100%;
  margin-bottom: initial;
  border-radius: 0;
}

.overlay-enter-active {
  transition: all 0.25s ease-out;
}
.overlay-enter {
  opacity: 0;
}
.overlay-enter-to {
  opacity: 1;
}
.overlay-leave-active {
  transition: all 0.25s;
}
.overlay-leave {
  opacity: 1;
}
.overlay-leave-to {
  opacity: 0;
}
.sidebar__left-enter-active,
.sidebar__right-enter-active,
.sidebar__bottom-enter-active,
.sidebar__center-enter-active,
.sidebar__left-leave-active,
.sidebar__right-leave-active,
.sidebar__bottom-leave-active,
.sidebar__center-leave-active {
  transition: all 0.5s;
}
.sidebar__left-enter,
.sidebar__left-leave-to {
  transform: translateX(-100%);
}
.sidebar__right-enter,
.sidebar__right-leave-to {
  transform: translateX(100%);
}
.sidebar__left-enter-to,
.sidebar__left-leave .sidebar__right-enter-to,
.sidebar__right-leave {
  transform: translateX(0);
}
.sidebar__bottom-enter,
.sidebar__bottom-leave-to {
  transform: translateY(100%);
}
.sidebar__bottom-enter-to,
.sidebar__bottom-leave {
  transform: translateY(0);
}
.sidebar__center-enter,
.sidebar__center-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
.sidebar__center-enter-to,
.sidebar__center-leave {
  opacity: 1;
  transform: scale(1);
}

@media screen and (max-width: 500px) {
  /* TODO: hotfix for https://github.com/cosmos/cosmos.network/pull/832#pullrequestreview-413650189 */
  .sidebar__content.sidebar__content__side__center {
    margin-top: 3rem;
  }
}
</style>

<script>
/**
 * `Sidebar` is a sheet that transitions on top of the main content and allows
 * displaying auxiliary information, such as table of contents or global navigation.
 */
export default {
  props: {
    /**
     * Toggles visibility of the component.
     */
    visible: {
      type: Boolean,
      default: false
    },
    /**
     * Width of the sidebar.
     */
    width: {
      type: String,
      default: ''
    },
    /**
     * Maximum width of the sidebar.
     */
    maxWidth: {
      type: String,
      default: ''
    },
    /**
     * Height of the sidebar.
     */
    height: {
      type: String,
      default: ''
    },
    /**
     * Maximum height of the sidebar.
     */
    maxHeight: {
      type: String,
      default: ''
    },
    /**
     * `left` | `right` | `bottom`
     */
    side: {
      type: String,
      default: "left"
    },
    /**
     * CSS `background-color` of the overlay.
     */
    backgroundColor: {
      type: String,
      default: "rgba(0, 0, 0, 0.35)"
    },
    /**
     * CSS `box-shadow` of the sidebar sheet.
     */
    boxShadow: {
      type: String,
      default: "none"
    },
    /**
     * Vertical height of overlay
     */
    marginTop: {
      type: String,
      default: ''
    },
    /**
     * Go fullscreen when viewport is narrower than width
     */
    fullscreen: {
      type: Boolean,
      default: false
    },
    /**
     * Add default close button for centered modal
     */
    buttonClose: {
      type: Boolean,
      default: false
    }
  },
  watch: {
    visible(becomesVisible) {
      if (becomesVisible) {
        this.visibleLocal = true;
      }
    }
  },
  data: function() {
    return {
      visibleLocal: true,
      startX: null,
      startY: null,
      currentX: null,
      currentY: null,
      translateX: null,
      translateY: null,
      isScrolling: null,
      marginTopComputed: null,
      fullscreenComputed: null
    };
  },
  watch: {
    visible(newVal) {
      if (newVal) {
        this.visibleLocal = true;
      }
    }
  },
  computed: {
    deltaX() {
      return this.currentX - this.startX;
    },
    deltaY() {
      return this.currentY - this.startY;
    },
    style() {
      return {
        "--sidebar-max-width": this.maxWidth,
        "--sidebar-width": this.width,
        "--sidebar-max-height": this.maxHeight,
        "--sidebar-height": this.height,
        "--sidebar-box-shadow": this.boxShadow,
        "--sidebar-margin-top": this.marginTopComputed + "px",
        "--sidebar-translate-x": `${this.translateX || 0}px`,
        "--sidebar-translate-y": `${this.translateY || 0}px`
      };
    }
  },
  mounted() {
    this.adjustVertically();
    window.addEventListener("resize", this.adjustVertically);
    document.querySelector("body").style.overflow = "hidden";
  },
  methods: {
    adjustVertically() {
      if (!this.$refs.content) return;
      const content = this.$refs.content.offsetHeight,
        height = window.innerHeight,
        marginTop = parseInt(this.marginTop) || 100;
      if (this.side === "center") {
        this.marginTopComputed =
          content > height - 40 ? 20 : (height - content) / 2;
        this.fullscreenComputed =
          window.innerWidth <= parseInt(this.width) && this.fullscreen;
      }
      if (this.side === "bottom") {
        this.marginTopComputed =
          content > height - marginTop ? marginTop : height - content;
      }
    },
    sidebarClick(e) {
      if (this.side === "center") this.visibleLocal = null;
      if (this.side === "bottom") this.visibleLocal = null;
    },
    setScrolling(bool) {
      this.isScrolling = bool;
    },
    emitVisible() {
      document.querySelector("body").style.overflow = "";
      /**
       * Sends `false` when closing the sidebar.
       * @type {Event}
       */
      this.$emit("visible", false);
    },
    close(e) {
      this.visibleLocal = null;
      this.$refs.overlay.style["pointer-events"] = "none";
      if (e.clientX && e.clientY) {
        const doc = document.elementFromPoint(e.clientX, e.clientY);
        if (doc && doc.click) doc.click();
      }
    },
    touchstart(e) {
      this.$refs.sidebar.style.transition = "";
      this.currentX = this.startX = e.changedTouches[0].clientX;
      this.currentY = this.startY = e.changedTouches[0].clientY;
    },
    touchmove(e) {
      this.currentX = e.changedTouches[0].clientX;
      this.currentY = e.changedTouches[0].clientY;
      if (this.side === "left" && !this.isScrolling) {
        this.translateX = this.deltaX > 0 ? 0 : this.deltaX;
      }
      if (this.side === "right" && !this.isScrolling) {
        this.translateX = this.deltaX < 0 ? 0 : this.deltaX;
      }
    },
    touchend(e) {
      const overThresholdX =
          Math.abs((this.deltaX * 100) / window.screen.width) > 25,
        left = this.side === "left",
        right = this.side === "right";
      if (left) {
        this.translateX = this.deltaX > 0 ? 0 : this.deltaX;
      }
      if (right) {
        this.translateX = this.deltaX < 0 ? 0 : this.deltaX;
      }
      if (overThresholdX && !this.isScrolling && (left || right)) {
        this.close(e);
      } else if (this.$refs.sidebar) {
        this.$refs.sidebar.style.transition = "all .5s";
      }
      this.startX = null;
      this.startY = null;
      this.currentX = null;
      this.currentY = null;
      this.translateX = null;
      this.isScrolling = null;
    }
  }
};
</script>
