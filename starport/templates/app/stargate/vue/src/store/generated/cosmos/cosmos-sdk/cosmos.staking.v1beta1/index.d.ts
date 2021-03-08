declare const _default: {
    namespaced: boolean;
    state: {
        getValidators: (state: any) => (params?: {}) => any;
        getValidator: (state: any) => (params?: {}) => any;
        getValidatorDelegations: (state: any) => (params?: {}) => any;
        getValidatorUnbondingDelegations: (state: any) => (params?: {}) => any;
        getDelegation: (state: any) => (params?: {}) => any;
        getUnbondingDelegation: (state: any) => (params?: {}) => any;
        getDelegatorDelegations: (state: any) => (params?: {}) => any;
        getDelegatorUnbondingDelegations: (state: any) => (params?: {}) => any;
        getRedelegations: (state: any) => (params?: {}) => any;
        getDelegatorValidators: (state: any) => (params?: {}) => any;
        getDelegatorValidator: (state: any) => (params?: {}) => any;
        getHistoricalInfo: (state: any) => (params?: {}) => any;
        getPool: (state: any) => (params?: {}) => any;
        getParams: (state: any) => (params?: {}) => any;
        _Structure: {
            HistoricalInfo: {
                fields: any[];
            };
            CommissionRates: {
                fields: any[];
            };
            Commission: {
                fields: any[];
            };
            Description: {
                fields: any[];
            };
            Validator: {
                fields: any[];
            };
            ValAddresses: {
                fields: any[];
            };
            DVPair: {
                fields: any[];
            };
            DVPairs: {
                fields: any[];
            };
            DVVTriplet: {
                fields: any[];
            };
            DVVTriplets: {
                fields: any[];
            };
            Delegation: {
                fields: any[];
            };
            UnbondingDelegation: {
                fields: any[];
            };
            UnbondingDelegationEntry: {
                fields: any[];
            };
            RedelegationEntry: {
                fields: any[];
            };
            Redelegation: {
                fields: any[];
            };
            Params: {
                fields: any[];
            };
            DelegationResponse: {
                fields: any[];
            };
            RedelegationEntryResponse: {
                fields: any[];
            };
            RedelegationResponse: {
                fields: any[];
            };
            Pool: {
                fields: any[];
            };
            LastValidatorPower: {
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
        getValidators: (state: any) => (params?: {}) => any;
        getValidator: (state: any) => (params?: {}) => any;
        getValidatorDelegations: (state: any) => (params?: {}) => any;
        getValidatorUnbondingDelegations: (state: any) => (params?: {}) => any;
        getDelegation: (state: any) => (params?: {}) => any;
        getUnbondingDelegation: (state: any) => (params?: {}) => any;
        getDelegatorDelegations: (state: any) => (params?: {}) => any;
        getDelegatorUnbondingDelegations: (state: any) => (params?: {}) => any;
        getRedelegations: (state: any) => (params?: {}) => any;
        getDelegatorValidators: (state: any) => (params?: {}) => any;
        getDelegatorValidator: (state: any) => (params?: {}) => any;
        getHistoricalInfo: (state: any) => (params?: {}) => any;
        getPool: (state: any) => (params?: {}) => any;
        getParams: (state: any) => (params?: {}) => any;
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
        QueryValidators({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryValidator({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryValidatorDelegations({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryValidatorUnbondingDelegations({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDelegation({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryUnbondingDelegation({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDelegatorDelegations({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDelegatorUnbondingDelegations({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryRedelegations({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDelegatorValidators({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDelegatorValidator({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryHistoricalInfo({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryPool({ commit, rootGetters }: {
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
        MsgBeginRedelegate({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgEditValidator({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgDelegate({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgUndelegate({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgCreateValidator({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
    };
};
export default _default;
