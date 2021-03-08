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
        async QueryEvidence({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryEvidence.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Evidence', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryEvidence', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryEvidence', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QueryAllEvidence({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryAllEvidence.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'AllEvidence', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryAllEvidence', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryAllEvidence', 'API Node Unavailable. Could not perform query.'));
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
