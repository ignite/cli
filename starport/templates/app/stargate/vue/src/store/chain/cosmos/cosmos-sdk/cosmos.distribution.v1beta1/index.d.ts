declare const _default: {
    namespaced: boolean;
    state: {
        getParams: (state: any) => (params?: {}) => any;
        getValidatorOutstandingRewards: (state: any) => (params?: {}) => any;
        getValidatorCommission: (state: any) => (params?: {}) => any;
        getValidatorSlashes: (state: any) => (params?: {}) => any;
        getDelegationRewards: (state: any) => (params?: {}) => any;
        getDelegationTotalRewards: (state: any) => (params?: {}) => any;
        getDelegatorValidators: (state: any) => (params?: {}) => any;
        getDelegatorWithdrawAddress: (state: any) => (params?: {}) => any;
        getCommunityPool: (state: any) => (params?: {}) => any;
        _Structure: {
            Params: {
                fields: any[];
            };
            ValidatorHistoricalRewards: {
                fields: any[];
            };
            ValidatorCurrentRewards: {
                fields: any[];
            };
            ValidatorAccumulatedCommission: {
                fields: any[];
            };
            ValidatorOutstandingRewards: {
                fields: any[];
            };
            ValidatorSlashEvent: {
                fields: any[];
            };
            ValidatorSlashEvents: {
                fields: any[];
            };
            FeePool: {
                fields: any[];
            };
            CommunityPoolSpendProposal: {
                fields: any[];
            };
            DelegatorStartingInfo: {
                fields: any[];
            };
            DelegationDelegatorReward: {
                fields: any[];
            };
            CommunityPoolSpendProposalWithDeposit: {
                fields: any[];
            };
            DelegatorWithdrawInfo: {
                fields: any[];
            };
            ValidatorOutstandingRewardsRecord: {
                fields: any[];
            };
            ValidatorAccumulatedCommissionRecord: {
                fields: any[];
            };
            ValidatorHistoricalRewardsRecord: {
                fields: any[];
            };
            ValidatorCurrentRewardsRecord: {
                fields: any[];
            };
            DelegatorStartingInfoRecord: {
                fields: any[];
            };
            ValidatorSlashEventRecord: {
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
        getParams: (state: any) => (params?: {}) => any;
        getValidatorOutstandingRewards: (state: any) => (params?: {}) => any;
        getValidatorCommission: (state: any) => (params?: {}) => any;
        getValidatorSlashes: (state: any) => (params?: {}) => any;
        getDelegationRewards: (state: any) => (params?: {}) => any;
        getDelegationTotalRewards: (state: any) => (params?: {}) => any;
        getDelegatorValidators: (state: any) => (params?: {}) => any;
        getDelegatorWithdrawAddress: (state: any) => (params?: {}) => any;
        getCommunityPool: (state: any) => (params?: {}) => any;
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
        QueryParams({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryValidatorOutstandingRewards({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryValidatorCommission({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryValidatorSlashes({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDelegationRewards({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDelegationTotalRewards({ commit, rootGetters }: {
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
        QueryDelegatorWithdrawAddress({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryCommunityPool({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        MsgFundCommunityPool({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgSetWithdrawAddress({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgWithdrawDelegatorReward({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgWithdrawValidatorCommission({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
    };
};
export default _default;
