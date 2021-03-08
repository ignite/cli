import { txClient, queryClient } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { DelegatorWithdrawInfo } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorOutstandingRewardsRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorAccumulatedCommissionRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorHistoricalRewardsRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorCurrentRewardsRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { DelegatorStartingInfoRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
import { ValidatorSlashEventRecord } from "./module/types/cosmos/distribution/v1beta1/genesis"
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
        Params: {},
        ValidatorOutstandingRewards: {},
        ValidatorCommission: {},
        ValidatorSlashes: {},
        DelegationRewards: {},
        DelegationTotalRewards: {},
        DelegatorValidators: {},
        DelegatorWithdrawAddress: {},
        CommunityPool: {},
        
        _Structure: {
            DelegatorWithdrawInfo: getStructure(DelegatorWithdrawInfo.fromPartial({})),
            ValidatorOutstandingRewardsRecord: getStructure(ValidatorOutstandingRewardsRecord.fromPartial({})),
            ValidatorAccumulatedCommissionRecord: getStructure(ValidatorAccumulatedCommissionRecord.fromPartial({})),
            ValidatorHistoricalRewardsRecord: getStructure(ValidatorHistoricalRewardsRecord.fromPartial({})),
            ValidatorCurrentRewardsRecord: getStructure(ValidatorCurrentRewardsRecord.fromPartial({})),
            DelegatorStartingInfoRecord: getStructure(DelegatorStartingInfoRecord.fromPartial({})),
            ValidatorSlashEventRecord: getStructure(ValidatorSlashEventRecord.fromPartial({})),
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
			return state.Params[JSON.stringify(params)] ?? {}
		},
        getValidatorOutstandingRewards: (state) => (params = {}) => {
			return state.ValidatorOutstandingRewards[JSON.stringify(params)] ?? {}
		},
        getValidatorCommission: (state) => (params = {}) => {
			return state.ValidatorCommission[JSON.stringify(params)] ?? {}
		},
        getValidatorSlashes: (state) => (params = {}) => {
			return state.ValidatorSlashes[JSON.stringify(params)] ?? {}
		},
        getDelegationRewards: (state) => (params = {}) => {
			return state.DelegationRewards[JSON.stringify(params)] ?? {}
		},
        getDelegationTotalRewards: (state) => (params = {}) => {
			return state.DelegationTotalRewards[JSON.stringify(params)] ?? {}
		},
        getDelegatorValidators: (state) => (params = {}) => {
			return state.DelegatorValidators[JSON.stringify(params)] ?? {}
		},
        getDelegatorWithdrawAddress: (state) => (params = {}) => {
			return state.DelegatorWithdrawAddress[JSON.stringify(params)] ?? {}
		},
        getCommunityPool: (state) => (params = {}) => {
			return state.CommunityPool[JSON.stringify(params)] ?? {}
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
		async QueryParams({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryParams.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Params', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryParams', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryParams', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryValidatorOutstandingRewards({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidatorOutstandingRewards.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'ValidatorOutstandingRewards', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidatorOutstandingRewards', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidatorOutstandingRewards', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryValidatorCommission({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidatorCommission.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'ValidatorCommission', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidatorCommission', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidatorCommission', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryValidatorSlashes({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidatorSlashes.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'ValidatorSlashes', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidatorSlashes', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidatorSlashes', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDelegationRewards({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegationRewards.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DelegationRewards', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegationRewards', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegationRewards', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDelegationTotalRewards({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegationTotalRewards.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DelegationTotalRewards', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegationTotalRewards', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegationTotalRewards', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDelegatorValidators({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegatorValidators.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DelegatorValidators', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorValidators', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorValidators', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDelegatorWithdrawAddress({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegatorWithdrawAddress.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DelegatorWithdrawAddress', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorWithdrawAddress', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorWithdrawAddress', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryCommunityPool({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryCommunityPool.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'CommunityPool', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryCommunityPool', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryCommunityPool', 'API Node Unavailable. Could not perform query.'))
			}
		},
		
		async sendMsgSetWithdrawAddress({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSetWithdrawAddress(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSetWithdrawAddress:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSetWithdrawAddress:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgWithdrawDelegatorReward({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgWithdrawDelegatorReward(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgWithdrawDelegatorReward:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgWithdrawDelegatorReward:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgWithdrawValidatorCommission({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgWithdrawValidatorCommission(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgWithdrawValidatorCommission:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgWithdrawValidatorCommission:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgFundCommunityPool({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgFundCommunityPool(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgFundCommunityPool:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgFundCommunityPool:Send', 'Could not broadcast Tx.')
				}
			}
		},
		
		async MsgSetWithdrawAddress({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSetWithdrawAddress(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSetWithdrawAddress:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSetWithdrawAddress:Create', 'Could not create message.')
				}
			}
		},
		async MsgWithdrawDelegatorReward({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgWithdrawDelegatorReward(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgWithdrawDelegatorReward:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgWithdrawDelegatorReward:Create', 'Could not create message.')
				}
			}
		},
		async MsgWithdrawValidatorCommission({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgWithdrawValidatorCommission(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgWithdrawValidatorCommission:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgWithdrawValidatorCommission:Create', 'Could not create message.')
				}
			}
		},
		async MsgFundCommunityPool({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgFundCommunityPool(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgFundCommunityPool:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgFundCommunityPool:Create', 'Could not create message.')
				}
			}
		},
		
	}
}
