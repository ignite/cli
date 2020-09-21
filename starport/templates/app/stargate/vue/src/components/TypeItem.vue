<template>
  <div>
    <app-text type="h2">List of {{ value.type }} items</app-text>
    <div class="item" v-for="instance in instanceList" :key="instance.id">
      <div class="item__field" v-for="(value, key) in instance" :key="key">
        <div class="item__field__key">{{ key }}:</div>
        <div class="item__field__value">
          {{ value }}
        </div>
      </div>
    </div>
    <div class="card__empty" v-if="instanceList.length < 1">
      There are no {{ value.type }} items yet. Create one using the form below.
    </div>
    <app-text type="h2">New {{ value.type }}</app-text>
    <div v-for="field in value.fields" :key="field">
      <app-input
        v-model="fields[field]"
        type="text"
        :placeholder="title(field)"
        :disabled="flight"
      />
    </div>
    <button
      :class="['button', `button__valid__${!!valid && !flight && hasAddress}`]"
      @click="submit"
    >
      Create {{ value.type }}
      <div class="button__label" v-if="flight">
        <div class="button__label__icon">
          <icon-refresh />
        </div>
        Sending transaction...
      </div>
    </button>
  </div>
</template>

<style scoped>
button {
  background: none;
  border: none;
  color: rgba(0, 125, 255);
  padding: 0;
  font-size: inherit;
  font-weight: 800;
  font-family: inherit;
  text-transform: uppercase;
  margin-top: 0.5rem;
  cursor: pointer;
  transition: opacity 0.1s;
  letter-spacing: 0.03em;
  transition: color 0.25s;
  display: inline-flex;
  align-items: center;
}
.item {
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.1);
  margin-bottom: 1rem;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow: hidden;
}
.item__field {
  display: grid;
  line-height: 1.5;
  grid-template-columns: 15% 1fr;
  grid-template-rows: 1fr;
  word-break: break-all;
}
.item__field__key {
  color: rgba(0, 0, 0, 0.25);
  word-break: keep-all;
  overflow: hidden;
}
button:focus {
  opacity: 0.85;
  outline: none;
}
.button.button__valid__true:active {
  opacity: 0.65;
}
.button__label {
  display: inline-flex;
  align-items: center;
}
.button__label__icon {
  height: 1em;
  width: 1em;
  margin: 0 0.5em 0 0.5em;
  fill: rgba(0, 0, 0, 0.25);
  animation: rotate linear 4s infinite;
}
.button.button__valid__false {
  color: rgba(0, 0, 0, 0.25);
  cursor: not-allowed;
}
.card__empty {
  margin-bottom: 1rem;
  border: 1px dashed rgba(0, 0, 0, 0.1);
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  border-radius: 8px;
  color: rgba(0, 0, 0, 0.25);
  text-align: center;
  min-height: 8rem;
}
@keyframes rotate {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(-360deg);
  }
}
@media screen and (max-width: 980px) {
  .narrow {
    padding: 0;
  }
}
</style>

<script>
export default {
  props: ["value"],
  data: function () {
    return {
      fields: {},
      flight: false,
    };
  },
  created() {
    (this.value.fields || []).forEach((field) => {
      this.$set(this.fields, field, "");
    });
  },
  computed: {
    hasAddress() {
      return !!this.$store.state.account.address;
    },
    instanceList() {
      return this.$store.state.data[this.value.type] || [];
    },
    valid() {
      return Object.values(this.fields).every((el) => {
        return el.trim().length > 0;
      });
    },
  },
  methods: {
    title(string) {
      return string.charAt(0).toUpperCase() + string.slice(1);
    },
    async submit() {
      if (this.valid && !this.flight && this.hasAddress) {
        this.flight = true;
        const payload = { type: this.value.type, body: this.fields };
        await this.$store.dispatch("entitySubmit", payload);
        await this.$store.dispatch("entityFetch", payload);
        this.flight = false;
        Object.keys(this.fields).forEach((f) => {
          this.fields[f] = "";
        });
      }
    },
  },
};
</script>
