import { txClient, queryClient } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { Params } from "./module/types/cosmos/bank/v1beta1/bank"
import { SendEnabled } from "./module/types/cosmos/bank/v1beta1/bank"
import { Input } from "./module/types/cosmos/bank/v1beta1/bank"
import { Output } from "./module/types/cosmos/bank/v1beta1/bank"
import { Supply } from "./module/types/cosmos/bank/v1beta1/bank"
import { DenomUnit } from "./module/types/cosmos/bank/v1beta1/bank"
import { Metadata } from "./module/types/cosmos/bank/v1beta1/bank"
import { Balance } from "./module/types/cosmos/bank/v1beta1/genesis"


async function initTxClient(vuexGetters) {
	return await txClient(vuexGetters['common/wallet/signer'], {
		addr: vuexGetters['common/env/apiTendermint']
	})
}

async function initQueryClient(vuexGetters) {
	return await queryClient({
		addr: vuexGetters['common/env/apiCosmos']
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
        Balance: {},
        AllBalances: {},
        TotalSupply: {},
        SupplyOf: {},
        Params: {},
        DenomMetadata: {},
        DenomsMetadata: {},
        
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
			return state.Balance[JSON.stringify(params)] ?? {}
		},
        getAllBalances: (state) => (params = {}) => {
			return state.AllBalances[JSON.stringify(params)] ?? {}
		},
        getTotalSupply: (state) => (params = {}) => {
			return state.TotalSupply[JSON.stringify(params)] ?? {}
		},
        getSupplyOf: (state) => (params = {}) => {
			return state.SupplyOf[JSON.stringify(params)] ?? {}
		},
        getParams: (state) => (params = {}) => {
			return state.Params[JSON.stringify(params)] ?? {}
		},
        getDenomMetadata: (state) => (params = {}) => {
			return state.DenomMetadata[JSON.stringify(params)] ?? {}
		},
        getDenomsMetadata: (state) => (params = {}) => {
			return state.DenomsMetadata[JSON.stringify(params)] ?? {}
		},
        
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('init')
			if (rootGetters['common/env/client']) {
				rootGetters['common/env/client'].on('newblock', () => {
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
				commit('QUERY', { query: 'Balance', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryBalance', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryBalance', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryAllBalances({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryAllBalances.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'AllBalances', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryAllBalances', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryAllBalances', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryTotalSupply({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryTotalSupply.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'TotalSupply', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryTotalSupply', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryTotalSupply', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QuerySupplyOf({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).querySupplyOf.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'SupplyOf', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QuerySupplyOf', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QuerySupplyOf', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryParams({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryParams.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Params', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryParams', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryParams', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDenomMetadata({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDenomMetadata.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DenomMetadata', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDenomMetadata', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDenomMetadata', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDenomsMetadata({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDenomsMetadata.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DenomsMetadata', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDenomsMetadata', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDenomsMetadata', 'API Node Unavailable. Could not perform query.'))
			}
		},
		
		async sendMsgMultiSend({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgMultiSend(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgMultiSend:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgMultiSend:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgSend({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSend(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSend:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSend:Send', 'Could not broadcast Tx.')
				}
			}
		},
		
		async MsgMultiSend({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgMultiSend(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgMultiSend:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgMultiSend:Create', 'Could not create message.')
				}
			}
		},
		async MsgSend({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSend(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSend:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSend:Create', 'Could not create message.')
				}
			}
		},
		
	}
}
