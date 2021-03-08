import { txClient, queryClient } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { FungibleTokenPacketData } from "./module/types/ibc/applications/transfer/v1/transfer"
import { DenomTrace } from "./module/types/ibc/applications/transfer/v1/transfer"
import { Params } from "./module/types/ibc/applications/transfer/v1/transfer"


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
        DenomTrace: {},
        DenomTraces: {},
        Params: {},
        
        _Structure: {
            FungibleTokenPacketData: getStructure(FungibleTokenPacketData.fromPartial({})),
            DenomTrace: getStructure(DenomTrace.fromPartial({})),
            Params: getStructure(Params.fromPartial({})),
            
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
        getDenomTrace: (state) => (params = {}) => {
			return state.DenomTrace[JSON.stringify(params)] ?? {}
		},
        getDenomTraces: (state) => (params = {}) => {
			return state.DenomTraces[JSON.stringify(params)] ?? {}
		},
        getParams: (state) => (params = {}) => {
			return state.Params[JSON.stringify(params)] ?? {}
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
		async QueryDenomTrace({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDenomTrace.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DenomTrace', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDenomTrace', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDenomTrace', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDenomTraces({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDenomTraces.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DenomTraces', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDenomTraces', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDenomTraces', 'API Node Unavailable. Could not perform query.'))
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
		
		async sendMsgTransfer({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgTransfer(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgTransfer:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgTransfer:Send', 'Could not broadcast Tx.')
				}
			}
		},
		
		async MsgTransfer({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgTransfer(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgTransfer:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgTransfer:Create', 'Could not create message.')
				}
			}
		},
		
	}
}
