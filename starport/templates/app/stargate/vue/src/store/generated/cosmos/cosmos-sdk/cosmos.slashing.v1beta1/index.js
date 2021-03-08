import { txClient, queryClient } from './module';
// @ts-ignore
import { SpVuexError } from '@starport/vuex';
import { SigningInfo } from "./module/types/cosmos/slashing/v1beta1/genesis";
import { ValidatorMissedBlocks } from "./module/types/cosmos/slashing/v1beta1/genesis";
import { MissedBlock } from "./module/types/cosmos/slashing/v1beta1/genesis";
import { ValidatorSigningInfo } from "./module/types/cosmos/slashing/v1beta1/slashing";
import { Params } from "./module/types/cosmos/slashing/v1beta1/slashing";
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
        Params: {},
        SigningInfo: {},
        SigningInfos: {},
        _Structure: {
            SigningInfo: getStructure(SigningInfo.fromPartial({})),
            ValidatorMissedBlocks: getStructure(ValidatorMissedBlocks.fromPartial({})),
            MissedBlock: getStructure(MissedBlock.fromPartial({})),
            ValidatorSigningInfo: getStructure(ValidatorSigningInfo.fromPartial({})),
            Params: getStructure(Params.fromPartial({})),
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
        getParams: (state) => (params = {}) => {
            return state.Params[JSON.stringify(params)] ?? {};
        },
        getSigningInfo: (state) => (params = {}) => {
            return state.SigningInfo[JSON.stringify(params)] ?? {};
        },
        getSigningInfos: (state) => (params = {}) => {
            return state.SigningInfos[JSON.stringify(params)] ?? {};
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
        async QueryParams({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryParams.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Params', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryParams', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryParams', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QuerySigningInfo({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).querySigningInfo.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'SigningInfo', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QuerySigningInfo', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QuerySigningInfo', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QuerySigningInfos({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).querySigningInfos.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'SigningInfos', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QuerySigningInfos', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QuerySigningInfos', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async sendMsgUnjail({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgUnjail(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgUnjail:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgUnjail:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async MsgUnjail({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgUnjail(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgUnjail:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgUnjail:Create', 'Could not create message.');
                }
            }
        },
    }
};
