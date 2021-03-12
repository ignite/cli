import { txClient, queryClient } from './module';
// @ts-ignore
import { SpVuexError } from '@starport/vuex';
import { Equivocation } from "./module/types/cosmos/evidence/v1beta1/evidence";
async function initTxClient(vuexGetters) {
    return await txClient(vuexGetters['common/wallet/signer'], {
        addr: vuexGetters['common/env/apiTendermint']
    });
}
async function initQueryClient(vuexGetters) {
    return await queryClient({
        addr: vuexGetters['common/env/apiCosmos']
    });
}
function getStructure(template) {
    let structure = { fields: [] };
    for (const [key, value] of Object.entries(template)) {
        let field = {};
        field.name = key;
        field.type = typeof value;
        structure.fields.push(field);
    }
    return structure;
}
const getDefaultState = () => {
    return {
        Evidence: {},
        AllEvidence: {},
        _Structure: {
            Equivocation: getStructure(Equivocation.fromPartial({})),
        },
        _Subscriptions: new Set(),
    };
};
// initial state
const state = getDefaultState();
export default {
    namespaced: true,
    state,
    mutations: {
        RESET_STATE(state) {
            Object.assign(state, getDefaultState());
        },
        QUERY(state, { query, key, value }) {
            state[query][JSON.stringify(key)] = value;
        },
        SUBSCRIBE(state, subscription) {
            state._Subscriptions.add(subscription);
        },
        UNSUBSCRIBE(state, subscription) {
            state._Subscriptions.delete(subscription);
        }
    },
    getters: {
        getEvidence: (state) => (params = {}) => {
            return state.Evidence[JSON.stringify(params)] ?? {};
        },
        getAllEvidence: (state) => (params = {}) => {
            return state.AllEvidence[JSON.stringify(params)] ?? {};
        },
        getTypeStructure: (state) => (type) => {
            return state._Structure[type].fields;
        }
    },
    actions: {
        init({ dispatch, rootGetters }) {
            console.log('init');
            if (rootGetters['common/env/client']) {
                rootGetters['common/env/client'].on('newblock', () => {
                    dispatch('StoreUpdate');
                });
            }
        },
        resetState({ commit }) {
            commit('RESET_STATE');
        },
        unsubscribe({ commit }, subscription) {
            commit('UNSUBSCRIBE', subscription);
        },
        async StoreUpdate({ state, dispatch }) {
            state._Subscriptions.forEach((subscription) => {
                dispatch(subscription.action, subscription.payload);
            });
        },
        async QueryEvidence({ commit, rootGetters, getters, state }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryEvidence.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryEvidence.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'Evidence', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryEvidence', payload: { all, ...key } });
                return getters['getEvidence'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryEvidence', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryAllEvidence({ commit, rootGetters, getters, state }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryAllEvidence.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryAllEvidence.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'AllEvidence', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryAllEvidence', payload: { all, ...key } });
                return getters['getAllEvidence'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryAllEvidence', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async sendMsgSubmitEvidence({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgSubmitEvidence(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgSubmitEvidence:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgSubmitEvidence:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async MsgSubmitEvidence({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgSubmitEvidence(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgSubmitEvidence:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgSubmitEvidence:Create', 'Could not create message.');
                }
            }
        },
    }
};
