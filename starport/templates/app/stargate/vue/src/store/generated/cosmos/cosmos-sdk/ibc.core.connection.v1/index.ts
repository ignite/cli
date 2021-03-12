import { txClient, queryClient } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { ConnectionEnd } from "./module/types/ibc/core/connection/v1/connection"
import { IdentifiedConnection } from "./module/types/ibc/core/connection/v1/connection"
import { Counterparty } from "./module/types/ibc/core/connection/v1/connection"
import { ClientPaths } from "./module/types/ibc/core/connection/v1/connection"
import { ConnectionPaths } from "./module/types/ibc/core/connection/v1/connection"
import { Version } from "./module/types/ibc/core/connection/v1/connection"


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
        Connection: {},
        Connections: {},
        ClientConnections: {},
        ConnectionClientState: {},
        ConnectionConsensusState: {},
        
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
			return state.Connection[JSON.stringify(params)] ?? {}
		},
        getConnections: (state) => (params = {}) => {
			return state.Connections[JSON.stringify(params)] ?? {}
		},
        getClientConnections: (state) => (params = {}) => {
			return state.ClientConnections[JSON.stringify(params)] ?? {}
		},
        getConnectionClientState: (state) => (params = {}) => {
			return state.ConnectionClientState[JSON.stringify(params)] ?? {}
		},
        getConnectionConsensusState: (state) => (params = {}) => {
			return state.ConnectionConsensusState[JSON.stringify(params)] ?? {}
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
		async QueryConnection({ commit, rootGetters, getters, state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryConnection.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryConnection.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'Connection', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryConnection', payload: { all, ...key} })
				return getters['getConnection'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryConnection', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryConnections({ commit, rootGetters, getters, state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryConnections.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryConnections.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'Connections', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryConnections', payload: { all, ...key} })
				return getters['getConnections'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryConnections', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryClientConnections({ commit, rootGetters, getters, state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryClientConnections.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryClientConnections.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ClientConnections', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryClientConnections', payload: { all, ...key} })
				return getters['getClientConnections'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryClientConnections', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryConnectionClientState({ commit, rootGetters, getters, state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryConnectionClientState.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryConnectionClientState.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ConnectionClientState', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryConnectionClientState', payload: { all, ...key} })
				return getters['getConnectionClientState'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryConnectionClientState', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryConnectionConsensusState({ commit, rootGetters, getters, state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryConnectionConsensusState.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryConnectionConsensusState.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ConnectionConsensusState', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryConnectionConsensusState', payload: { all, ...key} })
				return getters['getConnectionConsensusState'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryConnectionConsensusState', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		
		async sendMsgConnectionOpenAck({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenAck(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgConnectionOpenAck:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgConnectionOpenAck:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgConnectionOpenInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenInit(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgConnectionOpenInit:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgConnectionOpenInit:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgConnectionOpenConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenConfirm(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgConnectionOpenConfirm:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgConnectionOpenConfirm:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgConnectionOpenTry({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenTry(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgConnectionOpenTry:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgConnectionOpenTry:Send', 'Could not broadcast Tx.')
				}
			}
		},
		
		async MsgConnectionOpenAck({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenAck(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgConnectionOpenAck:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgConnectionOpenAck:Create', 'Could not create message.')
				}
			}
		},
		async MsgConnectionOpenInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenInit(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgConnectionOpenInit:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgConnectionOpenInit:Create', 'Could not create message.')
				}
			}
		},
		async MsgConnectionOpenConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenConfirm(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgConnectionOpenConfirm:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgConnectionOpenConfirm:Create', 'Could not create message.')
				}
			}
		},
		async MsgConnectionOpenTry({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgConnectionOpenTry(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgConnectionOpenTry:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgConnectionOpenTry:Create', 'Could not create message.')
				}
			}
		},
		
	}
}
