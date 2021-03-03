import { txClient, queryClient } from './module'

import { IdentifiedClientState } from "./module/types/ibc/core/client/v1/client"
import { ConsensusStateWithHeight } from "./module/types/ibc/core/client/v1/client"
import { ClientConsensusStates } from "./module/types/ibc/core/client/v1/client"
import { ClientUpdateProposal } from "./module/types/ibc/core/client/v1/client"
import { Height } from "./module/types/ibc/core/client/v1/client"
import { Params } from "./module/types/ibc/core/client/v1/client"
import { GenesisMetadata } from "./module/types/ibc/core/client/v1/genesis"
import { IdentifiedGenesisMetadata } from "./module/types/ibc/core/client/v1/genesis"


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
        getClientState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getClientStates: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConsensusState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConsensusStates: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getClientParams: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        
        _Structure: {
            IdentifiedClientState: getStructure(IdentifiedClientState.fromPartial({})),
            ConsensusStateWithHeight: getStructure(ConsensusStateWithHeight.fromPartial({})),
            ClientConsensusStates: getStructure(ClientConsensusStates.fromPartial({})),
            ClientUpdateProposal: getStructure(ClientUpdateProposal.fromPartial({})),
            Height: getStructure(Height.fromPartial({})),
            Params: getStructure(Params.fromPartial({})),
            GenesisMetadata: getStructure(GenesisMetadata.fromPartial({})),
            IdentifiedGenesisMetadata: getStructure(IdentifiedGenesisMetadata.fromPartial({})),
            
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
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getClientStates: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConsensusState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConsensusStates: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getClientParams: (state) => (params = {}) => {
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
        async QueryClientState({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryClientState.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryClientStates({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryClientStates.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryConsensusState({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryConsensusState.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryConsensusStates({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryConsensusStates.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryClientParams({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryClientParams.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        
        async MsgSubmitMisbehaviour({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSubmitMisbehaviour(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgCreateClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgCreateClient(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgUpdateClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUpdateClient(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgUpgradeClient({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUpgradeClient(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        
	}
}
