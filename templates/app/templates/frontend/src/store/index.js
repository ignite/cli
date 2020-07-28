import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import app from "./app.js";
import {
  Secp256k1Pen,
  SigningCosmosClient,
  makeSignBytes,
} from "@cosmjs/sdk38";

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
    client: null,
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
    clientUpdate(state, { client }) {
      state.client = client;
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
        const wallet = await Secp256k1Pen.fromMnemonic(mnemonic);
        const address = wallet.address("cosmos");
        const url = `${API}/auth/accounts/${address}`;
        const acc = (await axios.get(url)).data;
        if (acc.result.value.address === address) {
          const account = acc.result.value;
          const client = new SigningCosmosClient(API, address, wallet);
          commit("accountUpdate", { account });
          commit("walletUpdate", { wallet });
          commit("clientUpdate", { client });
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
    async accountUpdate({ state, commit }) {
      const url = `${API}/auth/accounts/${state.client.senderAddress}`;
      const acc = (await axios.get(url)).data;
      const account = acc.result.value;
      commit("accountUpdate", { account });
    },
    async entitySubmit({ state }, { type, body }) {
      const accountURL = `${API}/auth/accounts/${state.client.senderAddress}`;
      const account = (await axios.get(accountURL)).data.result.value;
      const accountNumber = account.account_number;
      const sequence = account.sequence;
      const address = state.client.senderAddress;
      const chain_id = await state.client.getChainId();
      const req = {
        base_req: { chain_id, from: address },
        creator: address,
        ...body,
      };
      const { data } = await axios.post(`${API}/${chain_id}/${type}`, req);
      const { msg, fee, memo } = data.value;
      const signBytes = makeSignBytes(
        msg,
        fee,
        chain_id,
        memo,
        `${accountNumber}`,
        `${sequence}`
      );
      const signatures = [await state.wallet.sign(signBytes)];
      const signedTx = { msg: msg, fee, memo, signatures };
      return await state.client.postTx(signedTx);
    },
  },
});
