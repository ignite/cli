import { txClient, queryClient } from './module';
// @ts-ignore
import { SpVuexError } from '@starport/vuex';
import { TextProposal } from "./module/types/cosmos/gov/v1beta1/gov";
import { Deposit } from "./module/types/cosmos/gov/v1beta1/gov";
import { Proposal } from "./module/types/cosmos/gov/v1beta1/gov";
import { TallyResult } from "./module/types/cosmos/gov/v1beta1/gov";
import { Vote } from "./module/types/cosmos/gov/v1beta1/gov";
import { DepositParams } from "./module/types/cosmos/gov/v1beta1/gov";
import { VotingParams } from "./module/types/cosmos/gov/v1beta1/gov";
import { TallyParams } from "./module/types/cosmos/gov/v1beta1/gov";
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
        Proposal: {},
        Proposals: {},
        Vote: {},
        Votes: {},
        Params: {},
        Deposit: {},
        Deposits: {},
        TallyResult: {},
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
        getProposal: (state) => (params = {}) => {
            return state.Proposal[JSON.stringify(params)] ?? {};
        },
        getProposals: (state) => (params = {}) => {
            return state.Proposals[JSON.stringify(params)] ?? {};
        },
        getVote: (state) => (params = {}) => {
            return state.Vote[JSON.stringify(params)] ?? {};
        },
        getVotes: (state) => (params = {}) => {
            return state.Votes[JSON.stringify(params)] ?? {};
        },
        getParams: (state) => (params = {}) => {
            return state.Params[JSON.stringify(params)] ?? {};
        },
        getDeposit: (state) => (params = {}) => {
            return state.Deposit[JSON.stringify(params)] ?? {};
        },
        getDeposits: (state) => (params = {}) => {
            return state.Deposits[JSON.stringify(params)] ?? {};
        },
        getTallyResult: (state) => (params = {}) => {
            return state.TallyResult[JSON.stringify(params)] ?? {};
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
        async QueryProposal({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryProposal.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Proposal', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryProposal', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryProposal', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QueryProposals({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryProposals.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Proposals', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryProposals', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryProposals', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QueryVote({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryVote.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Vote', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryVote', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryVote', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QueryVotes({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryVotes.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Votes', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryVotes', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryVotes', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QueryParams({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryParams.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Params', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryParams', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryParams', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QueryDeposit({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryDeposit.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Deposit', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryDeposit', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryDeposit', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QueryDeposits({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryDeposits.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'Deposits', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryDeposits', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryDeposits', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async QueryTallyResult({ commit, rootGetters }, { subscribe = false, ...key }) {
            try {
                const value = (await (await initQueryClient(rootGetters)).queryTallyResult.apply(null, Object.values(key))).data;
                commit('QUERY', { query: 'TallyResult', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'QueryTallyResult', payload: key });
            }
            catch (e) {
                console.error(new SpVuexError('QueryClient:QueryTallyResult', 'API Node Unavailable. Could not perform query.'));
            }
        },
        async sendMsgVote({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgVote(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgVote:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgVote:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgSubmitProposal({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgSubmitProposal(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgSubmitProposal:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgSubmitProposal:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async sendMsgDeposit({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgDeposit(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgDeposit:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgDeposit:Send', 'Could not broadcast Tx.');
                }
            }
        },
        async MsgVote({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgVote(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgVote:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgVote:Create', 'Could not create message.');
                }
            }
        },
        async MsgSubmitProposal({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgSubmitProposal(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgSubmitProposal:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgSubmitProposal:Create', 'Could not create message.');
                }
            }
        },
        async MsgDeposit({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgDeposit(value);
                return msg;
            }
            catch (e) {
                if (e.toString() == 'wallet is required') {
                    throw new SpVuexError('TxClient:MsgDeposit:Init', 'Could not initialize signing client. Wallet is required.');
                }
                else {
                    throw new SpVuexError('TxClient:MsgDeposit:Create', 'Could not create message.');
                }
            }
        },
    }
};
