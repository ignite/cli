import { txClient, queryClient } from './module';
import { HistoricalInfo } from "./module/types/cosmos/staking/v1beta1/staking";
import { CommissionRates } from "./module/types/cosmos/staking/v1beta1/staking";
import { Commission } from "./module/types/cosmos/staking/v1beta1/staking";
import { Description } from "./module/types/cosmos/staking/v1beta1/staking";
import { Validator } from "./module/types/cosmos/staking/v1beta1/staking";
import { ValAddresses } from "./module/types/cosmos/staking/v1beta1/staking";
import { DVPair } from "./module/types/cosmos/staking/v1beta1/staking";
import { DVPairs } from "./module/types/cosmos/staking/v1beta1/staking";
import { DVVTriplet } from "./module/types/cosmos/staking/v1beta1/staking";
import { DVVTriplets } from "./module/types/cosmos/staking/v1beta1/staking";
import { Delegation } from "./module/types/cosmos/staking/v1beta1/staking";
import { UnbondingDelegation } from "./module/types/cosmos/staking/v1beta1/staking";
import { UnbondingDelegationEntry } from "./module/types/cosmos/staking/v1beta1/staking";
import { RedelegationEntry } from "./module/types/cosmos/staking/v1beta1/staking";
import { Redelegation } from "./module/types/cosmos/staking/v1beta1/staking";
import { Params } from "./module/types/cosmos/staking/v1beta1/staking";
import { DelegationResponse } from "./module/types/cosmos/staking/v1beta1/staking";
import { RedelegationEntryResponse } from "./module/types/cosmos/staking/v1beta1/staking";
import { RedelegationResponse } from "./module/types/cosmos/staking/v1beta1/staking";
import { Pool } from "./module/types/cosmos/staking/v1beta1/staking";
import { LastValidatorPower } from "./module/types/cosmos/staking/v1beta1/genesis";
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
        getValidators: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getValidator: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getValidatorDelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getValidatorUnbondingDelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegation: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getUnbondingDelegation: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegatorDelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegatorUnbondingDelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getRedelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegatorValidators: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegatorValidator: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getHistoricalInfo: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getPool: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getParams: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        _Structure: {
            HistoricalInfo: getStructure(HistoricalInfo.fromPartial({})),
            CommissionRates: getStructure(CommissionRates.fromPartial({})),
            Commission: getStructure(Commission.fromPartial({})),
            Description: getStructure(Description.fromPartial({})),
            Validator: getStructure(Validator.fromPartial({})),
            ValAddresses: getStructure(ValAddresses.fromPartial({})),
            DVPair: getStructure(DVPair.fromPartial({})),
            DVPairs: getStructure(DVPairs.fromPartial({})),
            DVVTriplet: getStructure(DVVTriplet.fromPartial({})),
            DVVTriplets: getStructure(DVVTriplets.fromPartial({})),
            Delegation: getStructure(Delegation.fromPartial({})),
            UnbondingDelegation: getStructure(UnbondingDelegation.fromPartial({})),
            UnbondingDelegationEntry: getStructure(UnbondingDelegationEntry.fromPartial({})),
            RedelegationEntry: getStructure(RedelegationEntry.fromPartial({})),
            Redelegation: getStructure(Redelegation.fromPartial({})),
            Params: getStructure(Params.fromPartial({})),
            DelegationResponse: getStructure(DelegationResponse.fromPartial({})),
            RedelegationEntryResponse: getStructure(RedelegationEntryResponse.fromPartial({})),
            RedelegationResponse: getStructure(RedelegationResponse.fromPartial({})),
            Pool: getStructure(Pool.fromPartial({})),
            LastValidatorPower: getStructure(LastValidatorPower.fromPartial({})),
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
        getValidators: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getValidator: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getValidatorDelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getValidatorUnbondingDelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegation: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getUnbondingDelegation: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegatorDelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegatorUnbondingDelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getRedelegations: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegatorValidators: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getDelegatorValidator: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getHistoricalInfo: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getPool: (state) => (params = {}) => {
            return state.Post[JSON.stringify(params)] ?? {};
        },
        getParams: (state) => (params = {}) => {
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
        async QueryValidators({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryValidators.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryValidator({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryValidator.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryValidatorDelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryValidatorDelegations.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryValidatorUnbondingDelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryValidatorUnbondingDelegations.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryDelegation({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryDelegation.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryUnbondingDelegation({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryUnbondingDelegation.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryDelegatorDelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryDelegatorDelegations.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryDelegatorUnbondingDelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryDelegatorUnbondingDelegations.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryRedelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryRedelegations.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryDelegatorValidators({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryDelegatorValidators.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryDelegatorValidator({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryDelegatorValidator.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryHistoricalInfo({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryHistoricalInfo.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
        },
        async QueryPool({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryPool.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Post', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPost', payload: key });
            }
            catch (e) {
                console.log('Query Failed: API node unavailable');
            }
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
        async MsgBeginRedelegate({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgBeginRedelegate(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                throw 'Failed to broadcast transaction: ' + e;
            }
        },
        async MsgEditValidator({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgEditValidator(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                throw 'Failed to broadcast transaction: ' + e;
            }
        },
        async MsgDelegate({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgDelegate(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                throw 'Failed to broadcast transaction: ' + e;
            }
        },
        async MsgUndelegate({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgUndelegate(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                throw 'Failed to broadcast transaction: ' + e;
            }
        },
        async MsgCreateValidator({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgCreateValidator(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                throw 'Failed to broadcast transaction: ' + e;
            }
        },
    }
};
