import Vue from "vue";
import Vuex from "vuex";

import explorer from "@tendermint/vue/src/store/common/explorer/explorer";

Vue.use(Vuex);

export default new Vuex.Store({
  namespaced: true,
  state: {},
  getters: {},
  mutations: {},
  actions: {},
  modules: {
    explorer,
  },
});
