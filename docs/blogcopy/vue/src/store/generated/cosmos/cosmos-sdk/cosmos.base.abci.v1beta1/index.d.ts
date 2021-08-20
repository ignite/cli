import { TxResponse } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { ABCIMessageLog } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { StringEvent } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { Attribute } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { GasInfo } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { Result } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { SimulationResponse } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { MsgData } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { TxMsgData } from "./module/types/cosmos/base/abci/v1beta1/abci";
import { SearchTxsResult } from "./module/types/cosmos/base/abci/v1beta1/abci";
export { TxResponse, ABCIMessageLog, StringEvent, Attribute, GasInfo, Result, SimulationResponse, MsgData, TxMsgData, SearchTxsResult };
declare const _default: {
    namespaced: boolean;
    state: {
        _Structure: {
            TxResponse: {
                fields: any[];
            };
            ABCIMessageLog: {
                fields: any[];
            };
            StringEvent: {
                fields: any[];
            };
            Attribute: {
                fields: any[];
            };
            GasInfo: {
                fields: any[];
            };
            Result: {
                fields: any[];
            };
            SimulationResponse: {
                fields: any[];
            };
            MsgData: {
                fields: any[];
            };
            TxMsgData: {
                fields: any[];
            };
            SearchTxsResult: {
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
    };
};
export default _default;
