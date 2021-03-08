import { txClient, queryClient } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { PacketSequence } from "./module/types/ibc/core/channel/v1/genesis"
import { Channel } from "./module/types/ibc/core/channel/v1/channel"
import { IdentifiedChannel } from "./module/types/ibc/core/channel/v1/channel"
import { Counterparty } from "./module/types/ibc/core/channel/v1/channel"
import { Packet } from "./module/types/ibc/core/channel/v1/channel"
import { PacketState } from "./module/types/ibc/core/channel/v1/channel"
import { Acknowledgement } from "./module/types/ibc/core/channel/v1/channel"


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
        Channel: {},
        Channels: {},
        ConnectionChannels: {},
        ChannelClientState: {},
        ChannelConsensusState: {},
        PacketCommitment: {},
        PacketCommitments: {},
        PacketReceipt: {},
        PacketAcknowledgement: {},
        PacketAcknowledgements: {},
        UnreceivedPackets: {},
        UnreceivedAcks: {},
        NextSequenceReceive: {},
        
        _Structure: {
            PacketSequence: getStructure(PacketSequence.fromPartial({})),
            Channel: getStructure(Channel.fromPartial({})),
            IdentifiedChannel: getStructure(IdentifiedChannel.fromPartial({})),
            Counterparty: getStructure(Counterparty.fromPartial({})),
            Packet: getStructure(Packet.fromPartial({})),
            PacketState: getStructure(PacketState.fromPartial({})),
            Acknowledgement: getStructure(Acknowledgement.fromPartial({})),
            
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
			return state.Channel[JSON.stringify(params)] ?? {}
		},
        getChannels: (state) => (params = {}) => {
			return state.Channels[JSON.stringify(params)] ?? {}
		},
        getConnectionChannels: (state) => (params = {}) => {
			return state.ConnectionChannels[JSON.stringify(params)] ?? {}
		},
        getChannelClientState: (state) => (params = {}) => {
			return state.ChannelClientState[JSON.stringify(params)] ?? {}
		},
        getChannelConsensusState: (state) => (params = {}) => {
			return state.ChannelConsensusState[JSON.stringify(params)] ?? {}
		},
        getPacketCommitment: (state) => (params = {}) => {
			return state.PacketCommitment[JSON.stringify(params)] ?? {}
		},
        getPacketCommitments: (state) => (params = {}) => {
			return state.PacketCommitments[JSON.stringify(params)] ?? {}
		},
        getPacketReceipt: (state) => (params = {}) => {
			return state.PacketReceipt[JSON.stringify(params)] ?? {}
		},
        getPacketAcknowledgement: (state) => (params = {}) => {
			return state.PacketAcknowledgement[JSON.stringify(params)] ?? {}
		},
        getPacketAcknowledgements: (state) => (params = {}) => {
			return state.PacketAcknowledgements[JSON.stringify(params)] ?? {}
		},
        getUnreceivedPackets: (state) => (params = {}) => {
			return state.UnreceivedPackets[JSON.stringify(params)] ?? {}
		},
        getUnreceivedAcks: (state) => (params = {}) => {
			return state.UnreceivedAcks[JSON.stringify(params)] ?? {}
		},
        getNextSequenceReceive: (state) => (params = {}) => {
			return state.NextSequenceReceive[JSON.stringify(params)] ?? {}
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
		async QueryChannel({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryChannel.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Channel', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryChannel', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryChannel', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryChannels({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryChannels.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'Channels', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryChannels', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryChannels', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryConnectionChannels({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryConnectionChannels.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'ConnectionChannels', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryConnectionChannels', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryConnectionChannels', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryChannelClientState({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryChannelClientState.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'ChannelClientState', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryChannelClientState', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryChannelClientState', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryChannelConsensusState({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryChannelConsensusState.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'ChannelConsensusState', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryChannelConsensusState', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryChannelConsensusState', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryPacketCommitment({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketCommitment.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'PacketCommitment', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPacketCommitment', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryPacketCommitment', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryPacketCommitments({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketCommitments.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'PacketCommitments', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPacketCommitments', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryPacketCommitments', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryPacketReceipt({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketReceipt.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'PacketReceipt', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPacketReceipt', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryPacketReceipt', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryPacketAcknowledgement({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketAcknowledgement.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'PacketAcknowledgement', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPacketAcknowledgement', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryPacketAcknowledgement', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryPacketAcknowledgements({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryPacketAcknowledgements.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'PacketAcknowledgements', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPacketAcknowledgements', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryPacketAcknowledgements', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryUnreceivedPackets({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryUnreceivedPackets.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'UnreceivedPackets', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryUnreceivedPackets', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryUnreceivedPackets', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryUnreceivedAcks({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryUnreceivedAcks.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'UnreceivedAcks', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryUnreceivedAcks', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryUnreceivedAcks', 'API Node Unavailable. Could not perform query.'))
			}
		},
		async QueryNextSequenceReceive({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).queryNextSequenceReceive.apply(null, Object.values(key))).data
				commit('QUERY', { query: 'NextSequenceReceive', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryNextSequenceReceive', payload: key })
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryNextSequenceReceive', 'API Node Unavailable. Could not perform query.'))
			}
		},
		
		async sendMsgChannelOpenInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenInit(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelOpenInit:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelOpenInit:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgChannelCloseConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelCloseConfirm(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelCloseConfirm:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelCloseConfirm:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgTimeout({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgTimeout(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgTimeout:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgTimeout:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgRecvPacket({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgRecvPacket(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgRecvPacket:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgRecvPacket:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgChannelCloseInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelCloseInit(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelCloseInit:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelCloseInit:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgChannelOpenConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenConfirm(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelOpenConfirm:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelOpenConfirm:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgChannelOpenAck({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenAck(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelOpenAck:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelOpenAck:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgChannelOpenTry({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenTry(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelOpenTry:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelOpenTry:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgAcknowledgement({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgAcknowledgement(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgAcknowledgement:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgAcknowledgement:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgTimeoutOnClose({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgTimeoutOnClose(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgTimeoutOnClose:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgTimeoutOnClose:Send', 'Could not broadcast Tx.')
				}
			}
		},
		
		async MsgChannelOpenInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenInit(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelOpenInit:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelOpenInit:Create', 'Could not create message.')
				}
			}
		},
		async MsgChannelCloseConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelCloseConfirm(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelCloseConfirm:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelCloseConfirm:Create', 'Could not create message.')
				}
			}
		},
		async MsgTimeout({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgTimeout(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgTimeout:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgTimeout:Create', 'Could not create message.')
				}
			}
		},
		async MsgRecvPacket({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgRecvPacket(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgRecvPacket:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgRecvPacket:Create', 'Could not create message.')
				}
			}
		},
		async MsgChannelCloseInit({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelCloseInit(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelCloseInit:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelCloseInit:Create', 'Could not create message.')
				}
			}
		},
		async MsgChannelOpenConfirm({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenConfirm(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelOpenConfirm:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelOpenConfirm:Create', 'Could not create message.')
				}
			}
		},
		async MsgChannelOpenAck({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenAck(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelOpenAck:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelOpenAck:Create', 'Could not create message.')
				}
			}
		},
		async MsgChannelOpenTry({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgChannelOpenTry(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgChannelOpenTry:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgChannelOpenTry:Create', 'Could not create message.')
				}
			}
		},
		async MsgAcknowledgement({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgAcknowledgement(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgAcknowledgement:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgAcknowledgement:Create', 'Could not create message.')
				}
			}
		},
		async MsgTimeoutOnClose({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgTimeoutOnClose(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgTimeoutOnClose:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgTimeoutOnClose:Create', 'Could not create message.')
				}
			}
		},
		
	}
}
