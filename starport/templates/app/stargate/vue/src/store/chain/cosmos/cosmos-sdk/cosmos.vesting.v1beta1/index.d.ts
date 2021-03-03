declare const _default: {
    namespaced: boolean;
    state: {
        _Structure: {
            BaseVestingAccount: {
                fields: any[];
            };
            ContinuousVestingAccount: {
                fields: any[];
            };
            DelayedVestingAccount: {
                fields: any[];
            };
            Period: {
                fields: any[];
            };
            PeriodicVestingAccount: {
                fields: any[];
            };
        };
        _Subscriptions: Set<unknown>;
    };
    mutations: {
        RESET_STATE(state: any): void;
        QUERY(state: any, { query, key, value }: {
            query: any;
            key: any;
            value: any;
        }): void;
        SUBSCRIBE(state: any, subscription: any): void;
        UNSUBSCRIBE(state: any, subscription: any): void;
    };
    getters: {
        getTypeStructure: (state: any) => (type: any) => any;
    };
    actions: {
        init({ dispatch, rootGetters }: {
            dispatch: any;
            rootGetters: any;
        }): void;
        resetState({ commit }: {
            commit: any;
        }): void;
        unsubscribe({ commit }: {
            commit: any;
        }, subscription: any): void;
        StoreUpdate({ state, dispatch }: {
            state: any;
            dispatch: any;
        }): Promise<void>;
        MsgCreateVestingAccount({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
    };
};
export default _default;
