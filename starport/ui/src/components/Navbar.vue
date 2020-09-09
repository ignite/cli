<template>
  <nav>
    <div class="nav__logo">
      <slot />
    </div>
    <div class="nav__main">
      <div class="tabs">
        <router-link
          class="tabs__item"
          to="/"
        >
          Welcome
        </router-link>
        <router-link
          class="tabs__item"
          to="/blocks"
        >
          Blocks
        </router-link>
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
  text-decoration: none;

  font-size: 1.5rem;
  font-weight: 400;
  color: var(--c-txt-secondary);
}
.tabs__item:not(:last-child) {
  margin-bottom: 1rem;
}
.tabs__item.router-link-exact-active {
  font-weight: 600;
  color: var(--c-txt-primary);
}

</style>
