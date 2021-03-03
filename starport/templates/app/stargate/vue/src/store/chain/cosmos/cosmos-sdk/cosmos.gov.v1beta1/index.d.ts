declare const _default: {
    namespaced: boolean;
    state: {
        getProposal: (state: any) => (params?: {}) => any;
        getProposals: (state: any) => (params?: {}) => any;
        getVote: (state: any) => (params?: {}) => any;
        getVotes: (state: any) => (params?: {}) => any;
        getParams: (state: any) => (params?: {}) => any;
        getDeposit: (state: any) => (params?: {}) => any;
        getDeposits: (state: any) => (params?: {}) => any;
        getTallyResult: (state: any) => (params?: {}) => any;
        _Structure: {
            TextProposal: {
                fields: any[];
            };
            Deposit: {
                fields: any[];
            };
            Proposal: {
                fields: any[];
            };
            TallyResult: {
                fields: any[];
            };
            Vote: {
                fields: any[];
            };
            DepositParams: {
                fields: any[];
            };
            VotingParams: {
                fields: any[];
            };
            TallyParams: {
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
        getProposal: (state: any) => (params?: {}) => any;
        getProposals: (state: any) => (params?: {}) => any;
        getVote: (state: any) => (params?: {}) => any;
        getVotes: (state: any) => (params?: {}) => any;
        getParams: (state: any) => (params?: {}) => any;
        getDeposit: (state: any) => (params?: {}) => any;
        getDeposits: (state: any) => (params?: {}) => any;
        getTallyResult: (state: any) => (params?: {}) => any;
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
        QueryProposal({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryProposals({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryVote({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryVotes({ commit, rootGetters }: {
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
        QueryDeposit({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryDeposits({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryTallyResult({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        MsgVote({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgSubmitProposal({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgDeposit({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
    };
};
export default _default;
