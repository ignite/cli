<template>
  <div>
    <input
      type="text"
      :placeholder="placeholder"
      :value="value"
      @input="input"
      :disabled="disabled"
    />
  </div>
</template>

<style scoped>
input {
  border: none;
  font-size: inherit;
  padding: 0.75rem 1rem;
  margin-bottom: 0.5rem;
  width: 100%;
  max-width: 500px;
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.1);
  border-radius: 6px;
  box-sizing: border-box;
  background: rgba(0, 0, 0, 0);
  font-family: inherit;
}
input:focus {
  outline: none;
  box-shadow: 0 0 0 2.5px #d5e7fe, inset 0 0 0 1px #a5cafd;
}
input:disabled {
  color: rgba(0, 0, 0, 0.25);
}
</style>

<script>
export default {
  props: {
    value: {
      default: "",
    },
    placeholder: {
      default: "",
    },
    disabled: {
      default: false,
    },
  },
  methods: {
    input(event) {
      this.$emit("input", event.target.value);
    },
  },
};
</script>
