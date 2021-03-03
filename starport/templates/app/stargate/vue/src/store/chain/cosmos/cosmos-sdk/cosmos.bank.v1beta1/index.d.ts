declare const _default: {
    namespaced: boolean;
    state: {
        getBalance: (state: any) => (params?: {}) => any;
        getAllBalances: (state: any) => (params?: {}) => any;
        getTotalSupply: (state: any) => (params?: {}) => any;
        getSupplyOf: (state: any) => (params?: {}) => any;
        getParams: (state: any) => (params?: {}) => any;
        getDenomMetadata: (state: any) => (params?: {}) => any;
        getDenomsMetadata: (state: any) => (params?: {}) => any;
        _Structure: {
            Params: {
                fields: any[];
            };
            SendEnabled: {
                fields: any[];
            };
            Input: {
                fields: any[];
            };
            Output: {
                fields: any[];
            };
            Supply: {
                fields: any[];
            };
            DenomUnit: {
                fields: any[];
            };
            Metadata: {
                fields: any[];
            };
            Balance: {
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
        getBalance: (state: any) => (params?: {}) => any;
        getAllBalances: (state: any) => (params?: {}) => any;
        getTotalSupply: (state: any) => (params?: {}) => any;
        getSupplyOf: (state: any) => (params?: {}) => any;
        getParams: (state: any) => (params?: {}) => any;
        getDenomMetadata: (state: any) => (params?: {}) => any;
        getDenomsMetadata: (state: any) => (params?: {}) => any;
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
        QueryBalance({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryAllBalances({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryTotalSupply({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QuerySupplyOf({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryParams({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDenomMetadata({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDenomsMetadata({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        MsgSend({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgMultiSend({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
    };
};
export default _default;
