import { txClient, queryClient } from './module'

import { Params } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { ValidatorHistoricalRewards } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { ValidatorCurrentRewards } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { ValidatorAccumulatedCommission } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { ValidatorOutstandingRewards } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { ValidatorSlashEvent } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { ValidatorSlashEvents } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { FeePool } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { CommunityPoolSpendProposal } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { DelegatorStartingInfo } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { DelegationDelegatorReward } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { CommunityPoolSpendProposalWithDeposit } from "./module/types/cosmos/distribution/v1beta1/distribution"
import { DelegatorWithdrawInfo } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorOutstandingRewardsRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorAccumulatedCommissionRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorHistoricalRewardsRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorCurrentRewardsRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { DelegatorStartingInfoRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorSlashEventRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"


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
        getParams: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getValidatorOutstandingRewards: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getValidatorCommission: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getValidatorSlashes: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDelegationRewards: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDelegationTotalRewards: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDelegatorValidators: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDelegatorWithdrawAddress: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getCommunityPool: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        
        _Structure: {
            Params: getStructure(Params.fromPartial({})),
            ValidatorHistoricalRewards: getStructure(ValidatorHistoricalRewards.fromPartial({})),
            ValidatorCurrentRewards: getStructure(ValidatorCurrentRewards.fromPartial({})),
            ValidatorAccumulatedCommission: getStructure(ValidatorAccumulatedCommission.fromPartial({})),
            ValidatorOutstandingRewards: getStructure(ValidatorOutstandingRewards.fromPartial({})),
            ValidatorSlashEvent: getStructure(ValidatorSlashEvent.fromPartial({})),
            ValidatorSlashEvents: getStructure(ValidatorSlashEvents.fromPartial({})),
            FeePool: getStructure(FeePool.fromPartial({})),
            CommunityPoolSpendProposal: getStructure(CommunityPoolSpendProposal.fromPartial({})),
            DelegatorStartingInfo: getStructure(DelegatorStartingInfo.fromPartial({})),
            DelegationDelegatorReward: getStructure(DelegationDelegatorReward.fromPartial({})),
            CommunityPoolSpendProposalWithDeposit: getStructure(CommunityPoolSpendProposalWithDeposit.fromPartial({})),
            DelegatorWithdrawInfo: getStructure(DelegatorWithdrawInfo.fromPartial({})),
            ValidatorOutstandingRewardsRecord: getStructure(ValidatorOutstandingRewardsRecord.fromPartial({})),
            ValidatorAccumulatedCommissionRecord: getStructure(ValidatorAccumulatedCommissionRecord.fromPartial({})),
            ValidatorHistoricalRewardsRecord: getStructure(ValidatorHistoricalRewardsRecord.fromPartial({})),
            ValidatorCurrentRewardsRecord: getStructure(ValidatorCurrentRewardsRecord.fromPartial({})),
            DelegatorStartingInfoRecord: getStructure(DelegatorStartingInfoRecord.fromPartial({})),
            ValidatorSlashEventRecord: getStructure(ValidatorSlashEventRecord.fromPartial({})),
            
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
        getParams: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getValidatorOutstandingRewards: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getValidatorCommission: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getValidatorSlashes: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDelegationRewards: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDelegationTotalRewards: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDelegatorValidators: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getDelegatorWithdrawAddress: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getCommunityPool: (state) => (params = {}) => {
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
        async QueryParams({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryParams.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryValidatorOutstandingRewards({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidatorOutstandingRewards.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryValidatorCommission({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidatorCommission.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryValidatorSlashes({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidatorSlashes.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryDelegationRewards({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegationRewards.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryDelegationTotalRewards({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegationTotalRewards.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryDelegatorValidators({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegatorValidators.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryDelegatorWithdrawAddress({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegatorWithdrawAddress.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryCommunityPool({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryCommunityPool.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        
        async MsgFundCommunityPool({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgFundCommunityPool(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgSetWithdrawAddress({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSetWithdrawAddress(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgWithdrawDelegatorReward({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgWithdrawDelegatorReward(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgWithdrawValidatorCommission({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgWithdrawValidatorCommission(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        
	}
}
