import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import {
  createWalletFromMnemonic,
  signTx,
  createBroadcastTx,
} from "@tendermint/sig";
import * as bip39 from "bip39";
import app from "./app.js";

Vue.use(Vuex);

const RPC = "https://localhost:26657";
const API = "http://localhost:8080";

export default new Vuex.Store({
  state: {
    app,
    account: {},
    chain_id: "",
    data: {},
  },
  mutations: {
    accountCreate(state, account) {
      state.account = account;
    },
    chainIdSet(state, { chain_id }) {
      state.chain_id = chain_id;
    },
    instanceSet(state, { instanceList, type }) {
      const updated = {};
      updated[type] = instanceList;
      state.data = { ...state.data, ...updated };
    },
  },
  actions: {
    async init({ dispatch, state }) {
      dispatch("accountCreate");
      await dispatch("chainIdFetch");
      state.app.types.forEach(({ type }) => {
        dispatch("instanceListFetch", { type });
      });
    },
    async chainIdFetch({ commit }) {
      const node_info = (await axios.get(`${API}/node_info`)).data.node_info;
      commit("chainIdSet", { chain_id: node_info.network });
    },
    async accountCreate({ commit }) {
      const mnemonic = bip39.generateMnemonic();
      const wallet = createWalletFromMnemonic(mnemonic);
      await axios.post(`${API}/faucet`, { address: wallet.address });
      commit("accountCreate", { wallet, mnemonic });
    },
    async instanceListFetch({ state, commit }, { type }) {
      const { chain_id } = state;
      const url = `${API}/${chain_id}/${type}`;
      const instanceList = (await axios.get(url)).data.result;
      commit("instanceSet", { type, instanceList });
    },
    instanceCreate({ state, dispatch, commit }, { fields, type }) {
      return new Promise((resolve, reject) => {
        const wallet = state.account.wallet;
        const chain_id = state.chain_id;
        axios.get(`${API}/auth/accounts/${wallet.address}`).then(({ data }) => {
          const account = data.result.value;
          const meta = {
            sequence: `${account.sequence}`,
            account_number: `${account.account_number}`,
            chain_id,
          };
          const req = {
            base_req: {
              chain_id,
              from: wallet.address,
            },
            creator: wallet.address,
            ...fields,
          };
          axios.post(`${API}/${chain_id}/${type}`, req).then(({ data }) => {
            const tx = data.value;
            const stdTx = signTx(tx, meta, wallet);
            const txBroadcast = createBroadcastTx(stdTx, "block");
            const params = {
              headers: {
                "Content-Type": "application/json",
              },
            };
            axios.post(`${API}/txs`, txBroadcast, params).then(() => {
              this.dispatch("instanceListFetch", { type });
              resolve(true);
            });
          });
        });
      });
    },
  },
});
