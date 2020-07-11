import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import {
  createWalletFromMnemonic,
  signTx,
  createBroadcastTx,
} from "@tendermint/sig";
import app from "./app.js";

Vue.use(Vuex);

const RPC = "https://localhost:26657";
const API = "http://localhost:8080";

export default new Vuex.Store({
  state: {
    app,
    wallet: {},
    account: {},
    chain_id: "",
    data: {},
  },
  mutations: {
    accountUpdate(state, { account }) {
      state.account = account;
    },
    walletUpdate(state, { wallet }) {
      state.wallet = wallet;
    },
    chainIdSet(state, { chain_id }) {
      state.chain_id = chain_id;
    },
    entitySet(state, { type, body }) {
      const updated = {};
      updated[type] = body;
      state.data = { ...state.data, ...updated };
    },
  },
  actions: {
    async init({ dispatch, state }) {
      await dispatch("chainIdFetch");
      state.app.types.forEach(({ type }) => {
        dispatch("entityFetch", { type });
      });
    },
    async chainIdFetch({ commit }) {
      const node_info = (await axios.get(`${API}/node_info`)).data.node_info;
      commit("chainIdSet", { chain_id: node_info.network });
    },
    async accountSignIn({ commit }, { mnemonic }) {
      return new Promise(async (resolve, reject) => {
        const wallet = createWalletFromMnemonic(mnemonic);
        const url = `${API}/auth/accounts/${wallet.address}`;
        const acc = (await axios.get(url)).data;
        if (acc.result.value.address === wallet.address) {
          const account = acc.result.value;
          commit("accountUpdate", { account });
          commit("walletUpdate", { wallet });
          resolve(account);
        } else {
          reject("Account doesn't exist.");
        }
      });
    },
    async entityFetch({ state, commit }, { type }) {
      const { chain_id } = state;
      const url = `${API}/${chain_id}/${type}`;
      const body = (await axios.get(url)).data.result;
      commit("entitySet", { type, body });
    },
    async entitySubmit({ state }, { type, body }) {
      return new Promise((resolve, reject) => {
        const wallet = state.wallet;
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
            ...body,
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
              resolve(true);
            });
          });
        });
      });
    },
  },
});
