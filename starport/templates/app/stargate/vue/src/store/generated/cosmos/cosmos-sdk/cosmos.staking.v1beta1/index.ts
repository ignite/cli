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
		async QueryValidators({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryValidators.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryValidators.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'Validators', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidators', payload: { all, ...key} })
				return getters['getValidators'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidators', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryValidator({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryValidator.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryValidator.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'Validator', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidator', payload: { all, ...key} })
				return getters['getValidator'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidator', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryValidatorDelegations({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryValidatorDelegations.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryValidatorDelegations.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ValidatorDelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidatorDelegations', payload: { all, ...key} })
				return getters['getValidatorDelegations'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidatorDelegations', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryValidatorUnbondingDelegations({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryValidatorUnbondingDelegations.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryValidatorUnbondingDelegations.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'ValidatorUnbondingDelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidatorUnbondingDelegations', payload: { all, ...key} })
				return getters['getValidatorUnbondingDelegations'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryValidatorUnbondingDelegations', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryDelegation({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryDelegation.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryDelegation.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'Delegation', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegation', payload: { all, ...key} })
				return getters['getDelegation'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegation', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryUnbondingDelegation({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryUnbondingDelegation.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryUnbondingDelegation.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'UnbondingDelegation', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryUnbondingDelegation', payload: { all, ...key} })
				return getters['getUnbondingDelegation'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryUnbondingDelegation', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryDelegatorDelegations({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryDelegatorDelegations.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryDelegatorDelegations.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'DelegatorDelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorDelegations', payload: { all, ...key} })
				return getters['getDelegatorDelegations'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorDelegations', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryDelegatorUnbondingDelegations({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryDelegatorUnbondingDelegations.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryDelegatorUnbondingDelegations.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'DelegatorUnbondingDelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorUnbondingDelegations', payload: { all, ...key} })
				return getters['getDelegatorUnbondingDelegations'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorUnbondingDelegations', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryRedelegations({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryRedelegations.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryRedelegations.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'Redelegations', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryRedelegations', payload: { all, ...key} })
				return getters['getRedelegations'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryRedelegations', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryDelegatorValidators({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryDelegatorValidators.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryDelegatorValidators.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'DelegatorValidators', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorValidators', payload: { all, ...key} })
				return getters['getDelegatorValidators'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorValidators', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryDelegatorValidator({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryDelegatorValidator.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryDelegatorValidator.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'DelegatorValidator', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDelegatorValidator', payload: { all, ...key} })
				return getters['getDelegatorValidator'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDelegatorValidator', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryHistoricalInfo({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryHistoricalInfo.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryHistoricalInfo.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'HistoricalInfo', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryHistoricalInfo', payload: { all, ...key} })
				return getters['getHistoricalInfo'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryHistoricalInfo', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryPool({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryPool.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryPool.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'Pool', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPool', payload: { all, ...key} })
				return getters['getPool'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryPool', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryParams({ commit, rootGetters, getters }, { subscribe = false, all=false, ...key }) {
			try {
				let params=Object.values(key)
				let value = (await (await initQueryClient(rootGetters)).queryParams.apply(null, params)).data
				while (all && value.pagination && value.pagination.next_key!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryParams.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				commit('QUERY', { query: 'Params', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryParams', payload: { all, ...key} })
				return getters['getParams'](key) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryParams', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		
		async sendMsgEditValidator({ rootGetters }, { value, fee, memo }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgEditValidator(value)
				const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgEditValidator:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgEditValidator:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgBeginRedelegate({ rootGetters }, { value, fee, memo }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgBeginRedelegate(value)
				const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgBeginRedelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgBeginRedelegate:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgUndelegate({ rootGetters }, { value, fee, memo }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgUndelegate(value)
				const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgUndelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgUndelegate:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgCreateValidator({ rootGetters }, { value, fee, memo }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgCreateValidator(value)
				const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgCreateValidator:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCreateValidator:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgDelegate({ rootGetters }, { value, fee, memo }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgDelegate(value)
				const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgDelegate:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgDelegate:Send', 'Could not broadcast Tx.')
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
		
	}
}
