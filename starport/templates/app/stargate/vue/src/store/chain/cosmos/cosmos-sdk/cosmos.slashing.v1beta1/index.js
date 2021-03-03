import { txClient, queryClient } from './module';
import { ValidatorSigningInfo } from "./module/types/cosmos/slashing/v1beta1/slashing";
import { Params } from "./module/types/cosmos/slashing/v1beta1/slashing";
import { SigningInfo } from "./module/types/cosmos/slashing/v1beta1/genesis";
import { ValidatorMissedBlocks } from "./module/types/cosmos/slashing/v1beta1/genesis";
import { MissedBlock } from "./module/types/cosmos/slashing/v1beta1/genesis";
async function initTxClient(vuexGetters) {
    return await txClient(vuexGetters['chain/common/wallet/signer'], {
        addr: vuexGetters['chain/common/env/apiTendermint']
    });
}
async function initQueryClient(vuexGetters) {
    return await queryClient({
        addr: vuexGetters['chain/common/env/apiCosmos']
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
        getParams: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getSigningInfo: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getSigningInfos: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        _Structure: {
            ValidatorSigningInfo: getStructure(ValidatorSigningInfo.fromPartial({})),
            Params: getStructure(Params.fromPartial({})),
            SigningInfo: getStructure(SigningInfo.fromPartial({})),
            ValidatorMissedBlocks: getStructure(ValidatorMissedBlocks.fromPartial({})),
            MissedBlock: getStructure(MissedBlock.fromPartial({})),
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
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getSigningInfo: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getSigningInfos: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getTypeStructure: (state) => (type) => {
            return state._Structure[type].fields;
        }
    },
    actions: {
        init({ dispatch, rootGetters }) {
            console.log('init');
            if (rootGetters['chain/common/env/client']) {
                rootGetters['chain/common/env/client'].on('newblock', () => {
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
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QuerySigningInfo({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).querySigningInfo.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QuerySigningInfos({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).querySigningInfos.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async MsgUnjail({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgUnjail(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                throw 'Failed to broadcast transaction: ' + e;
            }
        },
    }
};
