import { txClient, queryClient } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { GenesisMetadata } from "./module/types/ibc/core/client/v1/genesis"
import { IdentifiedGenesisMetadata } from "./module/types/ibc/core/client/v1/genesis"
import { IdentifiedClientState } from "./module/types/ibc/core/client/v1/client"
import { ConsensusStateWithHeight } from "./module/types/ibc/core/client/v1/client"
import { ClientConsensusStates } from "./module/types/ibc/core/client/v1/client"
import { ClientUpdateProposal } from "./module/types/ibc/core/client/v1/client"
import { Height } from "./module/types/ibc/core/client/v1/client"
import { Params } from "./module/types/ibc/core/client/v1/client"


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
        ClientState: {},
        ClientStates: {},
        ConsensusState: {},
        ConsensusStates: {},
        ClientParams: {},
        
        _Structure: {
            GenesisMetadata: getStructure(GenesisMetadata.fromPartial({})),
            IdentifiedGenesisMetadata: getStructure(IdentifiedGenesisMetadata.fromPartial({})),
            IdentifiedClientState: getStructure(IdentifiedClientState.fromPartial({})),
            ConsensusStateWithHeight: getStructure(ConsensusStateWithHeight.fromPartial({})),
            ClientConsensusStates: getStructure(ClientConsensusStates.fromPartial({})),
            ClientUpdateProposal: getStructure(ClientUpdateProposal.fromPartial({})),
            Height: getStructure(Height.fromPartial({})),
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
        getClientState: (state) => (params = {}) => {
			return state.ClientState[JSON.stringify(params)] ?? {}
		},
        getClientStates: (state) => (params = {}) => {
			return state.ClientStates[JSON.stringify(params)] ?? {}
		},
        getConsensusState: (state) => (params = {}) => {
			return state.ConsensusState[JSON.stringify(params)] ?? {}
		},
        getConsensusStates: (state) => (params = {}) => {
			return state.ConsensusStates[JSON.stringify(params)] ?? {}
		},
        getClientParams: (state) => (params = {}) => {
			return state.ClientParams[JSON.stringify(params)] ?? {}
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
		async QueryClientState({ commit, rootGetters,state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryClientState.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryClientState.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ClientState', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryClientState', payload: { all, ...key} })
				return state.ClientState[JSON.stringify(key)] ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryClientState', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryClientStates({ commit, rootGetters,state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryClientStates.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryClientStates.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ClientStates', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryClientStates', payload: { all, ...key} })
				return state.ClientStates[JSON.stringify(key)] ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryClientStates', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryConsensusState({ commit, rootGetters,state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryConsensusState.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryConsensusState.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ConsensusState', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryConsensusState', payload: { all, ...key} })
				return state.ConsensusState[JSON.stringify(key)] ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryConsensusState', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryConsensusStates({ commit, rootGetters,state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryConsensusStates.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryConsensusStates.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ConsensusStates', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryConsensusStates', payload: { all, ...key} })
				return state.ConsensusStates[JSON.stringify(key)] ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryConsensusStates', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryClientParams({ commit, rootGetters,state }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryClientParams.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryClientParams.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ClientParams', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryClientParams', payload: { all, ...key} })
				return state.ClientParams[JSON.stringify(key)] ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryClientParams', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		
		async sendMsgUpgradeClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUpgradeClient(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgUpgradeClient:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgUpgradeClient:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgSubmitMisbehaviour({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSubmitMisbehaviour(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSubmitMisbehaviour:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSubmitMisbehaviour:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgCreateClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgCreateClient(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgCreateClient:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCreateClient:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgUpdateClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUpdateClient(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgUpdateClient:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgUpdateClient:Send', 'Could not broadcast Tx.')
				}
			}
		},
		
		async MsgUpgradeClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUpgradeClient(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgUpgradeClient:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgUpgradeClient:Create', 'Could not create message.')
				}
			}
		},
		async MsgSubmitMisbehaviour({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSubmitMisbehaviour(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSubmitMisbehaviour:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSubmitMisbehaviour:Create', 'Could not create message.')
				}
			}
		},
		async MsgCreateClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgCreateClient(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgCreateClient:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCreateClient:Create', 'Could not create message.')
				}
			}
		},
		async MsgUpdateClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUpdateClient(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgUpdateClient:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgUpdateClient:Create', 'Could not create message.')
				}
			}
		},
		
	}
}
