import { Reader, Writer } from 'protobufjs/minimal';
export declare const protobufPackage = "cosmonaut.blog.blog";
/** this line is used by starport scaffolding # proto/tx/message */
export interface MsgCreatePost {
    creator: string;
    title: string;
    body: string;
}
export interface MsgCreatePostResponse {
}
export declare const MsgCreatePost: {
    encode(message: MsgCreatePost, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgCreatePost;
    fromJSON(object: any): MsgCreatePost;
    toJSON(message: MsgCreatePost): unknown;
    fromPartial(object: DeepPartial<MsgCreatePost>): MsgCreatePost;
};
export declare const MsgCreatePostResponse: {
    encode(_: MsgCreatePostResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgCreatePostResponse;
    fromJSON(_: any): MsgCreatePostResponse;
    toJSON(_: MsgCreatePostResponse): unknown;
    fromPartial(_: DeepPartial<MsgCreatePostResponse>): MsgCreatePostResponse;
};
/** Msg defines the Msg service. */
export interface Msg {
    /** this line is used by starport scaffolding # proto/tx/rpc */
    CreatePost(request: MsgCreatePost): Promise<MsgCreatePostResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    CreatePost(request: MsgCreatePost): Promise<MsgCreatePostResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
