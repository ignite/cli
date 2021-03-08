import { txClient, queryClient } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { LastValidatorPower } from "./module/types/cosmos/staking/v1beta1/genesis"
import { HistoricalInfo } from "./module/types/cosmos/staking/v1beta1/staking"
import { CommissionRates } from "./module/types/cosmos/staking/v1beta1/staking"
import { Commission } from "./module/types/cosmos/staking/v1beta1/staking"
import { Description } from "./module/types/cosmos/staking/v1beta1/staking"
import { Validator } from "./module/types/cosmos/staking/v1beta1/staking"
import { ValAddresses } from "./module/types/cosmos/staking/v1beta1/staking"
import { DVPair } from "./module/types/cosmos/staking/v1beta1/staking"
import { DVPairs } from "./module/types/cosmos/staking/v1beta1/staking"
import { DVVTriplet } from "./module/types/cosmos/staking/v1beta1/staking"
import { DVVTriplets } from "./module/types/cosmos/staking/v1beta1/staking"
import { Delegation } from "./module/types/cosmos/staking/v1beta1/staking"
import { UnbondingDelegation } from "./module/types/cosmos/staking/v1beta1/staking"
import { UnbondingDelegationEntry } from "./module/types/cosmos/staking/v1beta1/staking"
import { RedelegationEntry } from "./module/types/cosmos/staking/v1beta1/staking"
import { Redelegation } from "./module/types/cosmos/staking/v1beta1/staking"
import { Params } from "./module/types/cosmos/staking/v1beta1/staking"
import { DelegationResponse } from "./module/types/cosmos/staking/v1beta1/staking"
import { RedelegationEntryResponse } from "./module/types/cosmos/staking/v1beta1/staking"
import { RedelegationResponse } from "./module/types/cosmos/staking/v1beta1/staking"
import { Pool } from "./module/types/cosmos/staking/v1beta1/staking"


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
        Validators: {},
        Validator: {},
        ValidatorDelegations: {},
        ValidatorUnbondingDelegations: {},
        Delegation: {},
        UnbondingDelegation: {},
        DelegatorDelegations: {},
        DelegatorUnbondingDelegations: {},
        Redelegations: {},
        DelegatorValidators: {},
        DelegatorValidator: {},
        HistoricalInfo: {},
        Pool: {},
        Params: {},
        
        _Structure: {
            LastValidatorPower: getStructure(LastValidatorPower.fromPartial({})),
            HistoricalInfo: getStructure(HistoricalInfo.fromPartial({})),
            CommissionRates: getStructure(CommissionRates.fromPartial({})),
            Commission: getStructure(Commission.fromPartial({})),
            Description: getStructure(Description.fromPartial({})),
            Validator: getStructure(Validator.fromPartial({})),
            ValAddresses: getStructure(ValAddresses.fromPartial({})),
            DVPair: getStructure(DVPair.fromPartial({})),
            DVPairs: getStructure(DVPairs.fromPartial({})),
            DVVTriplet: getStructure(DVVTriplet.fromPartial({})),
            DVVTriplets: getStructure(DVVTriplets.fromPartial({})),
            Delegation: getStructure(Delegation.fromPartial({})),
            UnbondingDelegation: getStructure(UnbondingDelegation.fromPartial({})),
            UnbondingDelegationEntry: getStructure(UnbondingDelegationEntry.fromPartial({})),
            RedelegationEntry: getStructure(RedelegationEntry.fromPartial({})),
            Redelegation: getStructure(Redelegation.fromPartial({})),
            Params: getStructure(Params.fromPartial({})),
            DelegationResponse: getStructure(DelegationResponse.fromPartial({})),
            RedelegationEntryResponse: getStructure(RedelegationEntryResponse.fromPartial({})),
            RedelegationResponse: getStructure(RedelegationResponse.fromPartial({})),
            Pool: getStructure(Pool.fromPartial({})),
            
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
        getValidators: (state) => (params = {}) => {
			return state.Validators[JSON.stringify(params)] ?? {}
		},
        getValidator: (state) => (params = {}) => {
			return state.Validator[JSON.stringify(params)] ?? {}
		},
        getValidatorDelegations: (state) => (params = {}) => {
			return state.ValidatorDelegations[JSON.stringify(params)] ?? {}
		},
        getValidatorUnbondingDelegations: (state) => (params = {}) => {
			return state.ValidatorUnbondingDelegations[JSON.stringify(params)] ?? {}
		},
        getDelegation: (state) => (params = {}) => {
			return state.Delegation[JSON.stringify(params)] ?? {}
		},
        getUnbondingDelegation: (state) => (params = {}) => {
			return state.UnbondingDelegation[JSON.stringify(params)] ?? {}
		},
        getDelegatorDelegations: (state) => (params = {}) => {
			return state.DelegatorDelegations[JSON.stringify(params)] ?? {}
		},
        getDelegatorUnbondingDelegations: (state) => (params = {}) => {
			return state.DelegatorUnbondingDelegations[JSON.stringify(params)] ?? {}
		},
        getRedelegations: (state) => (params = {}) => {
			return state.Redelegations[JSON.stringify(params)] ?? {}
		},
        getDelegatorValidators: (state) => (params = {}) => {
			return state.DelegatorValidators[JSON.stringify(params)] ?? {}
		},
        getDelegatorValidator: (state) => (params = {}) => {
			return state.DelegatorValidator[JSON.stringify(params)] ?? {}
		},
        getHistoricalInfo: (state) => (params = {}) => {
			return state.HistoricalInfo[JSON.stringify(params)] ?? {}
		},
        getPool: (state) => (params = {}) => {
			return state.Pool[JSON.stringify(params)] ?? {}
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
		async QueryValidators({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidators.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Validators', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidators', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidators', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryValidator({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidator.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Validator', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidator', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidator', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryValidatorDelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidatorDelegations.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'ValidatorDelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidatorDelegations', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidatorDelegations', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryValidatorUnbondingDelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryValidatorUnbondingDelegations.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'ValidatorUnbondingDelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidatorUnbondingDelegations', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidatorUnbondingDelegations', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDelegation({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegation.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Delegation', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegation', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegation', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryUnbondingDelegation({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryUnbondingDelegation.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'UnbondingDelegation', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryUnbondingDelegation', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryUnbondingDelegation', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDelegatorDelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegatorDelegations.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DelegatorDelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorDelegations', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorDelegations', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryDelegatorUnbondingDelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegatorUnbondingDelegations.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DelegatorUnbondingDelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorUnbondingDelegations', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorUnbondingDelegations', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryRedelegations({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryRedelegations.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Redelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryRedelegations', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryRedelegations', 'API Node Unavailable. Could not perform query.'))
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
		async QueryDelegatorValidator({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryDelegatorValidator.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'DelegatorValidator', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorValidator', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorValidator', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryHistoricalInfo({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryHistoricalInfo.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'HistoricalInfo', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryHistoricalInfo', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryHistoricalInfo', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryPool({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPool.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Pool', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPool', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryPool', 'API Node Unavailable. Could not perform query.'))
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
		
		async sendMsgEditValidator({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgEditValidator(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgEditValidator:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgEditValidator:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgUndelegate({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUndelegate(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgUndelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgUndelegate:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgBeginRedelegate({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgBeginRedelegate(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgBeginRedelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgBeginRedelegate:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgDelegate({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgDelegate(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgDelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgDelegate:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgCreateValidator({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgCreateValidator(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgCreateValidator:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCreateValidator:Send', 'Could not broadcast Tx.')
				}
			}
		},
		
		async MsgEditValidator({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgEditValidator(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgEditValidator:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgEditValidator:Create', 'Could not create message.')
				}
			}
		},
		async MsgUndelegate({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUndelegate(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgUndelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgUndelegate:Create', 'Could not create message.')
				}
			}
		},
		async MsgBeginRedelegate({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgBeginRedelegate(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgBeginRedelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgBeginRedelegate:Create', 'Could not create message.')
				}
			}
		},
		async MsgDelegate({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgDelegate(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgDelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgDelegate:Create', 'Could not create message.')
				}
			}
		},
		async MsgCreateValidator({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgCreateValidator(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgCreateValidator:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCreateValidator:Create', 'Could not create message.')
				}
			}
		},
		
	}
}
