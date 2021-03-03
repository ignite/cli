declare const _default: {
    namespaced: boolean;
    state: {
        getConnection: (state: any) => (params?: {}) => any;
        getConnections: (state: any) => (params?: {}) => any;
        getClientConnections: (state: any) => (params?: {}) => any;
        getConnectionClientState: (state: any) => (params?: {}) => any;
        getConnectionConsensusState: (state: any) => (params?: {}) => any;
        _Structure: {
            ConnectionEnd: {
                fields: any[];
            };
            IdentifiedConnection: {
                fields: any[];
            };
            Counterparty: {
                fields: any[];
            };
            ClientPaths: {
                fields: any[];
            };
            ConnectionPaths: {
                fields: any[];
            };
            Version: {
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
        getConnection: (state: any) => (params?: {}) => any;
        getConnections: (state: any) => (params?: {}) => any;
        getClientConnections: (state: any) => (params?: {}) => any;
        getConnectionClientState: (state: any) => (params?: {}) => any;
        getConnectionConsensusState: (state: any) => (params?: {}) => any;
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
        QueryConnection({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryConnections({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryClientConnections({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryConnectionClientState({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryConnectionConsensusState({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        MsgConnectionOpenTry({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgConnectionOpenInit({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgConnectionOpenAck({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgConnectionOpenConfirm({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
    };
};
export default _default;
