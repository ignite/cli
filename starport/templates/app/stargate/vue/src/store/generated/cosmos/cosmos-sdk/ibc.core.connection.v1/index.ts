import { txClient, queryClient } from './module'

import { ConnectionEnd } from "./module/types/ibc/core/connection/v1/connection"
import { IdentifiedConnection } from "./module/types/ibc/core/connection/v1/connection"
import { Counterparty } from "./module/types/ibc/core/connection/v1/connection"
import { ClientPaths } from "./module/types/ibc/core/connection/v1/connection"
import { ConnectionPaths } from "./module/types/ibc/core/connection/v1/connection"
import { Version } from "./module/types/ibc/core/connection/v1/connection"


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
        getConnection: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConnections: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getClientConnections: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConnectionClientState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConnectionConsensusState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        
        _Structure: {
            ConnectionEnd: getStructure(ConnectionEnd.fromPartial({})),
            IdentifiedConnection: getStructure(IdentifiedConnection.fromPartial({})),
            Counterparty: getStructure(Counterparty.fromPartial({})),
            ClientPaths: getStructure(ClientPaths.fromPartial({})),
            ConnectionPaths: getStructure(ConnectionPaths.fromPartial({})),
            Version: getStructure(Version.fromPartial({})),
            
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
        getConnection: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConnections: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getClientConnections: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConnectionClientState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConnectionConsensusState: (state) => (params = {}) => {
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
        async QueryConnection({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryConnection.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryConnections({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryConnections.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryClientConnections({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryClientConnections.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryConnectionClientState({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryConnectionClientState.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryConnectionConsensusState({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryConnectionConsensusState.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        
        async MsgConnectionOpenTry({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenTry(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgConnectionOpenInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenInit(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgConnectionOpenAck({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenAck(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgConnectionOpenConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenConfirm(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        
	}
}
