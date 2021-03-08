import { txClient, queryClient } from './module';
import { BaseVestingAccount } from "./module/types/cosmos/vesting/v1beta1/vesting";
import { ContinuousVestingAccount } from "./module/types/cosmos/vesting/v1beta1/vesting";
import { DelayedVestingAccount } from "./module/types/cosmos/vesting/v1beta1/vesting";
import { Period } from "./module/types/cosmos/vesting/v1beta1/vesting";
import { PeriodicVestingAccount } from "./module/types/cosmos/vesting/v1beta1/vesting";
async function initTxClient(vuexGetters) {
    return await txClient(vuexGetters['chain/common/wallet/signer'], {
        addr: vuexGetters['chain/common/env/apiTendermint']
    });
}
async function initQueryClient(vuexGetters) {
    return await queryClient({
        addr: vuexGetters['chain/common/env/apiCosmos']
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
        _Structure: {
            BaseVestingAccount: getStructure(BaseVestingAccount.fromPartial({})),
            ContinuousVestingAccount: getStructure(ContinuousVestingAccount.fromPartial({})),
            DelayedVestingAccount: getStructure(DelayedVestingAccount.fromPartial({})),
            Period: getStructure(Period.fromPartial({})),
            PeriodicVestingAccount: getStructure(PeriodicVestingAccount.fromPartial({})),
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
        getTypeStructure: (state) => (type) => {
            return state._Structure[type].fields;
        }
    },
    actions: {
        init({ dispatch, rootGetters }) {
            console.log('init');
            if (rootGetters['chain/common/env/client']) {
                rootGetters['chain/common/env/client'].on('newblock', () => {
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
        async MsgCreateVestingAccount({ rootGetters }, { value }) {
            try {
                const msg = await (await initTxClient(rootGetters)).msgCreateVestingAccount(value);
                await (await initTxClient(rootGetters)).signAndBroadcast([msg]);
            }
            catch (e) {
                throw 'Failed to broadcast transaction: ' + e;
            }
        },
    }
};
