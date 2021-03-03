declare const _default: {
    namespaced: boolean;
    state: {
        getChannel: (state: any) => (params?: {}) => any;
        getChannels: (state: any) => (params?: {}) => any;
        getConnectionChannels: (state: any) => (params?: {}) => any;
        getChannelClientState: (state: any) => (params?: {}) => any;
        getChannelConsensusState: (state: any) => (params?: {}) => any;
        getPacketCommitment: (state: any) => (params?: {}) => any;
        getPacketCommitments: (state: any) => (params?: {}) => any;
        getPacketReceipt: (state: any) => (params?: {}) => any;
        getPacketAcknowledgement: (state: any) => (params?: {}) => any;
        getPacketAcknowledgements: (state: any) => (params?: {}) => any;
        getUnreceivedPackets: (state: any) => (params?: {}) => any;
        getUnreceivedAcks: (state: any) => (params?: {}) => any;
        getNextSequenceReceive: (state: any) => (params?: {}) => any;
        _Structure: {
            Channel: {
                fields: any[];
            };
            IdentifiedChannel: {
                fields: any[];
            };
            Counterparty: {
                fields: any[];
            };
            Packet: {
                fields: any[];
            };
            PacketState: {
                fields: any[];
            };
            Acknowledgement: {
                fields: any[];
            };
            PacketSequence: {
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
        getChannel: (state: any) => (params?: {}) => any;
        getChannels: (state: any) => (params?: {}) => any;
        getConnectionChannels: (state: any) => (params?: {}) => any;
        getChannelClientState: (state: any) => (params?: {}) => any;
        getChannelConsensusState: (state: any) => (params?: {}) => any;
        getPacketCommitment: (state: any) => (params?: {}) => any;
        getPacketCommitments: (state: any) => (params?: {}) => any;
        getPacketReceipt: (state: any) => (params?: {}) => any;
        getPacketAcknowledgement: (state: any) => (params?: {}) => any;
        getPacketAcknowledgements: (state: any) => (params?: {}) => any;
        getUnreceivedPackets: (state: any) => (params?: {}) => any;
        getUnreceivedAcks: (state: any) => (params?: {}) => any;
        getNextSequenceReceive: (state: any) => (params?: {}) => any;
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
        QueryChannel({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryChannels({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryConnectionChannels({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryChannelClientState({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryChannelConsensusState({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryPacketCommitment({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryPacketCommitments({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryPacketReceipt({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryPacketAcknowledgement({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryPacketAcknowledgements({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryUnreceivedPackets({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryUnreceivedAcks({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        QueryNextSequenceReceive({ commit, rootGetters }: {
            commit: any;
            rootGetters: any;
        }, { subscribe, ...key }: {
            [x: string]: any;
            subscribe?: boolean;
        }): Promise<void>;
        MsgTimeout({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgChannelOpenInit({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgTimeoutOnClose({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgAcknowledgement({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgChannelOpenTry({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgRecvPacket({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgChannelOpenAck({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgChannelCloseConfirm({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgChannelOpenConfirm({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
        MsgChannelCloseInit({ rootGetters }: {
            rootGetters: any;
        }, { value }: {
            value: any;
        }): Promise<void>;
    };
};
export default _default;
