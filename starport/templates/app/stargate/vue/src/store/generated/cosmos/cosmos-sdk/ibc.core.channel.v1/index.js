import { txClient, queryClient } from './module';
// @ts-ignore
import { SpVuexError } from '@starport/vuex';
import { Channel } from "./module/types/ibc/core/channel/v1/channel";
import { IdentifiedChannel } from "./module/types/ibc/core/channel/v1/channel";
import { Counterparty } from "./module/types/ibc/core/channel/v1/channel";
import { Packet } from "./module/types/ibc/core/channel/v1/channel";
import { PacketState } from "./module/types/ibc/core/channel/v1/channel";
import { Acknowledgement } from "./module/types/ibc/core/channel/v1/channel";
import { PacketSequence } from "./module/types/ibc/core/channel/v1/genesis";
async function initTxClient(vuexGetters) {
    return await txClient(vuexGetters['common/wallet/signer'], {
        addr: vuexGetters['common/env/apiTendermint']
    });
}
async function initQueryClient(vuexGetters) {
    return await queryClient({
        addr: vuexGetters['common/env/apiCosmos']
    });
}
function getStructure(template) {
    let structure = { fields: [] };
    for (const [key, value] of Object.entries(template)) {
        let field = {};
        field.name = key;
        field.type = typeof value;
        structure.fields.push(field);
    }
    return structure;
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
            Channel: getStructure(Channel.fromPartial({})),
            IdentifiedChannel: getStructure(IdentifiedChannel.fromPartial({})),
            Counterparty: getStructure(Counterparty.fromPartial({})),
            Packet: getStructure(Packet.fromPartial({})),
            PacketState: getStructure(PacketState.fromPartial({})),
            Acknowledgement: getStructure(Acknowledgement.fromPartial({})),
            PacketSequence: getStructure(PacketSequence.fromPartial({})),
        },
        _Subscriptions: new Set(),
    };
};
// initial state
const state = getDefaultState();
export default {
    namespaced: true,
    state,
    mutations: {
        RESET_STATE(state) {
            Object.assign(state, getDefaultState());
        },
        QUERY(state, { query, key, value }) {
            state[query][JSON.stringify(key)] = value;
        },
        SUBSCRIBE(state, subscription) {
            state._Subscriptions.add(subscription);
        },
        UNSUBSCRIBE(state, subscription) {
            state._Subscriptions.delete(subscription);
        }
    },
    getters: {
        getChannel: (state) => (params = {}) => {
            return state.Channel[JSON.stringify(params)] ?? {};
        },
        getChannels: (state) => (params = {}) => {
            return state.Channels[JSON.stringify(params)] ?? {};
        },
        getConnectionChannels: (state) => (params = {}) => {
            return state.ConnectionChannels[JSON.stringify(params)] ?? {};
        },
        getChannelClientState: (state) => (params = {}) => {
            return state.ChannelClientState[JSON.stringify(params)] ?? {};
        },
        getChannelConsensusState: (state) => (params = {}) => {
            return state.ChannelConsensusState[JSON.stringify(params)] ?? {};
        },
        getPacketCommitment: (state) => (params = {}) => {
            return state.PacketCommitment[JSON.stringify(params)] ?? {};
        },
        getPacketCommitments: (state) => (params = {}) => {
            return state.PacketCommitments[JSON.stringify(params)] ?? {};
        },
        getPacketReceipt: (state) => (params = {}) => {
            return state.PacketReceipt[JSON.stringify(params)] ?? {};
        },
        getPacketAcknowledgement: (state) => (params = {}) => {
            return state.PacketAcknowledgement[JSON.stringify(params)] ?? {};
        },
        getPacketAcknowledgements: (state) => (params = {}) => {
            return state.PacketAcknowledgements[JSON.stringify(params)] ?? {};
        },
        getUnreceivedPackets: (state) => (params = {}) => {
            return state.UnreceivedPackets[JSON.stringify(params)] ?? {};
        },
        getUnreceivedAcks: (state) => (params = {}) => {
            return state.UnreceivedAcks[JSON.stringify(params)] ?? {};
        },
        getNextSequenceReceive: (state) => (params = {}) => {
            return state.NextSequenceReceive[JSON.stringify(params)] ?? {};
        },
        getTypeStructure: (state) => (type) => {
            return state._Structure[type].fields;
        }
    },
    actions: {
        init({ dispatch, rootGetters }) {
            console.log('init');
            if (rootGetters['common/env/client']) {
                rootGetters['common/env/client'].on('newblock', () => {
                    dispatch('StoreUpdate');
                });
            }
        },
        resetState({ commit }) {
            commit('RESET_STATE');
        },
        unsubscribe({ commit }, subscription) {
            commit('UNSUBSCRIBE', subscription);
        },
        async StoreUpdate({ state, dispatch }) {
            state._Subscriptions.forEach((subscription) => {
                dispatch(subscription.action, subscription.payload);
            });
        },
        async QueryChannel({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryChannel.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryChannel.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'Channel', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryChannel', payload: { all, ...key } });
                return getters['getChannel'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryChannel', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryChannels({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryChannels.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryChannels.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'Channels', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryChannels', payload: { all, ...key } });
                return getters['getChannels'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryChannels', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryConnectionChannels({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryConnectionChannels.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryConnectionChannels.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'ConnectionChannels', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryConnectionChannels', payload: { all, ...key } });
                return getters['getConnectionChannels'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryConnectionChannels', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryChannelClientState({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryChannelClientState.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryChannelClientState.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'ChannelClientState', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryChannelClientState', payload: { all, ...key } });
                return getters['getChannelClientState'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryChannelClientState', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryChannelConsensusState({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryChannelConsensusState.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryChannelConsensusState.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'ChannelConsensusState', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryChannelConsensusState', payload: { all, ...key } });
                return getters['getChannelConsensusState'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryChannelConsensusState', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryPacketCommitment({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryPacketCommitment.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryPacketCommitment.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'PacketCommitment', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPacketCommitment', payload: { all, ...key } });
                return getters['getPacketCommitment'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryPacketCommitment', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryPacketCommitments({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryPacketCommitments.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryPacketCommitments.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'PacketCommitments', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPacketCommitments', payload: { all, ...key } });
                return getters['getPacketCommitments'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryPacketCommitments', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryPacketReceipt({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryPacketReceipt.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryPacketReceipt.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'PacketReceipt', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPacketReceipt', payload: { all, ...key } });
                return getters['getPacketReceipt'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryPacketReceipt', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryPacketAcknowledgement({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryPacketAcknowledgement.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryPacketAcknowledgement.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'PacketAcknowledgement', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPacketAcknowledgement', payload: { all, ...key } });
                return getters['getPacketAcknowledgement'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryPacketAcknowledgement', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryPacketAcknowledgements({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryPacketAcknowledgements.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryPacketAcknowledgements.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'PacketAcknowledgements', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryPacketAcknowledgements', payload: { all, ...key } });
                return getters['getPacketAcknowledgements'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryPacketAcknowledgements', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryUnreceivedPackets({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryUnreceivedPackets.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryUnreceivedPackets.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'UnreceivedPackets', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryUnreceivedPackets', payload: { all, ...key } });
                return getters['getUnreceivedPackets'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryUnreceivedPackets', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryUnreceivedAcks({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryUnreceivedAcks.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryUnreceivedAcks.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'UnreceivedAcks', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryUnreceivedAcks', payload: { all, ...key } });
                return getters['getUnreceivedAcks'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryUnreceivedAcks', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async QueryNextSequenceReceive({ commit, rootGetters, getters }, { subscribe = false, all = false, ...key }) {
            try {
                let params = Object.values(key);
                let value = (await (await initQueryClient(rootGetters)).queryNextSequenceReceive.apply(null, params)).data;
                while (all && value.pagination && value.pagination.next_key != null) {
                    let next_values = (await (await initQueryClient(rootGetters)).queryNextSequenceReceive.apply(null, [...params, { 'pagination.key': value.pagination.next_key }])).data;
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop] = [...value[prop], ...next_values[prop]];
                        }
                        else {
                            value[prop] = next_values[prop];
                        }
                    }
                }
                commit('QUERY', { query: 'NextSequenceReceive', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryNextSequenceReceive', payload: { all, ...key } });
                return getters['getNextSequenceReceive'](key) ?? {};
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryNextSequenceReceive', 'API Node Unavailable. Could not perform query.'));
                return {};
            }
        },
        async sendMsgTimeoutOnClose({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgTimeoutOnClose(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgTimeoutOnClose:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgTimeoutOnClose:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgChannelOpenConfirm({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelOpenConfirm(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelOpenConfirm:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelOpenConfirm:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgAcknowledgement({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgAcknowledgement(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgAcknowledgement:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgAcknowledgement:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgChannelOpenTry({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelOpenTry(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelOpenTry:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelOpenTry:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgTimeout({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgTimeout(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgTimeout:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgTimeout:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgChannelOpenAck({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelOpenAck(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelOpenAck:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelOpenAck:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgChannelCloseConfirm({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelCloseConfirm(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelCloseConfirm:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelCloseConfirm:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgRecvPacket({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgRecvPacket(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgRecvPacket:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgRecvPacket:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgChannelCloseInit({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelCloseInit(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelCloseInit:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelCloseInit:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgChannelOpenInit({ rootGetters }, { value, fee, memo }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelOpenInit(value);
                const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], { fee: { amount: fee,
                        gas: "200000" }, memo });
                return result;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelOpenInit:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelOpenInit:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async MsgTimeoutOnClose({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgTimeoutOnClose(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgTimeoutOnClose:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgTimeoutOnClose:Create', 'Could not create message.');
                }
            }
        },
        async MsgChannelOpenConfirm({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelOpenConfirm(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelOpenConfirm:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelOpenConfirm:Create', 'Could not create message.');
                }
            }
        },
        async MsgAcknowledgement({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgAcknowledgement(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgAcknowledgement:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgAcknowledgement:Create', 'Could not create message.');
                }
            }
        },
        async MsgChannelOpenTry({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelOpenTry(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelOpenTry:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelOpenTry:Create', 'Could not create message.');
                }
            }
        },
        async MsgTimeout({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgTimeout(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgTimeout:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgTimeout:Create', 'Could not create message.');
                }
            }
        },
        async MsgChannelOpenAck({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelOpenAck(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelOpenAck:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelOpenAck:Create', 'Could not create message.');
                }
            }
        },
        async MsgChannelCloseConfirm({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelCloseConfirm(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelCloseConfirm:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelCloseConfirm:Create', 'Could not create message.');
                }
            }
        },
        async MsgRecvPacket({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgRecvPacket(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgRecvPacket:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgRecvPacket:Create', 'Could not create message.');
                }
            }
        },
        async MsgChannelCloseInit({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelCloseInit(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelCloseInit:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelCloseInit:Create', 'Could not create message.');
                }
            }
        },
        async MsgChannelOpenInit({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgChannelOpenInit(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgChannelOpenInit:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgChannelOpenInit:Create', 'Could not create message.');
                }
            }
        },
    }
};
