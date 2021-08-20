import { Writer, Reader } from 'protobufjs/minimal';
export declare const protobufPackage = "cosmonaut.blog.blog";
/** GenesisState defines the blog module's genesis state. */
export interface GenesisState {
}
export declare const GenesisState: {
    encode(_: GenesisState, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(_: any): GenesisState;
    toJSON(_: GenesisState): unknown;
    fromPartial(_: DeepPartial<GenesisState>): GenesisState;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
