import { txClient, queryClient } from './module'

import { Channel } from "./module/types/ibc/core/channel/v1/channel"
import { IdentifiedChannel } from "./module/types/ibc/core/channel/v1/channel"
import { Counterparty } from "./module/types/ibc/core/channel/v1/channel"
import { Packet } from "./module/types/ibc/core/channel/v1/channel"
import { PacketState } from "./module/types/ibc/core/channel/v1/channel"
import { Acknowledgement } from "./module/types/ibc/core/channel/v1/channel"
import { PacketSequence } from "./module/types/ibc/core/channel/v1/genesis"


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
        getChannel: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getChannels: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConnectionChannels: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getChannelClientState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getChannelConsensusState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketCommitment: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketCommitments: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketReceipt: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketAcknowledgement: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketAcknowledgements: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getUnreceivedPackets: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getUnreceivedAcks: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getNextSequenceReceive: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        
        _Structure: {
            Channel: getStructure(Channel.fromPartial({})),
            IdentifiedChannel: getStructure(IdentifiedChannel.fromPartial({})),
            Counterparty: getStructure(Counterparty.fromPartial({})),
            Packet: getStructure(Packet.fromPartial({})),
            PacketState: getStructure(PacketState.fromPartial({})),
            Acknowledgement: getStructure(Acknowledgement.fromPartial({})),
            PacketSequence: getStructure(PacketSequence.fromPartial({})),
            
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
        getChannel: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getChannels: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getConnectionChannels: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getChannelClientState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getChannelConsensusState: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketCommitment: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketCommitments: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketReceipt: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketAcknowledgement: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getPacketAcknowledgements: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getUnreceivedPackets: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getUnreceivedAcks: (state) => (params = {}) => {
			return state.Post[JSON.stringify(params)] ?? {}
		},
        getNextSequenceReceive: (state) => (params = {}) => {
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
        async QueryChannel({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryChannel.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryChannels({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryChannels.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryConnectionChannels({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryConnectionChannels.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryChannelClientState({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryChannelClientState.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryChannelConsensusState({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryChannelConsensusState.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryPacketCommitment({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketCommitment.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryPacketCommitments({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketCommitments.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryPacketReceipt({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketReceipt.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryPacketAcknowledgement({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketAcknowledgement.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryPacketAcknowledgements({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketAcknowledgements.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryUnreceivedPackets({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryUnreceivedPackets.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryUnreceivedAcks({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryUnreceivedAcks.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        async QueryNextSequenceReceive({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryNextSequenceReceive.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Post', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPost', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        
        async MsgTimeout({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgTimeout(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgChannelOpenInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenInit(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgTimeoutOnClose({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgTimeoutOnClose(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgAcknowledgement({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgAcknowledgement(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgChannelOpenTry({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenTry(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgRecvPacket({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgRecvPacket(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgChannelOpenAck({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenAck(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgChannelCloseConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelCloseConfirm(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgChannelOpenConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenConfirm(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        async MsgChannelCloseInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelCloseInit(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        
	}
}
