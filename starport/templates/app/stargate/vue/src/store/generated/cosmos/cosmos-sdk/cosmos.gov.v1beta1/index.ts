import { txClient, queryClient } from './module'

import { TextProposal } from "./module/types/cosmos/gov/v1beta1/gov"
import { Deposit } from "./module/types/cosmos/gov/v1beta1/gov"
import { Proposal } from "./module/types/cosmos/gov/v1beta1/gov"
import { TallyResult } from "./module/types/cosmos/gov/v1beta1/gov"
import { Vote } from "./module/types/cosmos/gov/v1beta1/gov"
import { DepositParams } from "./module/types/cosmos/gov/v1beta1/gov"
import { VotingParams } from "./module/types/cosmos/gov/v1beta1/gov"
import { TallyParams } from "./module/types/cosmos/gov/v1beta1/gov"


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
        getProposal: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getProposals: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getVote: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getVotes: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getParams: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDeposit: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDeposits: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getTallyResult: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        
        _Structure: {
            TextProposal: getStructure(TextProposal.fromPartial({})),
            Deposit: getStructure(Deposit.fromPartial({})),
            Proposal: getStructure(Proposal.fromPartial({})),
            TallyResult: getStructure(TallyResult.fromPartial({})),
            Vote: getStructure(Vote.fromPartial({})),
            DepositParams: getStructure(DepositParams.fromPartial({})),
            VotingParams: getStructure(VotingParams.fromPartial({})),
            TallyParams: getStructure(TallyParams.fromPartial({})),
            
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
        getProposal: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getProposals: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getVote: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getVotes: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getParams: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDeposit: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDeposits: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getTallyResult: (state) => (params = {}) => {
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
        async QueryProposal({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryProposal.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryProposals({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryProposals.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryVote({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryVote.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryVotes({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryVotes.apply(null, Object.values(key))).data
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
        async QueryDeposit({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDeposit.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryDeposits({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDeposits.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryTallyResult({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryTallyResult.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        
        async MsgVote({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgVote(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgSubmitProposal({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSubmitProposal(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgDeposit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgDeposit(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        
	}
}
