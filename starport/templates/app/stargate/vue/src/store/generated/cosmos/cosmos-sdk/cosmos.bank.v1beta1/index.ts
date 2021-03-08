import { txClient, queryClient } from './module'

import { Params } from "./module/types/cosmos/bank/v1beta1/bank"
import { SendEnabled } from "./module/types/cosmos/bank/v1beta1/bank"
import { Input } from "./module/types/cosmos/bank/v1beta1/bank"
import { Output } from "./module/types/cosmos/bank/v1beta1/bank"
import { Supply } from "./module/types/cosmos/bank/v1beta1/bank"
import { DenomUnit } from "./module/types/cosmos/bank/v1beta1/bank"
import { Metadata } from "./module/types/cosmos/bank/v1beta1/bank"
import { Balance } from "./module/types/cosmos/bank/v1beta1/genesis"


async function initTxClient(vuexGetters) {
	return await txClient(vuexGetters['chain/common/wallet/signer'], {
		addr: vuexGetters['chain/common/env/apiTendermint']
	})
}

async function initQueryClient(vuexGetters) {
	return await queryClient({
		addr: vuexGetters['chain/common/env/apiCosmos']
	})
}

function getStructure(template) {
	let structure = { fields: [] }
	for (const [key, value] of Object.entries(template)) {
		let field: any = {}
		field.name = key
		field.type = typeof value
		structure.fields.push(field)
	}
	return structure
}

const getDefaultState = () => {
	return {
        getBalance: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getAllBalances: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getTotalSupply: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getSupplyOf: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getParams: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDenomMetadata: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDenomsMetadata: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        
        _Structure: {
            Params: getStructure(Params.fromPartial({})),
            SendEnabled: getStructure(SendEnabled.fromPartial({})),
            Input: getStructure(Input.fromPartial({})),
            Output: getStructure(Output.fromPartial({})),
            Supply: getStructure(Supply.fromPartial({})),
            DenomUnit: getStructure(DenomUnit.fromPartial({})),
            Metadata: getStructure(Metadata.fromPartial({})),
            Balance: getStructure(Balance.fromPartial({})),
            
		},
		_Subscriptions: new Set(),
	}
}

// initial state
const state = getDefaultState()

export default {
	namespaced: true,
	state,
	mutations: {
		RESET_STATE(state) {
			Object.assign(state, getDefaultState())
		},
		QUERY(state, { query, key, value }) {
			state[query][JSON.stringify(key)] = value
		},
		SUBSCRIBE(state, subscription) {
			state._Subscriptions.add(subscription)
		},
		UNSUBSCRIBE(state, subscription) {
			state._Subscriptions.delete(subscription)
		}
	},
	getters: {
        getBalance: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getAllBalances: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getTotalSupply: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getSupplyOf: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getParams: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDenomMetadata: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDenomsMetadata: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('init')
			if (rootGetters['chain/common/env/client']) {
				rootGetters['chain/common/env/client'].on('newblock', () => {
					dispatch('StoreUpdate')
				})
			}
		},
		resetState({ commit }) {
			commit('RESET_STATE')
		},
		unsubscribe({ commit }, subscription) {
			commit('UNSUBSCRIBE', subscription)
		},
		async StoreUpdate({ state, dispatch }) {
			state._Subscriptions.forEach((subscription) => {
				dispatch(subscription.action, subscription.payload)
			})
		},
        async QueryBalance({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryBalance.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryAllBalances({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryAllBalances.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryTotalSupply({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryTotalSupply.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QuerySupplyOf({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).querySupplyOf.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryParams({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryParams.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryDenomMetadata({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDenomMetadata.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryDenomsMetadata({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDenomsMetadata.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        
        async MsgSend({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSend(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgMultiSend({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgMultiSend(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        
	}
}
