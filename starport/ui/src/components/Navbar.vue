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

.nav__main {
  margin-top: -1rem;
}

.tabs__item {
  display: block;
  text-decoration: none;

  font-size: 1.25rem;
  font-weight: 300;
  color: var(--c-txt-grey);

  padding: 0.85rem 0 0.85rem 1.5rem;
  border-radius: 8px;
  transition: background-color .3s;
}
.tabs__item:hover {
  background-color: var(--c-bg-third);
  transition: background-color .3s;
}
.tabs__item:not(:last-child) {
  margin-bottom: 0.5rem;
}
.tabs__item.router-link-exact-active {
  font-weight: 600;
  color: var(--c-txt-secondary);
  background-color: var(--c-bg-third);
}

</style>
