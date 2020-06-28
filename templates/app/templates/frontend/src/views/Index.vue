<template>
  <div>
    <div class="container">
      <div class="narrow">
        <div class="h1">Blog</div>
        <div class="card">
          <div class="card__icon">
            <icon-user />
          </div>
          <div class="card__desc">account: {{ address }}</div>
        </div>
        <div class="type" v-for="type in typeList" :key="type.type">
          <type :value="type" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  padding: 4rem 1rem 12rem;
  max-width: 900px;
  width: 100%;
  margin-left: auto;
  margin-right: auto;
  box-sizing: border-box;
}
.narrow {
  padding-left: 40px;
  padding-right: 40px;
  box-sizing: border-box;
}
.h1 {
  font-size: 2rem;
  font-weight: 800;
  letter-spacing: 0.03em;
  margin-bottom: 2rem;
}
.type {
  margin-bottom: 3rem;
}
.card {
  background: #f5f5f7;
  padding: 1rem;
  margin-bottom: 1.5rem;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  overflow: hidden;
}
.card__icon {
  width: 2rem;
  height: 2rem;
  fill: rgba(0, 0, 0, 0.1);
  margin-right: 1rem;
  flex-shrink: 0;
}
.card__desc {
  color: rgba(0, 0, 0, 0.5);
  letter-spacing: 0.02em;
}
@media screen and (max-width: 980px) {
  .narrow {
    padding: 0;
  }
}
</style>

<script>
import Type from "@/components/Type.vue";
import IconUser from "@/components/IconUser.vue";

export default {
  components: {
    Type,
    IconUser,
  },
  computed: {
    address() {
      const { account } = this.$store.state;
      const wallet = account && account.wallet;
      const address = wallet && wallet.address;
      return address;
    },
    typeList() {
      return this.$store.state.app.types || [];
    },
  },
};
</script>
