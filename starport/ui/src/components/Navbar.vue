<template>
  <nav>
    <div class="nav__logo">
      <slot />
    </div>
    <div class="nav__main">
      <div class="tabs">
        <button class="tabs__item -is-active">Welcome</button>
        <button class="tabs__item">Blocks</button>
      </div>
    </div>
  </nav>
</template>

<style scoped>

nav {
  position: sticky;
  top: 0;
}

.tabs__item {
  display: block;
  background: transparent;
  border: none;

  font-size: 1.8rem;
  font-weight: 400;
}
.tabs__item:not(:last-child) {
  margin-bottom: 1rem;
}
.tabs__item.-is-active {
  font-weight: 600;
}

</style>
