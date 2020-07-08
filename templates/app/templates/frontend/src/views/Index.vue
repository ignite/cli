<template>
  <div>
    <div class="container">
      <div class="narrow">
        <div class="h1">Blog</div>
        <wallet/>
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
@media screen and (max-width: 980px) {
  .narrow {
    padding: 0;
  }
}
</style>

<script>
import Type from "@/components/Type.vue";
import Wallet from "@/components/Wallet.vue";

export default {
  components: {
    Type,
    Wallet
  },
  computed: {
    typeList() {
      return this.$store.state.app.types || [];
    }
  }
};
</script>
