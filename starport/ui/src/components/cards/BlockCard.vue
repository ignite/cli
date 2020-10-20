<template>
  <div :class="['card']">
    <div class="card__top">
      <p class="card__title">{{title}}</p>
      <span class="card__note">{{note}}</span>
    </div>
    <div class="card__bottom">
      <slot />
    </div>
  </div>
</template>

<script>
export default {
  props: {
    title: { type: String | Number, required: true },
    note: { type: String, required: true }
  }
}
</script>

<style scoped>

.card {
  padding: 1.5rem 0;
}

.card__top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.card__bottom {
  font-size: 0.8125rem;
  text-align: left;
}
.card__bottom * {
  margin-top: 0.5rem;
}

.card__title {
  font-size: 1.3125rem;
  font-weight: var(--f-w-bold);
  color: var(--c-txt-third);
}

.card__note {
  font-size: 0.75rem;
  color: var(--c-txt-third);
}

</style>