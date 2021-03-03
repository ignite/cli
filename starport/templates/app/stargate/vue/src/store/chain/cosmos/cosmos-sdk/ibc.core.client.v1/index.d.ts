declare const _default: {
    namespaced: boolean;
    state: {
        getClientState: (state: any) => (params?: {}) => any;
        getClientStates: (state: any) => (params?: {}) => any;
        getConsensusState: (state: any) => (params?: {}) => any;
        getConsensusStates: (state: any) => (params?: {}) => any;
        getClientParams: (state: any) => (params?: {}) => any;
        _Structure: {
            IdentifiedClientState: {
                fields: any[];
            };
            ConsensusStateWithHeight: {
                fields: any[];
            };
            ClientConsensusStates: {
                fields: any[];
            };
            ClientUpdateProposal: {
                fields: any[];
            };
            Height: {
                fields: any[];
            };
            Params: {
                fields: any[];
            };
            GenesisMetadata: {
                fields: any[];
            };
            IdentifiedGenesisMetadata: {
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
        getClientState: (state: any) => (params?: {}) => any;
        getClientStates: (state: any) => (params?: {}) => any;
        getConsensusState: (state: any) => (params?: {}) => any;
        getConsensusStates: (state: any) => (params?: {}) => any;
        getClientParams: (state: any) => (params?: {}) => any;
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
        QueryClientState({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryClientStates({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryConsensusState({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryConsensusStates({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryClientParams({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        MsgSubmitMisbehaviour({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgCreateClient({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgUpdateClient({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgUpgradeClient({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
    };
};
export default _default;
